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
	channelHandler "github.com/muixstudio/clio/internal/api/handler/channel"
	channelDomain "github.com/muixstudio/clio/internal/channel"
	"github.com/muixstudio/clio/internal/logger"
)

type dependencies interface {
	channelDomain.ChannelPersisterProvider
	channelDomain.ChannelMemberPersisterProvider

	logger.Logger
}

func Register(router *gin.RouterGroup, deps dependencies) {
	v1 := router.Group("/v1")
	handler := initHandler(deps)

	{
		v1.GET("/team/:team_id/channels", handler.List())
		v1.POST("/team/:team_id/channel", handler.Create())
		v1.GET("/team/:team_id/channel/:channel_id", handler.Get())
		v1.PATCH("/team/:team_id/channel/:channel_id", handler.Update())
		v1.DELETE("/team/:team_id/channel/:channel_id", handler.Delete())

		v1.GET("/team/:team_id/channel/:channel_id/members", handler.ListMembers())
		v1.POST("/team/:team_id/channel/:channel_id/member", handler.AddMember())
		v1.GET("/team/:team_id/channel/:channel_id/member/:user_id", handler.GetMember())
		v1.PATCH("/team/:team_id/channel/:channel_id/member/:user_id", handler.UpdateMemberRole())
		v1.DELETE("/team/:team_id/channel/:channel_id/member/:user_id", handler.RemoveMember())
	}
}

func initHandler(deps dependencies) *channelHandler.ChannelHandler {
	return channelHandler.NewChannelHandler(deps)
}
