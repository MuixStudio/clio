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
	regflow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
)

func (gs *GormStore) CreateRegistrationFlow(ctx context.Context, f *regflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(toRegistrationFlowModel(f, orgID)).Error
}

func (gs *GormStore) GetRegistrationFlow(ctx context.Context, id uuid.UUID) (*regflow.Flow, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.RegistrationFlow
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "registration flow not found").WithCause(err).WithMetadata(map[string]string{"flow_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	return fromRegistrationFlowModel(&m)
}

func (gs *GormStore) UpdateRegistrationFlow(ctx context.Context, f *regflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(toRegistrationFlowModel(f, orgID)).Error
}

func toRegistrationFlowModel(f *regflow.Flow, orgID uuid.NullUUID) *model.RegistrationFlow {
	return &model.RegistrationFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		ExpiresAt:      f.ExpiresAt,
		Active:         string(f.Active),
		CreatedAt:      f.CreatedAt,
	}
}

func fromRegistrationFlowModel(m *model.RegistrationFlow) (*regflow.Flow, error) {
	return &regflow.Flow{
		Base: flow.Base{
			ID:        m.ID,
			State:     flow.State(m.State),
			ExpiresAt: m.ExpiresAt,
			CreatedAt: m.CreatedAt,
		},
		Active: identity.CredentialsType(m.Active),
	}, nil
}
