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

type VerificationCode struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	FlowID         uuid.UUID     `gorm:"column:flow_id;type:uuid;not null;index;comment:FK to verification or recovery flow"`
	FlowType       string        `gorm:"column:flow_type;type:varchar(64);not null;comment:verification code flow type"`
	Address        string        `gorm:"column:address;type:varchar(255);not null;comment:target address"`
	Code           string        `gorm:"column:code;type:varchar(10);not null;comment:verification code"`
	ExpiresAt      time.Time     `gorm:"column:expires_at;not null;comment:code expiry time"`
	UsedAt         *time.Time    `gorm:"column:used_at;comment:code usage timestamp"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
}

func (VerificationCode) TableName() string { return "verification_codes" }
