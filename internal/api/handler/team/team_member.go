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

package team

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/orgctx"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type addMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type memberItem struct {
	ID       uuid.UUID `json:"id"`
	TeamID   uuid.UUID `json:"team_id"`
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	JoinedAt time.Time `json:"joined_at"`
}

type memberResponse struct {
	Count   int          `json:"count"`
	Members []memberItem `json:"members"`
}

// ListMembers handles GET /team/:team_id/members.
func (h *TeamHandler) ListMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		orgID, ok := orgctx.OrgIDFromCtx(c.Request.Context())
		if !ok {
			response.Fail(c, errors.Unauthorized("MISSING_ORG", "organization context is required"))
			return
		}

		records, err := h.d.TeamMemberPersister().ListMembers(c.Request.Context(), orgID, teamID)
		if err != nil {
			log.Error("list members: failed", zap.String("team_id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		members := make([]memberItem, 0, len(records))
		for _, m := range records {
			members = append(members, memberItem{
				ID:       m.ID,
				TeamID:   m.TeamID,
				UserID:   m.UserID,
				Email:    m.Email,
				JoinedAt: m.JoinedAt,
			})
		}
		response.SuccessWithData(c, memberResponse{
			Count:   len(members),
			Members: members,
		})
	}
}

// AddMember handles POST /team/:team_id/member.
func (h *TeamHandler) AddMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		var req addMemberRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("add member: invalid request", zap.String("team_id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		orgID, ok := orgctx.OrgIDFromCtx(c.Request.Context())
		if !ok {
			response.Fail(c, errors.Unauthorized("MISSING_ORG", "organization context is required"))
			return
		}

		if err := h.d.TeamMemberPersister().AddMember(c.Request.Context(), orgID, teamID, req.UserID); err != nil {
			log.Error("add member: failed", zap.String("team_id", tid), zap.Stringer("user_id", req.UserID), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("member added", zap.String("team_id", tid), zap.Stringer("user_id", req.UserID))
		response.SuccessOK(c)
	}
}

// RemoveMember handles DELETE /team/:team_id/member/:member_id.
func (h *TeamHandler) RemoveMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		mid := c.Param("member_id")
		memberID, err := uuid.Parse(mid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument member id: "+mid))
			return
		}

		orgID, ok := orgctx.OrgIDFromCtx(c.Request.Context())
		if !ok {
			response.Fail(c, errors.Unauthorized("MISSING_ORG", "organization context is required"))
			return
		}

		if err := h.d.TeamMemberPersister().RemoveMember(c.Request.Context(), orgID, memberID); err != nil {
			log.Error("remove member: failed", zap.String("team_id", tid), zap.String("member_id", mid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("member removed", zap.String("team_id", tid), zap.String("member_id", mid))
		response.SuccessOK(c)
	}
}
