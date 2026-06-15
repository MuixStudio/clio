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

type Identity struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	Traits         string        `gorm:"column:traits;type:text;not null;comment:identity traits as JSON"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (Identity) TableName() string { return "identities" }

type Credential struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	IdentityID     uuid.UUID     `gorm:"column:identity_id;type:uuid;not null;index;comment:FK to identities"`
	Type           string        `gorm:"column:type;type:varchar(64);not null;comment:credential type"`
	Config         string        `gorm:"column:config;type:text;not null;comment:credential configuration as JSON"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (Credential) TableName() string { return "identity_credentials" }

type CredentialIdentifier struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;not null;comment:primary key"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;uniqueIndex:udx_cred_identifier,option:NULLS NOT DISTINCT;type:uuid;default:null;comment:FK to organizations"`
	CredentialID   uuid.UUID     `gorm:"column:credential_id;type:uuid;not null;index;comment:FK to identity_credentials"`
	Type           string        `gorm:"column:type;type:varchar(64);not null;uniqueIndex:udx_cred_identifier,option:NULLS NOT DISTINCT;comment:identifier type"`
	Identifier     string        `gorm:"column:identifier;type:varchar(255);not null;uniqueIndex:udx_cred_identifier,option:NULLS NOT DISTINCT;comment:identifier value"`
}

func (CredentialIdentifier) TableName() string { return "identity_credential_identifiers" }
