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
	connectorHannder "github.com/muixstudio/clio/internal/api/handler/connector"
	"github.com/muixstudio/clio/internal/connector"
	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/logger"
)

type dependencies interface {
	connector.ConnectorPersisterProvider
	alert.AlertPersisterProvider
	connector.ConnectorProviderProvider

	logger.Logger

	config.Provider
}

func Register(router *gin.RouterGroup, deps dependencies) {
	v1 := router.Group("/v1")
	handler := initHandler(deps)

	{
		v1.GET("/team/:team_id/connectors", handler.List())
		v1.POST("/team/:team_id/connector", handler.Create())
		v1.GET("/team/:team_id/connector/count", handler.Count())
		v1.PATCH("/team/:team_id/connector/:connector_id", handler.Update())
		v1.DELETE("/team/:team_id/connector/:connector_id", handler.Delete())
	}
}

func initHandler(deps dependencies) *connectorHannder.ConnectorHandler {
	return connectorHannder.NewConnectorHandler(deps)
}
