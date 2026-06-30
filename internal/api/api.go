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

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/api/router"
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
	cbinding "github.com/muixstudio/clio/internal/infra/binding"
	"github.com/muixstudio/clio/internal/infra/orgctx"
	alertRouter "github.com/muixstudio/clio/internal/infra/zus/dispatch"
	"github.com/muixstudio/clio/internal/logger"
	"github.com/muixstudio/clio/internal/session"
)

type dependencies interface {
	connectorRepo.ConnectorPersisterProvider
	alertRepo.AlertPersisterProvider
	alertDomain.AlertPublisherProvider
	alertRouter.RouteWriterProvider
	alertRouter.RouteReloaderProvider
	alertRouter.RouteTreeLoaderProvider
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

func NewApiEngine(deps dependencies) (*gin.Engine, error) {
	v, err := cbinding.NewValidator()
	if err != nil {
		return nil, err
	}
	binding.Validator = v
	binding.JSON = cbinding.JSON

	if deps.Config().Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), corsMiddleware(), orgMiddleware())

	router.Register(engine, deps)
	return engine, nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Session-Token, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func orgMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := orgctx.WithOrgID(c.Request.Context(), uuid.NullUUID{})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
