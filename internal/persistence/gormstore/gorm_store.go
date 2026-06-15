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
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/orgctx"
	"gorm.io/gorm"
)

type GormStore struct{ db *gorm.DB }

func (gs *GormStore) orgIDFromCtx(ctx context.Context) (uuid.NullUUID, error) {
	id, ok := orgctx.OrgIDFromCtx(ctx)
	if !ok {
		return uuid.NullUUID{}, errors.Unauthorized("MISSING_ORG", "organization context is required")
	}
	return id, nil
}

func NewGormStore(db *gorm.DB) (*GormStore, error) {
	err := migrate(db)
	if err != nil {
		return nil, err
	}
	return &GormStore{db: db}, nil
}

// Migrate runs AutoMigrate for all persistence models.
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Connector{},
		&model.Team{},
		&model.TeamMember{},
		&model.Channel{},
		&model.ChannelMember{},
		&model.AlertRoute{},
		&model.AlertRouteChannel{},

		&model.LoginFlow{},
		&model.LogoutFlow{},
		&model.VerificationFlow{},
		&model.RegistrationFlow{},
		&model.RecoveryFlow{},

		&model.Session{},
		&model.TokenExchangeCode{},

		&model.Identity{},
		&model.Credential{},
		&model.CredentialIdentifier{},

		&model.VerifiableAddress{},
		&model.VerificationCode{},
		&model.TokenExchangeCode{},

		&model.AlertDefinition{},
		&model.Alert{},
		&model.AlertEvent{},
	)
}
