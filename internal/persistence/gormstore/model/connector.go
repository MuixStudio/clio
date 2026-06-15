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
	"gorm.io/datatypes"
)

type Connector struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	Type           string         `gorm:"column:type;not null;type:varchar(64);comment:connector type"`
	Name           string         `gorm:"column:name;not null;type:varchar(255);comment:connector name"`
	Token          string         `gorm:"column:token;not null;type:varchar(36);comment:connector token"`
	TeamID         uuid.UUID      `gorm:"column:team_id;index;type:varchar(36);comment:FK to teams"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	Labels         datatypes.JSON `gorm:"column:labels;type:text;comment:connector labels as JSON"`
	Enabled        bool           `gorm:"column:enabled;not null;default:true;comment:connector is enabled"`
	CreatedAt      time.Time      `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;comment:last update timestamp"`
}

func (Connector) TableName() string { return "connectors" }
