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
	teamEntity "github.com/muixstudio/clio/internal/domain/team/entity"
	teamRepo "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type createTeamRequest struct {
	Name        string  `json:"name" validate:"required,min=4,max=32"`
	Description *string `json:"description"`
}

type updateTeamRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=4,max=32"`
	Description *string `json:"description"`
}

type listTeamRequest struct {
	Page     *int `form:"page" validate:"omitempty,gte=1"`
	PageSize *int `form:"page_size" validate:"omitempty,gte=1,lte=100"`
}

type teamItem struct {
	ID uuid.UUID `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type teamResponse struct {
	Count int         `json:"count"`
	Teams *[]teamItem `json:"teams"`
}

// Get handles GET /team/:team_id.
func (h *TeamHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		t, err := h.d.TeamPersister().GetTeam(c.Request.Context(), teamID)
		if err != nil {
			log.Error("get team: failed", zap.String("id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, teamItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
}

// Create handles POST /team.
func (h *TeamHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		var req createTeamRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("create team: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		t := &teamEntity.Team{Name: req.Name}
		if req.Description != nil {
			t.Description = *req.Description
		}

		if err := h.d.TeamPersister().CreateTeam(c.Request.Context(), t); err != nil {
			log.Error("create team: persist failed", zap.String("name", req.Name), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("team created", zap.Stringer("id", t.ID), zap.String("name", t.Name))
		response.SuccessWithData(c, teamItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
}

// Update handles PATCH /team/:team_id.
func (h *TeamHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		var req updateTeamRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("update team: invalid request", zap.String("id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		update := teamRepo.TeamUpdate{
			Name:        req.Name,
			Description: req.Description,
		}
		if err := h.d.TeamPersister().UpdateTeam(c.Request.Context(), teamID, update); err != nil {
			log.Error("update team: failed", zap.String("id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("team updated", zap.String("id", tid))
		response.SuccessOK(c)
	}
}

// Delete handles DELETE /team/:team_id.
func (h *TeamHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		if err := h.d.TeamPersister().DeleteTeam(c.Request.Context(), teamID); err != nil {
			log.Error("delete team: failed", zap.String("id", tid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("team deleted", zap.String("id", tid))
		response.SuccessNoContent(c)
	}
}

// GetUserTeams handles GET /user/:user_id/teams.
func (h *TeamHandler) GetUserTeams() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		uid := c.Param("user_id")
		userID, err := uuid.Parse(uid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument user id: "+uid))
			return
		}

		records, err := h.d.TeamPersister().ListTeamsByUserID(c.Request.Context(), userID)
		if err != nil {
			log.Error("get user teams: failed", zap.String("user_id", uid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		teams := make([]teamItem, 0, len(records))
		for _, t := range records {
			teams = append(teams, teamItem{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   t.UpdatedAt,
			})
		}
		log.Debug("get user teams", zap.String("user_id", uid), zap.Int("count", len(teams)))
		response.SuccessWithData(c, teamResponse{
			Count: len(teams),
			Teams: &teams,
		})
	}
}

// List handles GET /teams?page=1&items_per_page=20.
func (h *TeamHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		var req listTeamRequest
		if err := c.ShouldBindWith(&req, binding.Query); err != nil {
			log.Warn("list teams: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		page, pageSize := 1, 20
		if req.Page != nil {
			page = *req.Page
		}
		if req.PageSize != nil {
			pageSize = *req.PageSize
		}

		records, err := h.d.TeamPersister().ListTeams(c.Request.Context())
		if err != nil {
			log.Error("list teams: query failed", zap.Error(err))
			response.Fail(c, err)
			return
		}

		start := (page - 1) * pageSize
		if start > len(records) {
			start = len(records)
		}
		end := start + pageSize
		if end > len(records) {
			end = len(records)
		}
		pageRecords := records[start:end]

		teams := make([]teamItem, 0, len(pageRecords))
		for _, t := range pageRecords {
			teams = append(teams, teamItem{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   t.UpdatedAt,
			})
		}
		log.Debug("list teams", zap.Int("count", len(teams)))
		response.SuccessWithData(c, teamResponse{
			Count: len(records),
			Teams: &teams,
		})
	}
}
