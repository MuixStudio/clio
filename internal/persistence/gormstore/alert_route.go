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

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/zus/dispatch"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/persistence/gormstore/model"
)

// rootAlertRoute builds the default root route for a freshly created team. It
// matches everything (empty matcher set), has no parent, and carries no
// receiver/options so that child routes start from DefaultRouteOpts. It is
// created alongside the team and may never be deleted.
func rootAlertRoute(orgID uuid.NullUUID, teamID uuid.UUID) *model.AlertRoute {
	return &model.AlertRoute{
		ID:             uuid.New(),
		OrganizationID: orgID,
		TeamID:         teamID,
		Name:           "root",
		ParentID:       uuid.NullUUID{}, // NULL => root
		Priority:       0,
		Matchers:       datatypes.JSON("[]"),
		Enabled:        true,
		Continue:       false,
	}
}

// ListAlertRoutes returns every route row for a team, scoped to the caller's
// organization. The rows are unordered; pass them to router.BuildTree (or use
// LoadAlertRouteTree) to assemble the tree.
func (gs *GormStore) ListAlertRoutes(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) ([]*model.AlertRoute, error) {
	var rows []*model.AlertRoute
	result := gs.db.WithContext(ctx).
		Where(map[string]any{"team_id": teamID, "organization_id": orgID}).
		Find(&rows)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return rows, nil
}

// LoadAlertRouteTree loads a team's routes and rebuilds the in-memory routing
// tree, assigning each node a fresh post-order Idx.
func (gs *GormStore) LoadAlertRouteTree(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) (*dispatch.Route, error) {
	rows, err := gs.ListAlertRoutes(ctx, orgID, teamID)
	if err != nil {
		return nil, err
	}
	tree, err := dispatch.BuildTree(rows)
	if err != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: invalid route tree").WithCause(err)
	}
	return tree, nil
}

// AddAlertRoute inserts a route under the given parent (or as the root when
// ParentID is NULL) and returns the persisted row. It validates the parent
// belongs to the same team, enforces a single root per team, and computes the
// new node's sibling Priority, all inside one transaction.
func (gs *GormStore) AddAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in dispatch.AddRouteInput) (uuid.UUID, error) {
	matchersJSON, err := dispatch.EncodeMatchers(in.Matchers)
	if err != nil {
		return uuid.Nil, errors.BadRequest("INVALID_MATCHERS", "invalid route matchers").WithCause(err)
	}
	groupByJSON, err := encodeJSONStrings(in.GroupBy)
	if err != nil {
		return uuid.Nil, errors.BadRequest("INVALID_GROUP_BY", "invalid group_by").WithCause(err)
	}
	muteJSON, err := encodeJSONStrings(in.MuteTimeIntervals)
	if err != nil {
		return uuid.Nil, errors.BadRequest("INVALID_MUTE_INTERVALS", "invalid mute_time_intervals").WithCause(err)
	}
	activeJSON, err := encodeJSONStrings(in.ActiveTimeIntervals)
	if err != nil {
		return uuid.Nil, errors.BadRequest("INVALID_ACTIVE_INTERVALS", "invalid active_time_intervals").WithCause(err)
	}

	id := in.RouteID
	if id == uuid.Nil {
		id = uuid.New()
	}
	row := &model.AlertRoute{
		ID:                  id,
		OrganizationID:      orgID,
		TeamID:              teamID,
		Name:                in.Name,
		ParentID:            in.ParentID,
		Matchers:            datatypes.JSON(matchersJSON),
		Enabled:             true,
		Continue:            in.Continue,
		Receiver:            in.Receiver,
		GroupBy:             groupByJSON,
		GroupByAll:          in.GroupByAll,
		GroupWait:           durationPtr(in.GroupWait),
		GroupInterval:       durationPtr(in.GroupInterval),
		RepeatInterval:      durationPtr(in.RepeatInterval),
		MuteTimeIntervals:   muteJSON,
		ActiveTimeIntervals: activeJSON,
	}

	err = gs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if in.ParentID.Valid {
			// Parent must exist within the same team/org.
			var cnt int64
			if err := tx.Model(&model.AlertRoute{}).
				Where(map[string]any{"id": in.ParentID.UUID, "team_id": teamID, "organization_id": orgID}).
				Count(&cnt).Error; err != nil {
				return err
			}
			if cnt == 0 {
				return errors.NotFound("RECORD_NOT_FOUND", "parent route not found: "+in.ParentID.UUID.String())
			}
		} else {
			// Only one root route is allowed per team.
			var cnt int64
			if err := tx.Model(&model.AlertRoute{}).
				Where("team_id = ? AND organization_id = ? AND parent_id IS NULL", teamID, orgID).
				Count(&cnt).Error; err != nil {
				return err
			}
			if cnt > 0 {
				return errors.Conflict("ROOT_ROUTE_EXISTS", "root route already exists for team")
			}
		}

		// Append after existing siblings: priority = max(sibling priority) + 1.
		var maxPriority *int
		if err := tx.Model(&model.AlertRoute{}).
			Where(map[string]any{"team_id": teamID, "organization_id": orgID, "parent_id": in.ParentID}).
			Select("MAX(priority)").
			Scan(&maxPriority).Error; err != nil {
			return err
		}
		if maxPriority != nil {
			row.Priority = *maxPriority + 1
		}

		return tx.Create(row).Error
	})
	if err != nil {
		var appErr *errors.Error
		if errors.As(err, &appErr) {
			return uuid.Nil, appErr
		}
		return uuid.Nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	return row.ID, nil
}

// DeleteAlertRoute removes a single route node. The root route (ParentID NULL)
// can never be deleted, and a node that still has children is rejected too
// (no cascade) — callers must move or delete the children first. The whole
// check-and-delete runs in one transaction.
// UpdateAlertRoute updates mutable fields of one route node. ParentID/Priority
// are intentionally not changed here; moving/reordering routes should be a
// separate operation that can validate cycles and sibling ordering explicitly.
func (gs *GormStore) UpdateAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in dispatch.UpdateRouteInput) error {
	updates := map[string]any{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Matchers != nil {
		matchersJSON, err := dispatch.EncodeMatchers(*in.Matchers)
		if err != nil {
			return errors.BadRequest("INVALID_MATCHERS", "invalid route matchers").WithCause(err)
		}
		updates["matchers"] = datatypes.JSON(matchersJSON)
	}
	if in.Continue != nil {
		updates["continue"] = *in.Continue
	}
	if in.Receiver != nil {
		updates["receiver"] = *in.Receiver
	}
	if in.GroupBy != nil {
		groupByJSON, err := encodeJSONStrings(*in.GroupBy)
		if err != nil {
			return errors.BadRequest("INVALID_GROUP_BY", "invalid group_by").WithCause(err)
		}
		updates["group_by"] = groupByJSON
	}
	if in.GroupByAll != nil {
		updates["group_by_all"] = *in.GroupByAll
	}
	if in.GroupWait != nil {
		updates["group_wait"] = int64(*in.GroupWait)
	}
	if in.GroupInterval != nil {
		updates["group_interval"] = int64(*in.GroupInterval)
	}
	if in.RepeatInterval != nil {
		updates["repeat_interval"] = int64(*in.RepeatInterval)
	}
	if in.MuteTimeIntervals != nil {
		muteJSON, err := encodeJSONStrings(*in.MuteTimeIntervals)
		if err != nil {
			return errors.BadRequest("INVALID_MUTE_INTERVALS", "invalid mute_time_intervals").WithCause(err)
		}
		updates["mute_time_intervals"] = muteJSON
	}
	if in.ActiveTimeIntervals != nil {
		activeJSON, err := encodeJSONStrings(*in.ActiveTimeIntervals)
		if err != nil {
			return errors.BadRequest("INVALID_ACTIVE_INTERVALS", "invalid active_time_intervals").WithCause(err)
		}
		updates["active_time_intervals"] = activeJSON
	}
	if len(updates) == 0 {
		return nil
	}

	result := gs.db.WithContext(ctx).Model(&model.AlertRoute{}).
		Where(map[string]any{"id": in.RouteID, "team_id": teamID, "organization_id": orgID}).
		Updates(updates)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "route not found: "+in.RouteID.String())
	}
	return nil
}

func (gs *GormStore) DeleteAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, routeID uuid.UUID) error {
	err := gs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var route model.AlertRoute
		if err := tx.Where(map[string]any{
			"id":              routeID,
			"team_id":         teamID,
			"organization_id": orgID,
		}).First(&route).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.NotFound("RECORD_NOT_FOUND", "route not found: "+routeID.String())
			}
			return err
		}

		if !route.ParentID.Valid {
			return errors.BadRequest("ROOT_ROUTE_UNDELETABLE", "root route cannot be deleted")
		}

		var childCount int64
		if err := tx.Model(&model.AlertRoute{}).
			Where(map[string]any{"parent_id": routeID, "organization_id": orgID}).
			Count(&childCount).Error; err != nil {
			return err
		}
		if childCount > 0 {
			return errors.Conflict("ROUTE_HAS_CHILDREN", "route has child routes; remove or move them first")
		}

		return tx.Delete(&model.AlertRoute{}, "id = ?", routeID).Error
	})
	if err != nil {
		var appErr *errors.Error
		if errors.As(err, &appErr) {
			return appErr
		}
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	return nil
}

// encodeJSONStrings serializes a string slice into the datatypes.JSON form used
// by the route columns. A nil/empty slice yields a nil payload, which the
// router decodes as "inherit from parent".
func durationPtr(d *time.Duration) *int64 {
	if d == nil {
		return nil
	}
	v := int64(*d)
	return &v
}

func encodeJSONStrings(values []string) (datatypes.JSON, error) {
	if len(values) == 0 {
		return nil, nil
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}
