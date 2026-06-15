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

package flow

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/jwt"

	"github.com/muixstudio/clio/internal/driver/config"
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

// Handler handles connector management endpoints.
type Handler struct {
	d dependencies
}

func (h *Handler) issueSession(ctx context.Context, i *identity.Identity, method identity.CredentialsType) (*session.Session, error) {
	sess := session.New(i, h.d.Config().Session.Lifespan)
	sess.CompletedLoginFor(method, identity.AuthenticatorAssuranceLevel1)
	sess.SetAuthenticatorAssuranceLevel()
	if err := h.d.SessionPersister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	at, err := jwt.GenerateAccessToken(sess, h.d.Config().Session.JWTSecret, h.d.Config().Session.AccessTokenLifespan)
	if err != nil {
		return nil, err
	}
	sess.AccessToken = at
	return sess, nil
}

func setSessionCookies(c *gin.Context, sess *session.Session) {
	maxAge := int(time.Until(sess.ExpiresAt).Seconds())
	for _, cookie := range []http.Cookie{
		{Name: "sid", Value: sess.ID.String()},
		{Name: "access_token", Value: sess.AccessToken},
		{Name: "refresh_token", Value: sess.RefreshToken},
	} {
		cookie.MaxAge = maxAge
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteLaxMode
		http.SetCookie(c.Writer, &cookie)
	}
}

// NewHandler returns a FlowHandler.
func NewHandler(d dependencies) *Handler {
	return &Handler{d: d}
}
