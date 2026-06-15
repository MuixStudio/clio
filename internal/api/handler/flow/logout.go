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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/flow"
	logoutFlow "github.com/muixstudio/clio/internal/flow/logout"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/jwt"
	"github.com/muixstudio/clio/internal/infra/response"
)

func (h *Handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token := c.GetHeader("X-Session-Token")
		if token == "" {
			response.Fail(c, errors.Unauthorized("MISSING_TOKEN", "X-Session-Token header is required"))
			return
		}

		claims, err := jwt.ParseAccessToken(token, h.d.Config().Session.JWTSecret)
		if err != nil {
			response.Fail(c, errors.Unauthorized("INVALID_TOKEN", "invalid or expired access token"))
			return
		}
		sess, err := h.d.SessionPersister().GetSessionByID(ctx, claims.SessionID)
		if err != nil || !sess.IsActive() {
			response.Fail(c, errors.Unauthorized("SESSION_NOT_FOUND", "session not found or already expired"))
			return
		}

		f := logoutFlow.New(h.d.Config().SelfService.Flows.Logout.Lifespan, sess.IdentityID)
		if err := h.d.LogoutFlowPersister().CreateLogoutFlow(ctx, f); err != nil {
			response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
			return
		}

		if err := h.d.LogoutHookExecutor().PreHook(ctx, f); err != nil {
			f.SetState(flow.StateFailed)
			_ = h.d.LogoutFlowPersister().UpdateLogoutFlow(ctx, f)
			response.Fail(c, errors.BadRequest("PRE_HOOK_FAILED", err.Error()))
			return
		}

		if err := h.d.SessionPersister().RevokeSession(ctx, sess.ID); err != nil {
			f.SetState(flow.StateFailed)
			_ = h.d.LogoutFlowPersister().UpdateLogoutFlow(ctx, f)
			response.Fail(c, errors.InternalServerError("REVOKE_FAILED", err.Error()))
			return
		}

		if err := h.d.LogoutHookExecutor().PostHook(ctx, f); err != nil {
			f.SetState(flow.StateFailed)
			_ = h.d.LogoutFlowPersister().UpdateLogoutFlow(ctx, f)
			response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
			return
		}

		f.SetState(flow.StateCompleted)
		_ = h.d.LogoutFlowPersister().UpdateLogoutFlow(ctx, f)

		c.Status(http.StatusNoContent)
		return
	}
}

func (h *Handler) GetLogoutFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow id"))
			return
		}

		f, err := h.d.LogoutFlowPersister().GetLogoutFlow(c.Request.Context(), id)
		if err != nil {
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, map[string]any{
			"flow_id":     f.GetID().String(),
			"state":       f.GetState(),
			"identity_id": f.IdentityID.String(),
		})
	}
}
