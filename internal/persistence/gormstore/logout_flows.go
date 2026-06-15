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

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/flow"
	logoutflow "github.com/muixstudio/clio/internal/flow/logout"
	"github.com/muixstudio/clio/internal/infra/errors"
)

func (gs *GormStore) CreateLogoutFlow(ctx context.Context, f *logoutflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(toLogoutFlowModel(f, orgID)).Error
}

func (gs *GormStore) GetLogoutFlow(ctx context.Context, id uuid.UUID) (*logoutflow.Flow, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.LogoutFlow
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "logout flow not found").WithCause(err).WithMetadata(map[string]string{"flow_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	return fromLogoutFlowModel(&m)
}

func (gs *GormStore) UpdateLogoutFlow(ctx context.Context, f *logoutflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(toLogoutFlowModel(f, orgID)).Error
}

func toLogoutFlowModel(f *logoutflow.Flow, orgID uuid.NullUUID) *model.LogoutFlow {
	return &model.LogoutFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		IdentityID:     f.IdentityID,
		State:          string(f.State),
		ExpiresAt:      f.ExpiresAt,
		CreatedAt:      f.CreatedAt,
	}
}

func fromLogoutFlowModel(m *model.LogoutFlow) (*logoutflow.Flow, error) {
	return &logoutflow.Flow{
		Base: flow.Base{
			ID:        m.ID,
			State:     flow.State(m.State),
			ExpiresAt: m.ExpiresAt,
			CreatedAt: m.CreatedAt,
		},
		IdentityID: m.IdentityID,
	}, nil
}
