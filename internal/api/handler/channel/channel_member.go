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
	channelEntity "github.com/muixstudio/clio/internal/domain/channel/entity"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type addChannelMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Role   *string   `json:"role" validate:"omitempty,oneof=owner admin member"`
}

type updateChannelMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=owner admin member"`
}

type channelMemberItem struct {
	ID        uuid.UUID `json:"id"`
	TeamID    uuid.UUID `json:"team_id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type channelMemberListResponse struct {
	Count   int                 `json:"count"`
	Total   int                 `json:"total"`
	Members []channelMemberItem `json:"members"`
}

func toChannelMemberItem(member *channelEntity.ChannelMember) channelMemberItem {
	return channelMemberItem{
		ID:        member.ID,
		TeamID:    member.TeamID,
		ChannelID: member.ChannelID,
		UserID:    member.UserID,
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

// ListMembers handles GET /team/:team_id/channel/:channel_id/members.
func (h *ChannelHandler) ListMembers() gin.HandlerFunc {
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

		records, total, err := h.d.ChannelMemberPersister().ListChannelMembers(c.Request.Context(), teamID, channelID)
		if err != nil {
			log.Error("list channel members: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		members := make([]channelMemberItem, 0, len(records))
		for _, record := range records {
			members = append(members, toChannelMemberItem(record))
		}
		response.SuccessWithData(c, channelMemberListResponse{
			Count:   len(members),
			Total:   total,
			Members: members,
		})
	}
}

// AddMember handles POST /team/:team_id/channel/:channel_id/member.
func (h *ChannelHandler) AddMember() gin.HandlerFunc {
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

		var req addChannelMemberRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("add channel member: invalid request", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		role := channelEntity.Member
		if req.Role != nil {
			role = channelEntity.MemberRole(*req.Role)
		}
		if err := h.d.ChannelMemberPersister().AddChannelMember(c.Request.Context(), teamID, channelID, req.UserID, role); err != nil {
			log.Error("add channel member: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", req.UserID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel member added", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", req.UserID), zap.String("role", string(role)))
		response.SuccessOK(c)
	}
}

// GetMember handles GET /team/:team_id/channel/:channel_id/member/:user_id.
func (h *ChannelHandler) GetMember() gin.HandlerFunc {
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
		userID, ok := parseUUIDParam(c, "user_id", "user id")
		if !ok {
			return
		}

		member, err := h.d.ChannelMemberPersister().GetChannelMember(c.Request.Context(), teamID, channelID, userID)
		if err != nil {
			log.Error("get channel member: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, toChannelMemberItem(member))
	}
}

// UpdateMemberRole handles PATCH /team/:team_id/channel/:channel_id/member/:user_id.
func (h *ChannelHandler) UpdateMemberRole() gin.HandlerFunc {
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
		userID, ok := parseUUIDParam(c, "user_id", "user id")
		if !ok {
			return
		}

		var req updateChannelMemberRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("update channel member role: invalid request", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		role := channelEntity.MemberRole(req.Role)
		if err := h.d.ChannelMemberPersister().UpdateChannelMemberRole(c.Request.Context(), teamID, channelID, userID, role); err != nil {
			log.Error("update channel member role: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID), zap.String("role", req.Role), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel member role updated", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID), zap.String("role", req.Role))
		response.SuccessOK(c)
	}
}

// RemoveMember handles DELETE /team/:team_id/channel/:channel_id/member/:user_id.
func (h *ChannelHandler) RemoveMember() gin.HandlerFunc {
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
		userID, ok := parseUUIDParam(c, "user_id", "user id")
		if !ok {
			return
		}

		if err := h.d.ChannelMemberPersister().RemoveChannelMember(c.Request.Context(), teamID, channelID, userID); err != nil {
			log.Error("remove channel member: failed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("channel member removed", zap.Stringer("team_id", teamID), zap.Stringer("channel_id", channelID), zap.Stringer("user_id", userID))
		response.SuccessOK(c)
	}
}
