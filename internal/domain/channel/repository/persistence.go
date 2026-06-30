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

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/domain/channel/entity"
)

type Update struct {
	Name        *string
	Description *string
	Visibility  *entity.VisibilityType
	Enabled     *bool
	IsDefault   *bool
}

type ListOptions struct {
	Page, PageSize int
	TeamID         *uuid.UUID
	StartAt        *time.Time // 筛选时间段的开始时间
	EndAt          *time.Time // 筛选时间段的结束时间
}

type (
	ChannelPersister interface {
		AddChannel(ctx context.Context, channel *entity.Channel) error
		GetChannelByID(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID) (*entity.Channel, error)
		ListChannels(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, opts *ListOptions) (channels []*entity.Channel, total int, err error)
		UpdateChannel(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID, updates *Update) error
		DeleteChannelByID(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID) error
	}
	ChannelPersisterProvider interface {
		ChannelPersister() ChannelPersister
	}

	ChannelMemberPersister interface {
		AddChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID, role entity.MemberRole) error
		GetChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID) (*entity.ChannelMember, error)
		ListChannelMembers(ctx context.Context, teamID, channelID uuid.UUID) (members []*entity.ChannelMember, total int, err error)
		UpdateChannelMemberRole(ctx context.Context, teamID, channelID, userID uuid.UUID, role entity.MemberRole) error
		RemoveChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID) error
	}
	ChannelMemberPersisterProvider interface {
		ChannelMemberPersister() ChannelMemberPersister
	}
)
