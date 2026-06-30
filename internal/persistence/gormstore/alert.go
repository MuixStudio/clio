/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gormstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
	alert2 "github.com/muixstudio/clio/internal/domain/alert/repository"
	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultFlapWindow = 5 * time.Minute

// Save 将标准化之后的告警写入三张表
// 幂等：重复推送不会产生新记录，但会记录重复次数
func (gs *GormStore) Save(ctx context.Context, alert *entity.NormalizedAlert) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.Transaction(func(tx *gorm.DB) error {
		return saveAlert(tx, orgID, alert)
	})
}

// SaveAlerts 批量保存标准化之后的告警，整批在同一个事务中完成
// 幂等：重复推送不会产生新记录，但会记录重复次数
func (gs *GormStore) SaveAlerts(ctx context.Context, alerts []*entity.NormalizedAlert) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.Transaction(func(tx *gorm.DB) error {
		for _, a := range alerts {
			if err := saveAlert(tx, orgID, a); err != nil {
				return err
			}
		}
		return nil
	})
}

func saveAlert(tx *gorm.DB, oid uuid.NullUUID, alert *entity.NormalizedAlert) error {
	defID, err := upsertDefinition(tx, oid, alert)
	if err != nil {
		return err
	}

	switch alert.Status {
	case "firing":
		return handleFiring(tx, oid, defID, alert)
	case "resolved":
		return handleResolved(tx, oid, alert)
	default:
		return errors.BadRequest("INVALID_ALERT_STATUS", "unknown alert status: "+string(alert.Status))
	}
}

// upsertDefinition 将告警定义写入 alert_definitions 表
// 逻辑：
//   - 以 (fingerprint, connector_id) 为唯一键
//   - 第一次见到这个告警：INSERT 完整记录，first_seen_at = last_seen_at = 当前 StartsAt
//   - 已存在：只更新 last_seen_at，其余字段（labels、severity 等）保持原样
//     原因：labels 是 fingerprint 的来源，理论上不会变；变了就是另一个 fingerprint
//
// 注意：
//   - ID 每次都生成新值；OnConflict 命中已存在记录时，数据库里的 id 不会被改写
//   - 因此用 clause.Returning 把数据库真实的 id 回填到 def.ID：
//     无论是新插入还是命中冲突（DO UPDATE 会"触碰"该行，RETURNING 返回的就是既有行），
//     def.ID 都指向 alert_definition 表中真实存在的主键，可作为返回值透传给 Alert.AlertDefinitionID。
//     （注意 PostgreSQL 语义：DO UPDATE 才会 RETURNING 命中行，DO NOTHING 不会。）
func upsertDefinition(tx *gorm.DB, oid uuid.NullUUID, alert *entity.NormalizedAlert) (uuid.UUID, error) {
	l, err := json.Marshal(alert.Labels)
	if err != nil {
		return uuid.Nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	var j datatypes.JSON = l

	def := model.AlertDefinition{
		ID:             uuid.New(),
		Fingerprint:    alert.Fingerprint,
		AlertName:      alert.Labels["alertname"],
		Labels:         j,
		Severity:       string(alert.Severity),
		ConnectorID:    alert.ConnectorID,
		OrganizationID: oid,
		TeamID:         alert.TeamID,
		FirstSeenAt:    alert.StartsAt,
		LastSeenAt:     alert.StartsAt,
	}

	if err := tx.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "fingerprint"}, {Name: "connector_id"}},
			// 已存在时只更新 last_seen_at，其他字段保持原样
			DoUpdates: clause.AssignmentColumns([]string{"last_seen_at"}),
		},
		// 把真实落库的 id 回填到 def.ID
		clause.Returning{Columns: []clause.Column{{Name: "id"}}},
	).Create(&def).Error; err != nil {
		return uuid.Nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	return def.ID, nil
}

// handleFiring 处理 status=firing 的告警
// 分三条路径：
//
//	路径一：重复推送（最常见）
//	  条件：(fingerprint, starts_at, organization_id) 已存在，且 status=firing
//	  处理：UPDATE repeat_count+1（原子自增），不写 alert_events
//	  原因：Prometheus 每隔 repeat_interval 会重复推送同一条 firing 告警；
//	        重复推送零信息，写事件表会撑爆量最大的 alert_events，只累加计数即可
//
//	路径二：reviveAlert —— 同一次触发从 resolved 变回 firing
//	  条件：(fingerprint, starts_at, organization_id) 已存在，但 status=resolved
//	  处理：UPDATE status=firing, ends_at=NULL, flap_count+1
//	        INSERT alert_events(flapping)
//	  原因：startsAt 没变说明 Prometheus 认为是同一次触发，
//	        中间出现 resolved 是短暂抖动，Alertmanager 的 resolve_timeout 触发的
//
//	路径三：全新触发 or 抖动窗口内的新 firing
//	  条件：(fingerprint, starts_at, organization_id) 不存在
//	  继续判断 isFlapping：
//	    isFlapping=true（距上次 resolved < FlapWindow=5min）：
//	      → handleFlapping：复活上一条 alert，不新建
//	        UPDATE status=firing, ends_at=NULL, flap_count+1
//	        INSERT alert_events(flapping)
//	    isFlapping=false（没有历史 or 距上次 resolved >= 5min）：
//	      → createNewAlert：新建一条完整的触发记录
//	        INSERT alerts
//	        INSERT alert_events(firing)
func handleFiring(tx *gorm.DB, oid uuid.NullUUID, defID uuid.UUID, alert *entity.NormalizedAlert) error {
	// 先检查是否已存在完全相同的 (fingerprint, starts_at) 记录
	var existing model.Alert
	err := tx.Where(map[string]any{
		"fingerprint":     alert.Fingerprint,
		"starts_at":       alert.StartsAt,
		"organization_id": oid,
	}).First(&existing).Error

	if err == nil {
		// 记录已存在
		if existing.Status == "firing" {
			// 重复推送：只把这次触发的重复推送计数 +1，不写 alert_events。
			// 用 gorm.Expr 原子自增（repeat_count = repeat_count + 1），避免并发推送丢计数。
			if err := tx.Model(&existing).
				Update("repeat_count", gorm.Expr("repeat_count + 1")).Error; err != nil {
				return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
			}
			return nil
		}
		// 同一次触发：状态是 resolved 但又收到 firing（短暂抖动后 Prometheus 重新推）
		return reviveAlert(tx, oid, &existing, alert)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	// 记录不存在，判断是否属于抖动窗口内的新 firing
	if isFlapping(tx, oid, alert) {
		return handleFlapping(tx, oid, alert)
	}

	// 全新触发
	return createNewAlert(tx, oid, defID, alert)
}

// handleResolved 处理 status=resolved 的告警
// 逻辑：
//   - 查找同 fingerprint、同 organization_id、ends_at IS NULL 的记录（当前 firing 中的告警）
//   - 找到：UPDATE status=resolved, ends_at=实际结束时间
//     INSERT alert_events(resolved)
//   - 找不到：静默忽略，return nil
//     原因：服务重启后可能收到遗留的 resolved 推送，对应的 firing 记录不存在或已 resolved
//
// ends_at 取值优先级：
//  1. alert.EndsAt 不为 nil：用 Prometheus/Zabbix 给的实际恢复时间
//  2. alert.EndsAt 为 nil：兜底用 time.Now()
//     注意：不能用 alert.StartsAt 作为兜底，StartsAt 是告警开始时间，不是结束时间
func handleResolved(tx *gorm.DB, oid uuid.NullUUID, alert *entity.NormalizedAlert) error {
	// 找到当前 firing 中的记录
	var existing model.Alert
	err := tx.Where(map[string]any{
		"fingerprint":     alert.Fingerprint,
		"ends_at":         nil,
		"organization_id": oid,
	}).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 收到 resolved 但没有对应的 firing 记录
		// 可能是服务重启后收到的遗留推送，忽略
		return nil
	}
	if err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	endsAt := time.Now()
	if alert.EndsAt != nil {
		endsAt = *alert.EndsAt
	}

	// 更新 alert 状态
	if err := tx.Model(&existing).Updates(map[string]any{
		"status":  "resolved",
		"ends_at": endsAt,
	}).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	// 写入 resolved 事件
	// occurred_at 取恢复时间 endsAt，而非 StartsAt：
	// resolved 事件发生在告警恢复的那一刻，用 StartsAt 会让它和 firing 事件挤在同一时间点，
	// 破坏 alert_events 的时间线排序（ORDER BY occurred_at）。
	return insertEvent(tx, oid, existing.ID, model.AlertEventResolved, endsAt, alert.RawPayload)
}

// createNewAlert 创建一条全新的告警触发记录
// 调用时机：handleFiring 中确认是全新触发（非重复推送、非抖动）时
// 逻辑：
//
//	INSERT alerts {
//	    status    = firing
//	    ends_at   = NULL        ← nil 表示仍在 firing，resolved 后由 handleResolved 回填
//	    flap_count = 0          ← 初始为 0，每次抖动 +1
//	    annotations = alert.Annotations  ← 注意：这里应存 Annotations 而非 Labels
//	                                        Labels 已经存在 alert_definitions 里了
//	}
//	INSERT alert_events(firing)  ← 记录这次触发的起点，raw_payload 保留原始推送数据
func createNewAlert(tx *gorm.DB, oid uuid.NullUUID, defID uuid.UUID, alert *entity.NormalizedAlert) error {
	l, err := json.Marshal(alert.Labels)
	if err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	a, err := json.Marshal(alert.Annotations)
	if err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	var ja datatypes.JSON = a
	var jl datatypes.JSON = l

	newAlert := model.Alert{
		ID:                alert.ID,
		AlertDefinitionID: defID,
		Fingerprint:       alert.Fingerprint,
		Status:            "firing",
		Severity:          string(alert.Severity),
		Labels:            jl,
		Annotations:       ja,
		StartsAt:          alert.StartsAt,
		EndsAt:            nil,
		FlapCount:         0,
		ConnectorID:       alert.ConnectorID,
		OrganizationID:    oid,
		TeamID:            alert.TeamID,
	}

	if err := tx.Create(&newAlert).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	return insertEvent(tx, oid, newAlert.ID, model.AlertEventFiring, alert.StartsAt, alert.RawPayload)
}

// isFlapping
// 判断是否在抖动窗口内
// 条件：同 fingerprint 最近一次 alert 在 FlapWindow 内刚 resolved
func isFlapping(tx *gorm.DB, oid uuid.NullUUID, alert *entity.NormalizedAlert) bool {
	var lastAlert model.Alert
	err := tx.Where(map[string]any{
		"fingerprint":     alert.Fingerprint,
		"organization_id": oid,
	}).
		Not(map[string]any{"ends_at": nil}).
		Order("ends_at DESC").
		First(&lastAlert).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		return false
	}

	return alert.StartsAt.Sub(*lastAlert.EndsAt) < DefaultFlapWindow
}

// handleFlapping
// 抖动：复活上一条 alert，flap_count +1
func handleFlapping(tx *gorm.DB, oid uuid.NullUUID, alert *entity.NormalizedAlert) error {
	var lastAlert model.Alert
	if err := tx.Where(map[string]any{
		"fingerprint":     alert.Fingerprint,
		"organization_id": oid,
	}).
		Not(map[string]any{"ends_at": nil}).
		Order("ends_at DESC").
		First(&lastAlert).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	// 复活：清空 ends_at，状态改回 firing，flap_count +1
	if err := tx.Model(&lastAlert).Updates(map[string]any{
		"status":     "firing",
		"ends_at":    nil,
		"flap_count": lastAlert.FlapCount + 1,
	}).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	return insertEvent(tx, oid, lastAlert.ID, model.AlertEventFlapping, alert.StartsAt, alert.RawPayload)
}

// reviveAlert
// 同一次触发（相同 starts_at）从 resolved 变回 firing
// 这种情况只在 Prometheus 重复推送时出现
func reviveAlert(tx *gorm.DB, oid uuid.NullUUID, existing *model.Alert, alert *entity.NormalizedAlert) error {
	if err := tx.Model(existing).Updates(map[string]any{
		"status":     "firing",
		"ends_at":    nil,
		"flap_count": existing.FlapCount + 1,
	}).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	return insertEvent(tx, oid, existing.ID, model.AlertEventFlapping, alert.StartsAt, alert.RawPayload)
}

// insertEvent
func insertEvent(tx *gorm.DB, oid uuid.NullUUID, alertID uuid.UUID, eventType model.AlertEventType, occurredAt time.Time, rawPayload []byte) error {
	event := model.AlertEvent{
		ID:             uuid.New(),
		AlertID:        alertID,
		OrganizationID: oid,
		EventType:      eventType,
		OccurredAt:     occurredAt,
		RawPayload:     rawPayload,
	}
	if err := tx.Create(&event).Error; err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	return nil
}

func (gs *GormStore) ListAlert(ctx context.Context, opts *alert2.ListOptions) (alerts []*entity.NormalizedAlert, total int, err error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}

	if opts == nil {
		opts = &alert2.ListOptions{}
	}

	page := opts.Page
	if page <= 0 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := gs.db.WithContext(ctx).Model(&model.Alert{}).Where(map[string]any{
		"organization_id": orgID,
	})
	if opts.TeamID != nil {
		query = query.Where(map[string]any{"team_id": *opts.TeamID})
	}
	if opts.ConnectorID != nil {
		query = query.Where(map[string]any{"connector_id": *opts.ConnectorID})
	}
	if opts.Status != nil {
		query = query.Where(map[string]any{"status": string(*opts.Status)})
	}
	if opts.Severity != nil {
		query = query.Where(map[string]any{"severity": string(*opts.Severity)})
	}
	if opts.StartAt != nil {
		query = query.Where("starts_at >= ?", *opts.StartAt)
	}
	if opts.EndAt != nil {
		query = query.Where("starts_at <= ?", *opts.EndAt)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	var records []model.Alert
	if err := query.
		Order("starts_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	alerts = make([]*entity.NormalizedAlert, 0, len(records))
	for _, record := range records {
		labels := map[string]string{}
		if len(record.Labels) > 0 {
			if err := json.Unmarshal(record.Labels, &labels); err != nil {
				return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
			}
		}

		annotations := map[string]string{}
		if len(record.Annotations) > 0 {
			if err := json.Unmarshal(record.Annotations, &annotations); err != nil {
				return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
			}
		}

		alerts = append(alerts, &entity.NormalizedAlert{
			ID:          record.ID,
			Fingerprint: record.Fingerprint,
			ConnectorID: record.ConnectorID,
			TeamID:      record.TeamID,
			Status:      entity.Status(record.Status),
			Severity:    entity.Severity(record.Severity),
			Labels:      labels,
			Annotations: annotations,
			StartsAt:    record.StartsAt,
			EndsAt:      record.EndsAt,
			CreatedAt:   record.CreatedAt,
			UpdatedAt:   record.UpdatedAt,
		})
	}

	return alerts, int(count), nil
}
