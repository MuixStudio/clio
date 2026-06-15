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

// Channel 表示一个告警通知渠道。
//
// 一条告警最终会通过某个 Channel 发出去，例如 webhook、email、Slack、飞书、钉钉等。
// Visibility 控制频道可见范围：public 表示 team 内所有人可见，private 表示只有 ChannelMember 成员可见。
// IsDefault 标记该频道是否为 team 默认频道，可用于没有命中任何路由规则时的兜底投递。
// Config 保存不同渠道各自需要的配置，比如 webhook URL、请求头、收件人、模板参数等。
type Channel struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID         uuid.UUID     `gorm:"column:team_id;type:uuid;not null;index:idx_channels_team_id;comment:FK to teams"`
	Name           string        `gorm:"column:name;type:varchar(128);not null;comment:channel name"`
	Description    string        `gorm:"column:description;type:text;comment:channel description"`
	Visibility     string        `gorm:"column:visibility;type:varchar(16);not null;default:'private';index:idx_channels_visibility;comment:public or private"`
	Enabled        bool          `gorm:"column:enabled;not null;default:true;comment:channel is enabled"`
	IsDefault      bool          `gorm:"column:is_default;not null;default:false;index:idx_channels_is_default;comment:default channel for team"`
	CreatedBy      uuid.UUID     `gorm:"column:created_by;type:uuid;not null;index:idx_channels_created_by;comment:channel creator"`
	CreatedAt      time.Time     `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;autoUpdateTime;comment:record last update time"`
}

func (Channel) TableName() string { return "channels" }

// ChannelMember 表示私有频道的成员权限。
//
// public Channel 不需要依赖这张表，team 内所有人都能看到。
// private Channel 只有这张表中存在成员记录的用户能看到；Role 用于区分 owner / admin / member。
type ChannelMember struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID         uuid.UUID     `gorm:"column:team_id;type:uuid;not null;index:idx_channel_members_team_id;comment:FK to teams"`
	ChannelID      uuid.UUID     `gorm:"column:channel_id;type:uuid;not null;uniqueIndex:udx_channel_members_channel_user,option:NULLS NOT DISTINCT;index:idx_channel_members_channel_id;comment:FK to channels"`
	UserID         uuid.UUID     `gorm:"column:user_id;type:uuid;not null;uniqueIndex:udx_channel_members_channel_user,option:NULLS NOT DISTINCT;index:idx_channel_members_user_id;comment:FK to users"`
	Role           string        `gorm:"column:role;type:varchar(32);not null;default:'member';comment:member role in channel"`
	CreatedAt      time.Time     `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;autoUpdateTime;comment:record last update time"`
}

func (ChannelMember) TableName() string { return "channel_members" }

// AlertRoute 表示一条告警路由规则。
//
// 路由规则用于判断某条告警应该发送到哪些 Channel。
// Matchers 保存匹配条件，通常匹配告警 labels / severity / connector_id 等字段。
// Priority 控制匹配顺序，值越小越先匹配；Continue 表示命中本规则后是否继续匹配后续规则。
type AlertRoute struct {
	ID             uuid.UUID      `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	OrganizationID uuid.NullUUID  `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	TeamID         uuid.UUID      `gorm:"column:team_id;type:uuid;not null;index:idx_alert_routes_team_id;comment:FK to teams"`
	Name           string         `gorm:"column:name;type:varchar(128);not null;comment:route name"`
	Priority       int            `gorm:"column:priority;not null;default:0;index:idx_alert_routes_priority;comment:lower value matches first"`
	Matchers       datatypes.JSON `gorm:"column:matchers;type:jsonb;not null;comment:alert label matchers"`
	Enabled        bool           `gorm:"column:enabled;not null;default:true;comment:route is enabled"`
	Continue       bool           `gorm:"column:continue;not null;default:false;comment:continue matching following routes"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:record last update time"`
}

func (AlertRoute) TableName() string { return "alert_routes" }

// AlertRouteChannel 表示 AlertRoute 和 Channel 的关联关系。
//
// 一个路由规则可以绑定多个 Channel，同一个 Channel 也可以被多个路由规则复用。
// 这张表只负责保存多对多关系，实际路由匹配逻辑由 AlertRoute 决定。
type AlertRouteChannel struct {
	ID             uuid.UUID     `gorm:"primaryKey;column:id;type:uuid;comment:primary key"`
	OrganizationID uuid.NullUUID `gorm:"column:organization_id;index;type:uuid;default:null;comment:FK to organizations"`
	RouteID        uuid.UUID     `gorm:"column:route_id;type:uuid;not null;uniqueIndex:udx_alert_route_channels_route_channel,option:NULLS NOT DISTINCT;comment:FK to alert_routes"`
	ChannelID      uuid.UUID     `gorm:"column:channel_id;type:uuid;not null;uniqueIndex:udx_alert_route_channels_route_channel,option:NULLS NOT DISTINCT;index:idx_alert_route_channels_channel_id;comment:FK to channels"`
	CreatedAt      time.Time     `gorm:"column:created_at;autoCreateTime;comment:record creation time"`
}

func (AlertRouteChannel) TableName() string { return "alert_route_channels" }
