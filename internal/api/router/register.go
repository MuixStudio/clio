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

package router

import (
	"github.com/gin-gonic/gin"
	v1Alert "github.com/muixstudio/clio/internal/api/router/alert/v1"
	v1Channel "github.com/muixstudio/clio/internal/api/router/channel/v1"
	v1Connector "github.com/muixstudio/clio/internal/api/router/connector/v1"
	v1Flow "github.com/muixstudio/clio/internal/api/router/flow/v1"
	v1Session "github.com/muixstudio/clio/internal/api/router/session/v1"
	v1Team "github.com/muixstudio/clio/internal/api/router/team/v1"
	v1Webhook "github.com/muixstudio/clio/internal/api/router/webhook/v1"
	alertDomain "github.com/muixstudio/clio/internal/domain/alert"
	alertRepo "github.com/muixstudio/clio/internal/domain/alert/repository"
	"github.com/muixstudio/clio/internal/domain/channel/repository"
	connectorDomain "github.com/muixstudio/clio/internal/domain/connector"
	connectorRepo "github.com/muixstudio/clio/internal/domain/connector/repository"
	teamRepo "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/driver/config"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/logger"
	"github.com/muixstudio/clio/internal/session"
)

type connectorDependencies interface {
	connectorRepo.ConnectorPersisterProvider
	alertRepo.AlertPersisterProvider
	alertDomain.AlertPublisherProvider
	connectorDomain.ConnectorProviderProvider
	repository.ChannelPersisterProvider
	repository.ChannelMemberPersisterProvider
	teamRepo.TeamPersisterProvider
	teamRepo.TeamMemberPersisterProvider

	session.SessionPersisterProvider
	session.SessionTokenExchangeCodePersisterProvider

	loginFlow.FlowPersistenceProvider
	loginFlow.StrategyProvider
	loginFlow.HookExecutorProvider

	logoutLfow.FlowPersistenceProvider
	logoutLfow.HookExecutorProvider

	recoveryFlow.FlowPersistenceProvider
	recoveryFlow.StrategyProvider
	recoveryFlow.HookExecutorProvider

	registrationFlow.FlowPersistenceProvider
	registrationFlow.StrategyProvider
	registrationFlow.HookExecutorProvider

	verificationFlow.FlowPersistenceProvider
	verificationFlow.StrategyProvider
	verificationFlow.HookExecutorProvider

	logger.Logger

	config.Provider
}

func Register(eg *gin.Engine, deps connectorDependencies) {
	apiGroup := eg.Group("/api")

	v1Flow.Register(apiGroup, deps)
	v1Session.Register(apiGroup, deps)

	v1Connector.Register(apiGroup, deps)
	v1Webhook.Register(apiGroup, deps)
	v1Alert.Register(apiGroup, deps)
	v1Channel.Register(apiGroup, deps)

	v1Team.Register(apiGroup, deps)

}
