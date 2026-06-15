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

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/session"
)

func (gs *GormStore) CreateSession(ctx context.Context, sess *session.Session) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Create(toSessionModel(sess, orgID)).Error
}

func (gs *GormStore) GetSessionByToken(ctx context.Context, token string) (*session.Session, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.Session
	err = gs.db.WithContext(ctx).Where(map[string]any{"refresh_token": token, "organization_id": orgID}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "session not found").WithCause(err)
	}
	if err != nil {
		return nil, err
	}

	sess := fromSessionModel(&m)

	var im model.Identity
	if err := gs.db.WithContext(ctx).
		Where(map[string]any{"id": m.IdentityID, "organization_id": orgID}).
		First(&im).Error; err == nil {
		if i, err := fromIdentityModel(&im); err == nil {
			sess.Identity = i
		}
	}

	return sess, nil
}

func (gs *GormStore) UpdateSession(ctx context.Context, sess *session.Session) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).Save(toSessionModel(sess, orgID)).Error
}

func (gs *GormStore) RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, oldToken, newToken string) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).
		Model(&model.Session{}).
		Where(map[string]any{
			"id":              sessionID,
			"organization_id": orgID,
			"refresh_token":   oldToken,
			"active":          true,
		}).
		Updates(map[string]any{
			"refresh_token": newToken,
			"updated_at":    time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "session not found").WithMetadata(map[string]string{"session_id": sessionID.String()})
	}
	return nil
}

func (gs *GormStore) GetSessionByID(ctx context.Context, id uuid.UUID) (*session.Session, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.Session
	err = gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "session not found").WithCause(err).WithMetadata(map[string]string{"session_id": id.String()})
	}
	if err != nil {
		return nil, err
	}
	sess := fromSessionModel(&m)

	var im model.Identity
	if err := gs.db.WithContext(ctx).
		Where(map[string]any{"id": m.IdentityID, "organization_id": orgID}).
		First(&im).Error; err == nil {
		if i, err := fromIdentityModel(&im); err == nil {
			sess.Identity = i
		}
	}
	return sess, nil
}

func (gs *GormStore) ListSessionsByIdentityID(ctx context.Context, identityID uuid.UUID) ([]*session.Session, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var models []model.Session
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"identity_id":     identityID,
			"organization_id": orgID,
			"active":          true,
		}).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	sessions := make([]*session.Session, 0, len(models))
	for i := range models {
		sessions = append(sessions, fromSessionModel(&models[i]))
	}
	return sessions, nil
}

func (gs *GormStore) RevokeAllSessionsByIdentityID(ctx context.Context, identityID uuid.UUID, exceptID uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	return gs.db.WithContext(ctx).
		Model(&model.Session{}).
		Where(map[string]any{
			"identity_id":     identityID,
			"organization_id": orgID,
			"active":          true,
		}).
		Not(map[string]any{"id": exceptID}).
		Update("active", false).Error
}

func (gs *GormStore) RevokeSession(ctx context.Context, id uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).
		Model(&model.Session{}).
		Where(map[string]any{
			"id":              id,
			"organization_id": orgID,
		}).
		Update("active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "session not found").WithMetadata(map[string]string{"session_id": id.String()})
	}
	return nil
}

func toSessionModel(sess *session.Session, orgID uuid.NullUUID) *model.Session {
	amrJSON, _ := json.Marshal(sess.AMR)
	return &model.Session{
		ID:                          sess.ID,
		OrganizationID:              orgID,
		IdentityID:                  sess.IdentityID,
		RefreshToken:                sess.RefreshToken,
		ExpiresAt:                   sess.ExpiresAt,
		IssuedAt:                    sess.IssuedAt,
		AuthenticatedAt:             sess.AuthenticatedAt,
		Active:                      sess.Active,
		AuthenticatorAssuranceLevel: string(sess.AuthenticatorAssuranceLevel),
		AMR:                         string(amrJSON),
		CreatedAt:                   sess.CreatedAt,
		UpdatedAt:                   sess.UpdatedAt,
	}
}

func fromSessionModel(m *model.Session) *session.Session {
	var amr []session.AuthenticationMethod
	if m.AMR != "" {
		_ = json.Unmarshal([]byte(m.AMR), &amr)
	}
	if amr == nil {
		amr = []session.AuthenticationMethod{}
	}

	return &session.Session{
		ID:                          m.ID,
		IdentityID:                  m.IdentityID,
		RefreshToken:                m.RefreshToken,
		ExpiresAt:                   m.ExpiresAt,
		IssuedAt:                    m.IssuedAt,
		AuthenticatedAt:             m.AuthenticatedAt,
		Active:                      m.Active,
		AuthenticatorAssuranceLevel: identity.AuthenticatorAssuranceLevel(m.AuthenticatorAssuranceLevel),
		AMR:                         amr,
		CreatedAt:                   m.CreatedAt,
		UpdatedAt:                   m.UpdatedAt,
	}
}
