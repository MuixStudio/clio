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
	teamHandler "github.com/muixstudio/clio/internal/api/handler/team"
	"github.com/muixstudio/clio/internal/domain/alert/repository"
	connectorDomain "github.com/muixstudio/clio/internal/domain/connector"
	connectorRepo "github.com/muixstudio/clio/internal/domain/connector/repository"
	teamRepo "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/logger"
)

type dependencies interface {
	connectorRepo.ConnectorPersisterProvider
	repository.AlertPersisterProvider
	connectorDomain.ConnectorProviderProvider
	teamRepo.TeamPersisterProvider
	teamRepo.TeamMemberPersisterProvider

	logger.Logger

	config.Provider
}

func Register(router *gin.RouterGroup, deps dependencies) {
	v1 := router.Group("/v1")
	handler := initHandler(deps)
	{
		v1.GET("/teams", handler.List())
		v1.POST("/team", handler.Create())
		v1.GET("/team/:team_id", handler.Get())
		v1.PATCH("/team/:team_id", handler.Update())
		v1.DELETE("/team/:team_id", handler.Delete())

		v1.GET("/team/:team_id/members", handler.ListMembers())
		v1.POST("/team/:team_id/member", handler.AddMember())
		v1.DELETE("/team/:team_id/member/:member_id", handler.RemoveMember())

		v1.GET("/user/:user_id/teams", handler.GetUserTeams()) // Deprecation Warning
	}
}

func initHandler(deps dependencies) *teamHandler.TeamHandler {
	return teamHandler.NewTeamHandler(deps)
}
