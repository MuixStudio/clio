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

package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                          uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID              uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	IdentityID                  uuid.UUID     `gorm:"column:identity_id;type:uuid;not null;index;comment:FK to identities"`
	RefreshToken                string        `gorm:"column:refresh_token;type:varchar(64);not null;uniqueIndex;comment:opaque refresh token"`
	ExpiresAt                   time.Time     `gorm:"column:expires_at;not null;comment:session expiry time"`
	IssuedAt                    time.Time     `gorm:"column:issued_at;not null;comment:session issue time"`
	AuthenticatedAt             time.Time     `gorm:"column:authenticated_at;not null;comment:authentication time"`
	Active                      bool          `gorm:"column:active;not null;default:true;comment:session is active"`
	AuthenticatorAssuranceLevel string        `gorm:"column:authenticator_assurance_level;type:varchar(16);not null;default:'';comment:AAL level"`
	AMR                         string        `gorm:"column:amr;type:text;not null;default:'[]';comment:authentication methods as JSON array"`
	CreatedAt                   time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt                   time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (Session) TableName() string { return "sessions" }

type TokenExchangeCode struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	Code           string        `gorm:"column:code;type:varchar(64);not null;uniqueIndex;comment:token exchange code"`
	SessionID      uuid.UUID     `gorm:"column:session_id;type:uuid;not null;index;comment:FK to sessions"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:code expiry time"`
	UsedAt         *time.Time    `gorm:"column:used_at;comment:code usage timestamp"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (TokenExchangeCode) TableName() string { return "session_token_exchange_codes" }
