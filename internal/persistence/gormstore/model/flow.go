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

type LoginFlow struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	State          string        `gorm:"column:state;type:varchar(64);not null;comment:flow state"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:flow expiry time"`
	Active         string        `gorm:"column:active;type:varchar(64);comment:active credentials type"`
	Refresh        bool          `gorm:"column:refresh;not null;default:false;comment:is refresh login flow"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (LoginFlow) TableName() string { return "login_flows" }

type RegistrationFlow struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	State          string        `gorm:"column:state;type:varchar(64);not null;comment:flow state"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:flow expiry time"`
	Active         string        `gorm:"column:active;type:varchar(64);comment:active credentials type"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (RegistrationFlow) TableName() string { return "registration_flows" }

type LogoutFlow struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	IdentityID     uuid.UUID     `gorm:"column:identity_id;type:uuid;not null;index;comment:FK to identities"`
	State          string        `gorm:"column:state;type:varchar(64);not null;comment:flow state"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:flow expiry time"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (LogoutFlow) TableName() string { return "logout_flows" }

type RecoveryFlow struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	State          string        `gorm:"column:state;type:varchar(64);not null;comment:flow state"`
	Email          string        `gorm:"column:email;type:varchar(255);comment:recovery email address"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:flow expiry time"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (RecoveryFlow) TableName() string { return "recovery_flows" }

type VerificationFlow struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	State          string        `gorm:"column:state;type:varchar(64);not null;comment:flow state"`
	Email          string        `gorm:"column:email;type:varchar(255);comment:verification email address"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:flow expiry time"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (VerificationFlow) TableName() string { return "verification_flows" }
