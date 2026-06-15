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

type VerifiableAddress struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;uniqueIndex:udx_via_value,option:NULLS NOT DISTINCT;type:uuid;default:null;comment:FK to organizations"`
	IdentityID     uuid.UUID     `gorm:"column:identity_id;type:uuid;not null;index;comment:FK to identities"`
	Value          string        `gorm:"column:value;type:varchar(255);not null;uniqueIndex:udx_via_value,option:NULLS NOT DISTINCT;comment:verifiable address value"`
	Via            string        `gorm:"column:via;type:varchar(16);not null;uniqueIndex:udx_via_value,option:NULLS NOT DISTINCT;comment:verification channel"`
	Verified       bool          `gorm:"column:verified;not null;default:false;comment:address is verified"`
	VerifiedAt     *time.Time    `gorm:"column:verified_at;comment:verification timestamp"`
	Status         string        `gorm:"column:status;type:varchar(32);not null;default:'pending';comment:verification status"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (VerifiableAddress) TableName() string { return "identity_verifiable_addresses" }
