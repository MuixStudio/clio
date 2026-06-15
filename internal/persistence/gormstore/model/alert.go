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

type AlertModel struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key;not null"`
	Fingerprint    string         `gorm:"column:fingerprint;uniqueIndex:idx_alert_fingerprint;type:varchar(128);comment:alert fingerprint"`
	Source         string         `gorm:"column:source;type:varchar(64);comment:alert source"`
	ConnectorID    uuid.UUID      `gorm:"column:connector_id;uniqueIndex:idx_alert_fingerprint;type:varchar(36);comment:FK to connectors"`
	TeamID         uuid.UUID      `gorm:"column:team_id;index;type:varchar(36);comment:FK to teams"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	Severity       string         `gorm:"column:severity;type:varchar(32);comment:alert severity"`
	Status         string         `gorm:"column:status;type:varchar(32);comment:alert status"`
	StartsAt       time.Time      `gorm:"column:starts_at;comment:alert start time"`
	EndsAt         time.Time      `gorm:"column:ends_at;comment:alert end time"`
	EntityName     string         `gorm:"column:entity_name;type:varchar(255);comment:alert entity name"`
	EntityIP       string         `gorm:"column:entity_ip;type:varchar(64);comment:alert entity IP"`
	EntitySvc      string         `gorm:"column:entity_svc;type:varchar(255);comment:alert entity service"`
	Labels         datatypes.JSON `gorm:"column:labels;type:text;comment:alert labels as JSON"`
	Annotations    datatypes.JSON `gorm:"column:annotations;type:text;comment:alert annotations as JSON"`
	RawRef         string         `gorm:"column:raw_ref;type:text;comment:raw alert reference"`
	CreatedAt      time.Time      `gorm:"column:created_at;comment:creation timestamp"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;comment:last update timestamp"`
}

func (AlertModel) TableName() string { return "alerts" }

type AlertDefinition struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	Fingerprint    string         `gorm:"column:fingerprint;uniqueIndex:idx_alert_definition_fingerprint,option:NULLS NOT DISTINCT;type:varchar(128);comment:alert fingerprint, hash of labels"`
	AlertName      string         `gorm:"column:alert_name;type:varchar(128);not null;comment:alert rule name"`
	Labels         datatypes.JSON `gorm:"column:labels;type:jsonb;not null;comment:full label set, source of fingerprint"`
	Severity       string         `gorm:"column:severity;type:varchar(32);comment:critical / warning / info"`
	ConnectorID    uuid.UUID      `gorm:"column:connector_id;uniqueIndex:idx_alert_definition_fingerprint,option:NULLS NOT DISTINCT;type:uuid;comment:FK to connectors"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID         uuid.UUID      `gorm:"column:team_id;type:uuid;default:null;comment:FK to teams"`
	FirstSeenAt    time.Time      `gorm:"column:first_seen_at;not null;comment:first time this fingerprint was seen"`
	LastSeenAt     time.Time      `gorm:"column:last_seen_at;not null;comment:last time this fingerprint was seen, updated on every trigger"`
}

func (AlertDefinition) TableName() string { return "alert_definition" }

type Alert struct {
	ID                uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	AlertDefinitionID uuid.UUID      `gorm:"column:alert_definition_id;type:uuid;not null;index:idx_alerts_alert_id;comment:FK to alerts"`
	Fingerprint       string         `gorm:"column:fingerprint;type:varchar(128);not null;uniqueIndex:udx_alerts_fingerprint_starts_at_connector_id,option:NULLS NOT DISTINCT;index:idx_alerts_fingerprint;comment:FK to alert_definitions"`
	Status            string         `gorm:"column:status;type:varchar(32);not null;comment:firing / resolved"`
	Severity          string         `gorm:"column:severity;type:varchar(32);comment:critical / warning / info"`
	Labels            datatypes.JSON `gorm:"column:labels;type:jsonb;comment:labels"`
	Annotations       datatypes.JSON `gorm:"column:annotations;type:jsonb;comment:human-readable annotations"`
	StartsAt          time.Time      `gorm:"column:starts_at;not null;uniqueIndex:udx_alerts_fingerprint_starts_at_connector_id,option:NULLS NOT DISTINCT;comment:time when this firing began"`
	EndsAt            *time.Time     `gorm:"column:ends_at;default:null;comment:time when resolved, NULL means still firing"`
	FlapCount         int            `gorm:"column:flap_count;default:0;comment:number of firing<->resolved transitions"`
	RepeatCount       int            `gorm:"column:repeat_count;default:0;comment:number of duplicate firing pushes for this trigger"`
	ConnectorID       uuid.UUID      `gorm:"column:connector_id;uniqueIndex:udx_alerts_fingerprint_starts_at_connector_id,option:NULLS NOT DISTINCT;type:uuid;comment:FK to connectors"`
	OrganizationID    uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID            uuid.UUID      `gorm:"column:team_id;type:uuid;default:null;index:idx_alerts_team_id;comment:FK to teams"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:record last update time"`
}

func (Alert) TableName() string { return "alerts" }

type AlertEventType string

const (
	AlertEventFiring     AlertEventType = "firing"
	AlertEventResolved   AlertEventType = "resolved"
	AlertEventFlapping   AlertEventType = "flapping"   // 短时间内反复横跳
	AlertEventSuppressed AlertEventType = "suppressed" // 被 silence / inhibit 压制
	AlertEventNotified   AlertEventType = "notified"   // 已发出通知（用于审计）
	AlertEventEscalated  AlertEventType = "escalated"  // 升级处理
)

type AlertEvent struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	AlertID        uuid.UUID      `gorm:"column:alert_id;type:uuid;not null;index:idx_alert_events_alert_id;comment:FK to alerts"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	EventType      AlertEventType `gorm:"column:event_type;type:varchar(32);not null;index:idx_alert_events_event_type;comment:firing / resolved / flapping / suppressed / notified / escalated"`
	OccurredAt     time.Time      `gorm:"column:occurred_at;not null;index:idx_alert_events_occurred_at;comment:time this event occurred"`
	RawPayload     datatypes.JSON `gorm:"column:raw_payload;type:jsonb;comment:original push payload, for audit and incident replay"`
}

func (AlertEvent) TableName() string { return "alert_events" }
