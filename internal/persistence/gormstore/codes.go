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
	"time"

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/session"
	codestrategy "github.com/muixstudio/clio/internal/strategy/code"
)

func (gs *GormStore) CreateTokenExchangeCode(ctx context.Context, c *session.TokenExchangeCode) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	m := &model.TokenExchangeCode{
		ID:             c.ID,
		OrganizationID: orgID,
		Code:           c.Code,
		SessionID:      c.SessionID,
		ExpiresAt:      c.ExpiresAt,
		UsedAt:         c.UsedAt,
		CreatedAt:      c.CreatedAt,
	}
	return gs.db.WithContext(ctx).Create(m).Error
}

func (gs *GormStore) ConsumeTokenExchangeCode(ctx context.Context, code string) (*session.TokenExchangeCode, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var m model.TokenExchangeCode
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"code":            code,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "exchange code not found").WithCause(err)
	}
	if err != nil {
		return nil, err
	}
	if now.After(m.ExpiresAt) {
		return nil, errors.Gone("CODE_EXPIRED", "exchange code has expired")
	}
	if m.UsedAt != nil {
		return nil, errors.BadRequest("CODE_ALREADY_USED", "exchange code has already been used")
	}

	result := gs.db.WithContext(ctx).
		Model(&model.TokenExchangeCode{}).
		Where(map[string]any{
			"id":              m.ID,
			"organization_id": orgID,
			"used_at":         nil,
		}).
		Update("used_at", now)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.BadRequest("CODE_ALREADY_USED", "exchange code has already been used")
	}

	m.UsedAt = &now
	return &session.TokenExchangeCode{
		ID:        m.ID,
		SessionID: m.SessionID,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (gs *GormStore) CreateVerificationCode(ctx context.Context, c *codestrategy.VerificationCode) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	m := &model.VerificationCode{
		ID:             c.ID,
		OrganizationID: orgID,
		FlowID:         c.FlowID,
		FlowType:       c.FlowType,
		Address:        c.Address,
		Code:           c.Code,
		ExpiresAt:      c.ExpiresAt,
		UsedAt:         c.UsedAt,
	}
	return gs.db.WithContext(ctx).Create(m).Error
}

func (gs *GormStore) FindVerificationCode(ctx context.Context, flowID uuid.UUID, code string) (*codestrategy.VerificationCode, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.VerificationCode
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"flow_id":         flowID,
			"organization_id": orgID,
			"code":            code,
		}).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "verification code not found").WithCause(err).WithMetadata(map[string]string{"flow_id": flowID.String()})
	}
	if err != nil {
		return nil, err
	}
	return &codestrategy.VerificationCode{
		ID:        m.ID,
		FlowID:    m.FlowID,
		FlowType:  m.FlowType,
		Address:   m.Address,
		Code:      m.Code,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
	}, nil
}

func (gs *GormStore) UseVerificationCode(ctx context.Context, id uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	return gs.db.WithContext(ctx).
		Model(&model.VerificationCode{}).
		Where(map[string]any{
			"id":              id,
			"organization_id": orgID,
		}).
		Update("used_at", now).Error
}
