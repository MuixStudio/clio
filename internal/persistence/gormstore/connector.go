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

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/connector"
	"github.com/muixstudio/clio/internal/infra/errors"
	"gorm.io/gorm"
)

func toConnectorModel(c *connector.Connector) *model.Connector {
	labels, _ := json.Marshal(c.Labels)
	return &model.Connector{
		ID:      c.ID,
		Type:    c.Type,
		Name:    c.Name,
		Token:   c.Token,
		TeamID:  c.TeamID,
		Labels:  labels,
		Enabled: c.Enabled,
	}
}

func fromConnectorModel(m *model.Connector) *connector.Connector {
	var labels map[string]string
	_ = json.Unmarshal(m.Labels, &labels)
	return &connector.Connector{
		ID:        m.ID,
		Type:      m.Type,
		Name:      m.Name,
		Token:     m.Token,
		TeamID:    m.TeamID,
		Labels:    labels,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (gs *GormStore) CreateConnector(ctx context.Context, c *connector.Connector) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	if c.ID == (uuid.UUID{}) {
		c.ID = uuid.New()
	}
	m := toConnectorModel(c)
	m.OrganizationID = orgID
	result := gs.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return nil
}

func (gs *GormStore) GetConnector(ctx context.Context, id uuid.UUID) (*connector.Connector, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.Connector
	result := gs.db.WithContext(ctx).Where(
		map[string]any{
			"id":              id,
			"organization_id": orgID,
		},
	).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("RECORD_NOT_FOUND", "connector not found: "+id.String())
		}
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return fromConnectorModel(&m), nil
}

func (gs *GormStore) UpdateConnector(ctx context.Context, id uuid.UUID, update *connector.Update) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	updates := map[string]any{}
	if update != nil {
		if update.Name != nil {
			updates["name"] = *update.Name
		}
		if update.Enabled != nil {
			updates["enabled"] = *update.Enabled
		}
	}

	if len(updates) == 0 {
		return nil
	}

	result := gs.db.WithContext(ctx).Model(&model.Connector{}).
		Where(map[string]any{
			"id":              id,
			"organization_id": orgID,
		}).
		Updates(updates)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "connector not found: "+id.String())
	}
	return nil
}

func (gs *GormStore) Count(ctx context.Context, options connector.CountOptions) (int64, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	conditions := map[string]any{
		"organization_id": orgID,
	}
	if options.Type != nil {
		conditions["type"] = *options.Type
	}
	if options.Enable != nil {
		conditions["enabled"] = *options.Enable
	}
	if options.TeamID != nil {
		conditions["team_id"] = *options.TeamID
	}

	var count int64
	result := gs.db.WithContext(ctx).Model(&model.Connector{}).
		Where(conditions).
		Count(&count)
	if result.Error != nil {
		return 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return count, nil
}

func (gs *GormStore) ListConnectors(ctx context.Context, options *connector.ListOptions) ([]*connector.Connector, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	conditions := map[string]any{
		"organization_id": orgID,
	}
	if options.Enable != nil {
		conditions["enabled"] = *options.Enable
	}
	if options.Type != nil {
		conditions["type"] = *options.Type
	}
	if options.TeamID != nil {
		conditions["team_id"] = *options.TeamID
	}

	query := gs.db.WithContext(ctx).Where(conditions)
	if options.PageSize > 0 {
		query = query.Limit(options.PageSize)
		if options.Page > 0 {
			query = query.Offset((options.Page - 1) * options.PageSize)
		}
	}

	var models []model.Connector
	result := query.Find(&models)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	connectors := make([]*connector.Connector, 0, len(models))
	for i := range models {
		connectors = append(connectors, fromConnectorModel(&models[i]))
	}
	return connectors, nil
}

func (gs *GormStore) DeleteConnector(ctx context.Context, id uuid.UUID, teamID uuid.UUID) error {

	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
		"team_id":         teamID,
	}).Delete(&model.Connector{})
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "team not found: "+id.String())
	}
	return nil

}
