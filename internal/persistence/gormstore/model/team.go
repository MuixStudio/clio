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

type Team struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	Name           string        `gorm:"column:name;not null;type:varchar(255);comment:team name"`
	Description    string        `gorm:"column:description;type:text;comment:team description"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (Team) TableName() string { return "teams" }

//type teamRoleModel struct {
//	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
//	TeamID         uuid.UUID     `gorm:"not null;index;type:varchar(36)"`
//	Name           string        `gorm:"not null;type:varchar(128)"`
//	Description    string        `gorm:"type:text"`
//	OrganizationID uuid.NullUUID `gorm:"index;type:uuid;default:null"`
//	IsSystemRole   bool          `gorm:"not null;default:false"`
//	CreatedAt      time.Time
//	UpdatedAt      time.Time
//}
//
//func (teamRoleModel) TableName() string { return "team_roles" }

type TeamMember struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	TeamID         uuid.UUID     `gorm:"column:team_id;not null;index;type:varchar(36);comment:FK to teams"`
	UserID         uuid.UUID     `gorm:"column:user_id;not null;type:varchar(36);comment:FK to users"`
	Email          string        `gorm:"column:email;not null;type:varchar(255);comment:member email"`
	RoleID         uuid.UUID     `gorm:"column:role_id;type:varchar(36);comment:FK to team roles"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	JoinedAt       time.Time     `gorm:"column:joined_at;comment:team join timestamp"`
	CreatedAt      time.Time     `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;comment:last update timestamp"`
}

func (TeamMember) TableName() string { return "team_members" }
