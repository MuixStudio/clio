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

package alert

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	alertEntity "github.com/muixstudio/clio/internal/domain/alert/entity"
	alertRepo "github.com/muixstudio/clio/internal/domain/alert/repository"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type listAlertRequest struct {
	Page        *int       `form:"page" validate:"omitempty,gte=1"`
	PageSize    *int       `form:"page_size" validate:"omitempty,gte=1,lte=100"`
	Status      *string    `form:"status" validate:"omitempty,oneof=firing resolved"`
	Severity    *string    `form:"severity" validate:"omitempty,oneof=critical high medium low"`
	ConnectorID *string    `form:"connector_id" validate:"omitempty,uuid"`
	StartAt     *time.Time `form:"start_at" time_format:"2006-01-02T15:04:05Z07:00"`
	EndAt       *time.Time `form:"end_at" time_format:"2006-01-02T15:04:05Z07:00"`
}

type alertItem struct {
	ID          uuid.UUID         `json:"id"`
	Fingerprint string            `json:"fingerprint"`
	ConnectorID uuid.UUID         `json:"connector_id"`
	TeamID      uuid.UUID         `json:"team_id"`
	Severity    string            `json:"severity"`
	Status      string            `json:"status"`
	StartsAt    time.Time         `json:"starts_at"`
	EndsAt      time.Time         `json:"ends_at"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type alertResponse struct {
	Count  int         `json:"count"`
	Total  int         `json:"total"`
	Alerts []alertItem `json:"alerts"`
}

func toAlertItem(a *alertEntity.NormalizedAlert) alertItem {
	var endsAt time.Time
	if a.EndsAt != nil {
		endsAt = *a.EndsAt
	}

	return alertItem{
		ID:          a.ID,
		Fingerprint: a.Fingerprint,
		ConnectorID: a.ConnectorID,
		TeamID:      a.TeamID,
		Severity:    string(a.Severity),
		Status:      string(a.Status),
		StartsAt:    a.StartsAt,
		EndsAt:      endsAt,
		Labels:      a.Labels,
		Annotations: a.Annotations,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// List handles GET /team/:team_id/alerts?status=&severity=&connector_id=&start_at=&end_at=&page=&page_size=.
func (h *AlertHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		var req listAlertRequest
		if err := c.ShouldBindWith(&req, binding.Query); err != nil {
			log.Warn("list alerts: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		options := alertRepo.ListOptions{
			Page:     1,
			PageSize: 20,
			TeamID:   &teamID,
		}
		if req.Page != nil {
			options.Page = *req.Page
		}
		if req.PageSize != nil {
			options.PageSize = *req.PageSize
		}
		if req.Status != nil {
			status := alertEntity.Status(*req.Status)
			options.Status = &status
		}
		if req.Severity != nil {
			severity := alertEntity.Severity(*req.Severity)
			options.Severity = &severity
		}
		if req.ConnectorID != nil {
			connectorID, err := uuid.Parse(*req.ConnectorID)
			if err != nil {
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector_id: "+*req.ConnectorID))
				return
			}
			options.ConnectorID = &connectorID
		}

		options.StartAt = req.StartAt
		options.EndAt = req.EndAt

		records, total, err := h.d.AlertPersister().ListAlert(c.Request.Context(), &options)
		if err != nil {
			log.Error("list alerts: query failed", zap.Error(err))
			response.Fail(c, err)
			return
		}

		alerts := make([]alertItem, 0, len(records))
		for _, a := range records {
			alerts = append(alerts, toAlertItem(a))
		}
		log.Debug("list alerts", zap.Int("count", total))
		response.SuccessWithData(c, alertResponse{
			Count:  len(alerts),
			Total:  total,
			Alerts: alerts,
		})
	}
}
