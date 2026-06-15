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
	recflow "github.com/muixstudio/clio/internal/flow/recovery"
	verflow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/infra/errors"
)

func (gs *GormStore) CreateRecoveryFlow(ctx context.Context, f *recflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(&model.RecoveryFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		Email:          f.Email,
		ExpiresAt:      f.ExpiresAt,
		CreatedAt:      f.CreatedAt,
	}).Error
}

func (gs *GormStore) GetRecoveryFlow(ctx context.Context, id uuid.UUID) (*recflow.Flow, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.RecoveryFlow
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "recovery flow not found").WithCause(err).WithMetadata(map[string]string{"flow_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	return &recflow.Flow{
		Base: flow.Base{
			ID:        m.ID,
			State:     flow.State(m.State),
			ExpiresAt: m.ExpiresAt,
			CreatedAt: m.CreatedAt,
		},
		Email: m.Email,
	}, nil
}

func (gs *GormStore) UpdateRecoveryFlow(ctx context.Context, f *recflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(&model.RecoveryFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		Email:          f.Email,
		ExpiresAt:      f.ExpiresAt,
		CreatedAt:      f.CreatedAt,
	}).Error
}

func (gs *GormStore) CreateVerificationFlow(ctx context.Context, f *verflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(&model.VerificationFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		Email:          f.Email,
		ExpiresAt:      f.ExpiresAt,
		CreatedAt:      f.CreatedAt,
	}).Error
}

func (gs *GormStore) GetVerificationFlow(ctx context.Context, id uuid.UUID) (*verflow.Flow, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.VerificationFlow
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "verification flow not found").WithCause(err).WithMetadata(map[string]string{"flow_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	return &verflow.Flow{
		Base: flow.Base{
			ID:        m.ID,
			State:     flow.State(m.State),
			ExpiresAt: m.ExpiresAt,
			CreatedAt: m.CreatedAt,
		},
		Email: m.Email,
	}, nil
}

func (gs *GormStore) UpdateVerificationFlow(ctx context.Context, f *verflow.Flow) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(&model.VerificationFlow{
		ID:             f.ID,
		OrganizationID: orgID,
		State:          string(f.State),
		Email:          f.Email,
		ExpiresAt:      f.ExpiresAt,
		CreatedAt:      f.CreatedAt,
	}).Error
}
