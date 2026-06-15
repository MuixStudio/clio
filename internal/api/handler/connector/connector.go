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

package connector

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/connector"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

type createConnectorRequest struct {
	Name   string            `json:"name" validate:"required"`
	Type   string            `json:"type" validate:"required"`
	Labels map[string]string `json:"labels"`
}

type updateConnectorRequest struct {
	Name    *string `json:"name" validate:"omitempty,min=4,max=32"`
	Enabled *bool   `json:"enabled"`
}

type listConnectorRequest struct {
	Page     *int    `form:"page" validate:"omitempty,gte=1"`
	PageSize *int    `form:"page_size" validate:"omitempty,gte=1,lte=100"`
	Status   *string `form:"status" validate:"omitempty,oneof=enable disable"`
	Type     *string `form:"type" validate:"omitempty"`
}

type connectorResponse struct {
	ID uuid.UUID `json:"id"`

	Type    string            `json:"type"`
	Name    string            `json:"name"`
	Token   string            `json:"token,omitempty"`
	Labels  map[string]string `json:"labels"`
	Enabled bool              `json:"enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create handles POST /:team_id/connector.
func (h *ConnectorHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector team id: "+tid))
			return
		}

		var req createConnectorRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("create connector: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		_, ok := h.d.ConnectorProviders()[req.Type]
		if !ok {
			log.Warn("webhook receive: no provider for type", zap.String("type", req.Type))
			response.Fail(
				c,
				errors.BadRequest(
					"NO_PROVIDER",
					fmt.Sprintf("no provider for type %s", req.Type),
				),
			)
			return
		}

		conn := connector.NewConnector(
			req.Name,
			req.Type,
			connector.WithTeamID(teamID),
			connector.WithLabels(req.Labels),
		)

		if err := h.d.ConnectorPersister().CreateConnector(c.Request.Context(), conn); err != nil {
			log.Error("create connector: persist failed", zap.String("name", req.Name), zap.String("type", req.Type), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("connector created", zap.Stringer("id", conn.ID), zap.String("name", conn.Name), zap.String("type", conn.Type))
		response.SuccessWithData(c, connectorResponse{
			ID:      conn.ID,
			Type:    conn.Type,
			Name:    conn.Name,
			Token:   conn.Token,
			Enabled: conn.Enabled,
		})
	}
}

// Update handles PATCH /:team_id/connector/:connector_id.
func (h *ConnectorHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		cid := c.Param("connector_id")
		connectorID, err := uuid.Parse(cid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector id: "+cid))
			return
		}
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector team id: "+tid))
			return
		}

		var req updateConnectorRequest
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			log.Warn("update connector: invalid request", zap.String("id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		updates := connector.Update{
			Name:    req.Name,
			Enabled: req.Enabled,
			TeamID:  &teamID,
		}
		if err := h.d.ConnectorPersister().UpdateConnector(c.Request.Context(), connectorID, &updates); err != nil {
			log.Error("update connector: failed", zap.String("id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("connector updated", zap.String("id", cid))
		response.SuccessOK(c)
	}
}

// CountEnabled handles GET /:team_id/connector/count?type=xxx&status=enable|disable.
func (h *ConnectorHandler) Count() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		connectorType := c.Query("type")

		if connectorType == "" {
			response.Fail(c, errors.BadRequest("MISSING_TYPE", "query parameter 'type' is required"))
			return
		}

		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector team id: "+tid))
			return
		}

		options := connector.CountOptions{
			Type:   &connectorType,
			TeamID: &teamID,
		}
		switch c.Query("status") {
		case "enable":
			enabled := true
			options.Enable = &enabled
		case "disable":
			enabled := false
			options.Enable = &enabled
		}

		count, err := h.d.ConnectorPersister().Count(c.Request.Context(), options)
		if err != nil {
			log.Error("count connectors: failed", zap.String("type", connectorType), zap.Error(err))
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, gin.H{"type": connectorType, "count": count})
	}
}

func (h *ConnectorHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		cid := c.Param("connector_id")
		connectorID, err := uuid.Parse(cid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector id: "+cid))
			return
		}
		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument team id: "+tid))
			return
		}

		if err := h.d.ConnectorPersister().DeleteConnector(c.Request.Context(), connectorID, teamID); err != nil {
			log.Error("delete connector: failed", zap.String("team_id", tid), zap.String("connector_id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		log.Info("connector deleted", zap.String("id", cid))
		response.SuccessOK(c)
	}
}

// List handles GET /:team_id/connectors?page=1&items_per_page=20&status=enable|disable.
func (h *ConnectorHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()

		tid := c.Param("team_id")
		teamID, err := uuid.Parse(tid)
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid argument connector team id: "+tid))
			return
		}

		var req listConnectorRequest
		if err := c.ShouldBindWith(&req, binding.Query); err != nil {
			log.Warn("list connectors: invalid request", zap.Error(err))
			response.Fail(c, err)
			return
		}

		options := connector.ListOptions{
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
		if req.Type != nil {
			options.Type = req.Type
		}

		if req.Status != nil {
			switch *req.Status {
			case "enable":
				enabled := true
				options.Enable = &enabled
			case "disable":
				enabled := false
				options.Enable = &enabled
			default:
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "query parameter 'status' must be enable or disable"))
				return
			}
		}

		records, err := h.d.ConnectorPersister().ListConnectors(c.Request.Context(), &options)
		if err != nil {
			log.Error("list connectors: query failed", zap.Error(err))
			response.Fail(c, err)
			return
		}

		resp := make([]connectorResponse, 0, len(records))
		for _, r := range records {
			resp = append(resp, connectorResponse{
				ID:        r.ID,
				Type:      r.Type,
				Name:      r.Name,
				Labels:    r.Labels,
				Enabled:   r.Enabled,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
		}
		log.Debug("list connectors", zap.Int("count", len(resp)))
		response.SuccessWithData(c, map[string]interface{}{
			"count":      len(resp),
			"connectors": resp,
		})
	}
}
