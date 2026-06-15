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
	flowHandler "github.com/muixstudio/clio/internal/api/handler/flow"
	"github.com/muixstudio/clio/internal/driver/config"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/logger"
	"github.com/muixstudio/clio/internal/session"
)

type dependencies interface {
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

	session.SessionPersisterProvider

	logger.Logger
	config.Provider
}

func Register(router *gin.RouterGroup, deps dependencies) {
	v1 := router.Group("/v1")
	handler := initHandler(deps)

	// login flow
	{
		v1.POST("/login", handler.Login())
		v1.GET("/login/flows", handler.GetLoginFlow())
	}

	// logout flow
	{
		v1.POST("/logout", handler.Logout())
		v1.GET("/logout/flows", handler.GetLogoutFlow())
	}

	// recovery flow
	{
		v1.POST("/recovery", handler.Recovery())
		v1.GET("/recovery/flows", handler.GetRecoveryFlow())
	}

	// registration flow
	{
		v1.POST("/registration", handler.Registration())
		v1.GET("/registration/flows", handler.GetRegistrationFlow())
	}

	// verification flow
	{
		v1.POST("/verification", handler.Verification())
		v1.GET("/verification/flows", handler.GetVerificationFlow())
	}

}

func initHandler(deps dependencies) *flowHandler.Handler {
	return flowHandler.NewHandler(deps)
}
