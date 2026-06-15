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
	loginflow "github.com/muixstudio/clio/internal/flow/login"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
)

func (gs *GormStore) CreateLoginFlow(ctx context.Context, f *loginflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(toLoginFlowModel(f, orgID)).Error
}

func (gs *GormStore) GetLoginFlow(ctx context.Context, id uuid.UUID) (*loginflow.Flow, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.LoginFlow
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "login flow not found").WithCause(err).WithMetadata(map[string]string{"flow_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	return fromLoginFlowModel(&m)
}

func (gs *GormStore) UpdateLoginFlow(ctx context.Context, f *loginflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(toLoginFlowModel(f, orgID)).Error
}

func toLoginFlowModel(f *loginflow.Flow, orgID uuid.NullUUID) *model.LoginFlow {
	return &model.LoginFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		ExpiresAt:      f.ExpiresAt,
		Active:         string(f.Active),
		Refresh:        f.Refresh,
		CreatedAt:      f.CreatedAt,
	}
}

func fromLoginFlowModel(m *model.LoginFlow) (*loginflow.Flow, error) {
	return &loginflow.Flow{
		Base: flow.Base{
			ID:        m.ID,
			State:     flow.State(m.State),
			ExpiresAt: m.ExpiresAt,
			CreatedAt: m.CreatedAt,
		},
		Active:  identity.CredentialsType(m.Active),
		Refresh: m.Refresh,
	}, nil
}
