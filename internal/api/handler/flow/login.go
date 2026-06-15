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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	cbinding "github.com/muixstudio/clio/internal/infra/binding"
	"github.com/muixstudio/clio/internal/infra/jwt"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/muixstudio/clio/internal/flow"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/session"
)

type loginRequest struct {
	Method identity.CredentialsType `json:"method" validate:"required"`
	// FlowID, when set, resumes an existing flow instead of creating a new one.
	// Required for two-step strategies (e.g. code): the client receives a
	// flow_id in Step 1 and must echo it back in Step 2.
	FlowID string `json:"flow_id,omitempty"`
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		refresh := c.Query("refresh") == "true"

		// Fetch the caller's existing session (if any) and pass it to strategies.
		// sess is forwarded so AAL2 strategies know whose second factor to verify.
		var existingSess *session.Session
		if token := c.GetHeader("X-Session-Token"); token != "" {
			claims, err := jwt.ParseAccessToken(token, h.d.Config().Session.JWTSecret)
			if err == nil {
				s, err := h.d.SessionPersister().GetSessionByID(ctx, claims.SessionID)
				if err == nil && s.IsActive() {
					if !refresh {
						response.Fail(c, errors.BadRequest("ALREADY_LOGGED_IN", "already logged in; pass ?refresh=true to re-authenticate"))
						return
					}
					existingSess = s
				}
			}
		}

		// Decode the envelope first so we can inspect flow_id.
		var req loginRequest
		if err := c.ShouldBindBodyWith(&req, cbinding.JSON); err != nil {
			response.Fail(c, err)
			return
		}

		var f *loginFlow.Flow
		if req.FlowID != "" {
			// Resume an existing flow — needed by two-step strategies (e.g. code).
			id, err := uuid.Parse(req.FlowID)
			if err != nil {
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow_id"))
				return
			}
			existing, err := h.d.LoginFlowPersister().GetLoginFlow(ctx, id)
			if err != nil {
				response.Fail(c, err)
				return
			}
			if existing.IsExpired() {
				response.Fail(c, errors.Gone("FLOW_EXPIRED", "flow has expired, please start over"))
				return
			}
			f = existing
		} else {
			// Create a fresh flow (single-step strategies: password, oidc).
			f = loginFlow.New(h.d.Config().SelfService.Flows.Login.Lifespan, refresh)
			if err := h.d.LoginHookExecutor().PreHook(ctx, f); err != nil {
				response.Fail(c, errors.BadRequest("PRE_HOOK_FAILED", err.Error()))
				return
			}
			if err := h.d.LoginFlowPersister().CreateLoginFlow(ctx, f); err != nil {
				response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
				return
			}
		}

		loginStrategies := h.d.LoginStrategies()
		s, ok := loginStrategies[req.Method]
		if !ok {
			response.Fail(c, errors.BadRequest("METHOD_NOT_FOUND", "no strategy matched method: "+req.Method.String()))
			return
		}
		ident, err := s.Login(c, f, existingSess)
		if ident == nil && err == nil {
			return
		}
		if err != nil {
			f.SetState(flow.StateFailed)
			_ = h.d.LoginFlowPersister().UpdateLoginFlow(ctx, f)
			response.Fail(c, err)
			return
		}
		f.Active = s.ID()

		sess, err := h.issueSession(ctx, ident, f.Active)
		if err != nil {
			response.Fail(c, errors.InternalServerError("SESSION_ISSUE_FAILED", err.Error()))
			return
		}

		if err := h.d.LoginHookExecutor().PostHook(ctx, f, ident, sess); err != nil {
			response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
			return
		}

		f.SetState(flow.StateCompleted)
		_ = h.d.LoginFlowPersister().UpdateLoginFlow(ctx, f)

		setSessionCookies(c, sess)
		response.SuccessWithData(c, struct {
			FlowID string     `json:"flow_id"`
			State  flow.State `json:"state"`
		}{
			FlowID: f.GetID().String(),
			State:  f.GetState(),
		})
	}
}

func (h *Handler) GetLoginFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow id"))
			return
		}

		f, err := h.d.LoginFlowPersister().GetLoginFlow(c.Request.Context(), id)
		if err != nil {
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, map[string]any{
			"flow_id": f.GetID().String(),
			"state":   f.GetState(),
		})
	}
}
