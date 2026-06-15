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

package channel

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	channelDomain "github.com/muixstudio/clio/internal/channel"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type createChannelRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=128"`
	Description *string   `json:"description"`
	Visibility  string    `json:"visibility" validate:"required,oneof=public private"`
	CreatedBy   uuid.UUID `json:"created_by" validate:"required"`
	IsDefault   *bool     `json:"is_default"`
}

type updateChannelRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=128"`
	Description *string `json:"description"`
	Visibility  *string `json:"visibility" validate:"omitempty,oneof=public private"`
	Enabled     *bool   `json:"enabled"`
	IsDefault   *bool   `json:"is_default"`
}

type listChannelRequest struct {
	UserID   uuid.UUID  `form:"user_id" validate:"required"`
	Page     *int       `form:"page" validate:"omitempty,gte=1"`
	PageSize *int       `form:"page_size" validate:"omitempty,gte=1,lte=100"`
	StartAt  *time.Time `form:"start_at" time_format:"2006-01-02T15:04:05Z07:00"`
	EndAt    *time.Time `form:"end_at" time_format:"2006-01-02T15:04:05Z07:00"`
}

type channelItem struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	TeamID      uuid.UUID `json:"team_id"`
	Name        string    `json:"name"`
	Visibility  string    `json:"visibility"`
	Enabled     bool      `json:"enabled"`
	IsDefault   bool      `json:"is_default"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type channelListResponse struct {
	Count    int           `json:"count"`
	Total    int           `json:"total"`
	Channels []channelItem `json:"channels"`
}

func parseUUIDParam(c *gin.Context, param string, label string) (uuid.UUID, bool) {
	raw := c.Param(param)
	id, err := uuid.Parse(raw)
	if err != nil {
		response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument "+label+": "+raw))
		return uuid.Nil, false
	}
	return id, true
}

func toChannelItem(ch *channelDomain.Channel) channelItem {
	return channelItem{
		ID:          ch.ID,
		Description: ch.Description,
		TeamID:      ch.TeamID,
		Name:        ch.Name,
		Visibility:  string(ch.Visibility),
		Enabled:     ch.Enabled,
		IsDefault:   ch.IsDefault,
		CreatedBy:   ch.CreatedBy,
		CreatedAt:   ch.CreatedAt,
		UpdatedAt:   ch.UpdatedAt,
	}
}

// List handles GET /team/:team_id/channels.
func (h *ChannelHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}

		var req listChannelRequest
		if err := c.ShouldBindWith(&req, binding.Query); err != nil {
			log.Warn("list channels: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		options := channelDomain.ListOptions{
			Page:     1,
			PageSize: 20,
			TeamID:   &teamID,
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
		}
		if req.Page != nil {
			options.Page = *req.Page
		}
		if req.PageSize != nil {
			options.PageSize = *req.PageSize
		}

		records, total, err := h.d.ChannelPersister().ListChannels(c.Request.Context(), teamID, req.UserID, &options)
		if err != nil {
			log.Error("list channels: failed", zap.Stringer("team_id", teamID), zap.Stringer("user_id", req.UserID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		items := make([]channelItem, 0, len(records))
		for _, record := range records {
			items = append(items, toChannelItem(record))
		}
		response.SuccessWithData(c, channelListResponse{
			Count:    len(items),
			Total:    total,
			Channels: items,
		})
	}
}

// Get handles GET /team/:team_id/channel/:channel_id.
func (h *ChannelHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}
		channelID, ok := parseUUIDParam(c, "channel_id", "channel id")
		if !ok {
			return
		}

		record, err := h.d.ChannelPersister().GetChannelByID(c.Request.Context(), teamID, channelID)
		if err != nil {
			log.Error("get channel: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, toChannelItem(record))
	}
}

// Create handles POST /team/:team_id/channel.
func (h *ChannelHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}

		var req createChannelRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("create channel: invalid request", zap.Stringer("team_id", teamID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		var opts []channelDomain.Option
		if req.Description != nil {
			opts = append(opts, channelDomain.WithDescription(*req.Description))
		}

		var ch *channelDomain.Channel
		if req.IsDefault != nil && *req.IsDefault {
			ch = channelDomain.NewDefaultChannel(teamID, req.CreatedBy, opts...)
		} else {
			switch channelDomain.VisibilityType(req.Visibility) {
			case channelDomain.Public:
				ch = channelDomain.NewPublicChannel(teamID, req.CreatedBy, req.Name, opts...)
			case channelDomain.Private:
				ch = channelDomain.NewPrivateChannel(teamID, req.CreatedBy, req.Name, opts...)
			default:
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "visibility must be public or private"))
				return
			}
		}

		if err := h.d.ChannelPersister().AddChannel(c.Request.Context(), ch); err != nil {
			log.Error("create channel: failed", zap.Stringer("team_id", teamID), zap.String("name", req.Name), zap.Error(err))
			response.Fail(c, err)
			return
		}

		if err := h.d.ChannelMemberPersister().AddChannelMember(c.Request.Context(), teamID, ch.ID, req.CreatedBy, channelDomain.Owner); err != nil {
			log.Error("create channel owner: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", ch.ID), zap.Stringer("user_id", req.CreatedBy), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel created", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", ch.ID), zap.String("name", ch.Name))
		response.SuccessWithData(c, toChannelItem(ch))
	}
}

// Update handles PATCH /team/:team_id/channel/:channel_id.
func (h *ChannelHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}
		channelID, ok := parseUUIDParam(c, "channel_id", "channel id")
		if !ok {
			return
		}

		var req updateChannelRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("update channel: invalid request", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		var visibility *channelDomain.VisibilityType
		if req.Visibility != nil {
			v := channelDomain.VisibilityType(*req.Visibility)
			visibility = &v
		}
		updates := channelDomain.Update{
			Name:        req.Name,
			Description: req.Description,
			Visibility:  visibility,
			Enabled:     req.Enabled,
			IsDefault:   req.IsDefault,
		}

		if err := h.d.ChannelPersister().UpdateChannel(c.Request.Context(), teamID, channelID, &updates); err != nil {
			log.Error("update channel: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel updated", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID))
		response.SuccessOK(c)
	}
}

// Delete handles DELETE /team/:team_id/channel/:channel_id.
func (h *ChannelHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		teamID, ok := parseUUIDParam(c, "team_id", "team id")
		if !ok {
			return
		}
		channelID, ok := parseUUIDParam(c, "channel_id", "channel id")
		if !ok {
			return
		}

		if err := h.d.ChannelPersister().DeleteChannelByID(c.Request.Context(), teamID, channelID); err != nil {
			log.Error("delete channel: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel deleted", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID))
		response.SuccessOK(c)
	}
}
