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

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/alert"
	"github.com/muixstudio/clio/internal/api/handler/webhook"
	"github.com/muixstudio/clio/internal/connector"
	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/logger"
)

type connectorDependencies interface {
	connector.ConnectorPersisterProvider
	alert.AlertPersisterProvider
	connector.ConnectorProviderProvider

	logger.Logger

	config.Provider
}

func Register(router *gin.RouterGroup, deps connectorDependencies) {
	v1 := router.Group("/v1")
	handler := initHandler(deps)

	{
		v1.POST("/webhook/:type/:connector_id", handler.Receive())
	}
}

func initHandler(deps connectorDependencies) *webhook.WebhookHandler {
	return webhook.NewWebhookHandler(deps)
}
