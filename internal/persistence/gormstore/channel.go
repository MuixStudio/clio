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

package gormstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/channel"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/persistence/gormstore/model"
	"gorm.io/gorm"
)

func toChannelModel(c *channel.Channel) *model.Channel {
	return &model.Channel{
		ID:          c.ID,
		TeamID:      c.TeamID,
		Name:        c.Name,
		Description: c.Description,
		Visibility:  string(c.Visibility),
		Enabled:     c.Enabled,
		IsDefault:   c.IsDefault,
		CreatedBy:   c.CreatedBy,
	}
}

func fromChannelModel(m *model.Channel) *channel.Channel {
	return &channel.Channel{
		ID:          m.ID,
		Description: m.Description,
		TeamID:      m.TeamID,
		Name:        m.Name,
		Visibility:  channel.VisibilityType(m.Visibility),
		Enabled:     m.Enabled,
		IsDefault:   m.IsDefault,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toChannelMemberModel(teamID uuid.UUID, channelID uuid.UUID, userID uuid.UUID, role channel.MemberRole) *model.ChannelMember {
	return &model.ChannelMember{
		ID:        uuid.New(),
		TeamID:    teamID,
		ChannelID: channelID,
		UserID:    userID,
		Role:      string(role),
	}
}

func fromChannelMemberModel(m *model.ChannelMember) *channel.ChannelMember {
	return &channel.ChannelMember{
		ID:        m.ID,
		TeamID:    m.TeamID,
		ChannelID: m.ChannelID,
		UserID:    m.UserID,
		Role:      channel.MemberRole(m.Role),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (gs *GormStore) AddChannel(ctx context.Context, c *channel.Channel) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	if c.ID == (uuid.UUID{}) {
		c.ID = uuid.New()
	}
	m := toChannelModel(c)
	m.OrganizationID = orgID
	result := gs.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return nil
}

func (gs *GormStore) GetChannelByID(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID) (*channel.Channel, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.Channel
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              channelID,
		"team_id":         teamID,
		"organization_id": orgID,
	}).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("RECORD_NOT_FOUND", "channel not found: "+channelID.String())
		}
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return fromChannelModel(&m), nil
}

func (gs *GormStore) ListChannels(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, opts *channel.ListOptions) ([]*channel.Channel, int, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := gs.channelListQuery(ctx, orgID, opts).
		Where("team_id = ?", teamID).
		Where("visibility = ? OR EXISTS (?)",
			string(channel.Public),
			gs.db.WithContext(ctx).Model(&model.ChannelMember{}).
				Select("1").
				Where("channel_members.channel_id = channels.id").
				Where(map[string]any{
					"channel_members.organization_id": orgID,
					"channel_members.team_id":         teamID,
					"channel_members.user_id":         userID,
				}),
		)
	return gs.listChannels(query, opts)
}

func (gs *GormStore) UpdateChannel(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID, update *channel.Update) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	updates := map[string]any{}
	if update != nil {
		if update.Name != nil {
			updates["name"] = *update.Name
		}
		if update.Description != nil {
			updates["description"] = *update.Description
		}
		if update.Visibility != nil {
			updates["visibility"] = string(*update.Visibility)
		}
		if update.Enabled != nil {
			updates["enabled"] = *update.Enabled
		}
		if update.IsDefault != nil {
			updates["is_default"] = *update.IsDefault
		}
	}
	if len(updates) == 0 {
		return nil
	}

	result := gs.db.WithContext(ctx).Model(&model.Channel{}).
		Where(map[string]any{
			"id":              channelID,
			"team_id":         teamID,
			"organization_id": orgID,
		}).
		Updates(updates)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "channel not found: "+channelID.String())
	}
	return nil
}

func (gs *GormStore) DeleteChannelByID(ctx context.Context, teamID uuid.UUID, channelID uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              channelID,
		"team_id":         teamID,
		"organization_id": orgID,
	}).Delete(&model.Channel{})
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "channel not found: "+channelID.String())
	}
	return nil
}

func (gs *GormStore) AddChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID, role channel.MemberRole) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}
	if role == "" {
		role = channel.Member
	}

	m := toChannelMemberModel(teamID, channelID, userID, role)
	m.OrganizationID = orgID
	result := gs.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return nil
}

func (gs *GormStore) GetChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID) (*channel.ChannelMember, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var m model.ChannelMember
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"team_id":         teamID,
		"channel_id":      channelID,
		"user_id":         userID,
		"organization_id": orgID,
	}).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("RECORD_NOT_FOUND", "channel member not found: "+userID.String())
		}
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return fromChannelMemberModel(&m), nil
}

func (gs *GormStore) ListChannelMembers(ctx context.Context, teamID, channelID uuid.UUID) ([]*channel.ChannelMember, int, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := gs.db.WithContext(ctx).Model(&model.ChannelMember{}).Where(map[string]any{
		"team_id":         teamID,
		"channel_id":      channelID,
		"organization_id": orgID,
	})

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	var models []model.ChannelMember
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	members := make([]*channel.ChannelMember, 0, len(models))
	for i := range models {
		members = append(members, fromChannelMemberModel(&models[i]))
	}
	return members, int(count), nil
}

func (gs *GormStore) UpdateChannelMemberRole(ctx context.Context, teamID, channelID, userID uuid.UUID, role channel.MemberRole) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).Model(&model.ChannelMember{}).
		Where(map[string]any{
			"team_id":         teamID,
			"channel_id":      channelID,
			"user_id":         userID,
			"organization_id": orgID,
		}).
		Update("role", string(role))
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "channel member not found: "+userID.String())
	}
	return nil
}

func (gs *GormStore) RemoveChannelMember(ctx context.Context, teamID, channelID, userID uuid.UUID) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	result := gs.db.WithContext(ctx).Where(map[string]any{
		"team_id":         teamID,
		"channel_id":      channelID,
		"user_id":         userID,
		"organization_id": orgID,
	}).Delete(&model.ChannelMember{})
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "channel member not found: "+userID.String())
	}
	return nil
}

func (gs *GormStore) channelListQuery(ctx context.Context, orgID uuid.NullUUID, opts *channel.ListOptions) *gorm.DB {
	query := gs.db.WithContext(ctx).Model(&model.Channel{}).Where(map[string]any{
		"organization_id": orgID,
	})
	if opts == nil {
		return query
	}
	if opts.TeamID != nil {
		query = query.Where(map[string]any{"team_id": *opts.TeamID})
	}
	if opts.StartAt != nil {
		query = query.Where("created_at >= ?", *opts.StartAt)
	}
	if opts.EndAt != nil {
		query = query.Where("created_at <= ?", *opts.EndAt)
	}
	return query
}

func (gs *GormStore) listChannels(query *gorm.DB, opts *channel.ListOptions) ([]*channel.Channel, int, error) {
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}

	page, pageSize := 0, 0
	if opts != nil {
		page = opts.Page
		pageSize = opts.PageSize
	}
	if pageSize > 0 {
		query = query.Limit(pageSize)
		if page > 0 {
			query = query.Offset((page - 1) * pageSize)
		}
	}

	var models []model.Channel
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	channels := make([]*channel.Channel, 0, len(models))
	for i := range models {
		channels = append(channels, fromChannelModel(&models[i]))
	}
	return channels, int(count), nil
}
