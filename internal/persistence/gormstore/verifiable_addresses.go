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
	"fmt"
	"time"

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
)

func (gs *GormStore) FindIdentityByEmail(ctx context.Context, email string) (*identity.Identity, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var ci model.CredentialIdentifier
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"identifier":      email,
			"type":            []string{string(identity.CredentialsTypePassword), string(identity.CredentialsTypeCode)},
		}).
		First(&ci).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", fmt.Sprintf("identity not found for email: %s", email)).WithCause(err).WithMetadata(map[string]string{"email": email})
	}
	if err != nil {
		return nil, err
	}

	var cm model.Credential
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"id":              ci.CredentialID,
		}).
		First(&cm).Error
	if err != nil {
		return nil, err
	}

	var m model.Identity
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"id":              cm.IdentityID,
		}).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return fromIdentityModel(&m)
}

func (gs *GormStore) FindOrCreateVerifiableAddress(ctx context.Context, identityID uuid.UUID, via, value string) (*identity.VerifiableAddress, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.VerifiableAddress
	err = gs.db.WithContext(ctx).
		Where(map[string]any{"organization_id": orgID, "via": via, "value": value}).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		now := time.Now()
		m = model.VerifiableAddress{
			ID:             uuid.New(),
			OrganizationID: orgID,
			IdentityID:     identityID,
			Value:          value,
			Via:            via,
			Verified:       false,
			Status:         string(identity.VerifiableAddressStatusPending),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if createErr := gs.db.WithContext(ctx).Create(&m).Error; createErr != nil {
			return nil, createErr
		}
	} else if err != nil {
		return nil, err
	}

	return fromVerifiableAddressModel(&m), nil
}

func (gs *GormStore) MarkVerifiableAddressVerified(ctx context.Context, id uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	result := gs.db.WithContext(ctx).
		Model(&model.VerifiableAddress{}).
		Where(map[string]any{
			"id":              id,
			"organization_id": orgID,
		}).
		Updates(map[string]any{
			"verified":    true,
			"verified_at": now,
			"status":      string(identity.VerifiableAddressStatusCompleted),
			"updated_at":  now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "verifiable address not found").WithMetadata(map[string]string{"address_id": id.String()})
	}
	return nil
}

func fromVerifiableAddressModel(m *model.VerifiableAddress) *identity.VerifiableAddress {
	return &identity.VerifiableAddress{
		ID:         m.ID,
		Value:      m.Value,
		Via:        m.Via,
		Verified:   m.Verified,
		VerifiedAt: m.VerifiedAt,
		Status:     identity.VerifiableAddressStatus(m.Status),
		IdentityID: m.IdentityID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
