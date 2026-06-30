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

package webhook

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"go.uber.org/zap"
)

// Receive handles POST /webhook/:type/:id.
func (h *WebhookHandler) Receive() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.d.LoggerProvider()
		cid := c.Param("connector_id")
		connectorID, err := uuid.Parse(cid)
		if err != nil {
			response.Fail(c, errors.NotFound("RECORD_NOT_FOUND", "connector not found: "+cid))
			return
		}
		routeType := c.Param("type")

		connector, err := h.d.ConnectorPersister().GetConnector(c.Request.Context(), connectorID)
		if err != nil {
			log.Warn("webhook receive: connector not found", zap.String("connector_id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		connectorProvider, ok := h.d.ConnectorProviders()[routeType]
		if !ok {
			log.Warn("webhook receive: no provider for type", zap.String("type", routeType), zap.String("connector_id", cid))
			response.Fail(
				c,
				errors.BadRequest(
					"NO_PROVIDER",
					fmt.Sprintf("no provider for type %s", routeType),
				),
			)
			return
		}

		if err = h.authentication(c, connector.Token); err != nil {
			log.Warn("webhook receive: authentication failed", zap.String("type", routeType), zap.String("id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}
		bt, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		alerts, err := connectorProvider.Normalize(bt)
		if err != nil {
			log.Error("webhook receive: parse alerts failed", zap.String("type", routeType), zap.String("connector_id", cid), zap.Error(err))
			response.Fail(c, err)
			return
		}

		for _, a := range alerts {
			a.TeamID = connector.TeamID
			a.ConnectorID = connector.ID
		}

		err = h.d.AlertPersister().SaveAlerts(c.Request.Context(), alerts)
		if err != nil {
			log.Error("webhook receive: push alerts failed", zap.String("type", routeType), zap.Error(err))
			response.Fail(c, err)
			return
		}

		// Publish only after persistence succeeds. The database remains the source
		// of truth; dispatch can ack on receipt and rebuild/replay later if needed.

		err = h.d.AlertPersister().Save(c.Request.Context(), alerts[0])
		if err != nil {
			log.Error("webhook receive: publish alerts failed", zap.String("type", routeType), zap.Error(err))
			response.Fail(c, errors.InternalServerError("INTERNAL_SERVER_ERROR", "failed to publish alerts").WithCause(err))
			return
		}

		response.SuccessNoContent(c)
	}
}
func (h *WebhookHandler) authentication(c *gin.Context, token string) error {
	if token == "" {
		return errors.Unauthorized("UNAUTHORIZED", "unauthorized")
	}
	if c.Query("token") != token {
		return errors.Unauthorized("UNAUTHORIZED", "unauthorized")
	}
	return nil
}
