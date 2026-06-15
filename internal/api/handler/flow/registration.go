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
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/flow"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"github.com/muixstudio/clio/internal/session"
)

type registrationRequest struct {
	// Method selects the strategy (e.g. "password", "oidc", "code").
	Method string `json:"method"`
	// FlowID, when set, resumes an existing flow instead of creating a new one.
	// Required for two-step strategies (e.g. code): the client receives a
	// flow_id in Step 1 and must echo it back in Step 2.
	FlowID string `json:"flow_id,omitempty"`
	// Data is forwarded verbatim to the chosen strategy.
	Data json.RawMessage `json:"data"`
}

type registrationResponse struct {
	FlowID  string           `json:"flow_id"`
	State   flow.State       `json:"state"`
	Session *session.Session `json:"session,omitempty"`
}

func (h *Handler) Registration() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req registrationRequest
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
			response.Fail(c, err)
			return
		}

		var f *registrationFlow.Flow
		if req.FlowID != "" {
			// Resume an existing flow — needed by two-step strategies (e.g. code).
			id, err := uuid.Parse(req.FlowID)
			if err != nil {
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow_id").WithCause(err))
				return
			}
			existing, err := h.d.RegistrationFlowPersister().GetRegistrationFlow(ctx, id)
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
			f = registrationFlow.New(h.d.Config().SelfService.Flows.Registration.Lifespan)
			if err := h.d.RegistrationHookExecutor().PreHook(ctx, f); err != nil {
				response.Fail(c, errors.BadRequest("PRE_HOOK_FAILED", err.Error()))
				return
			}
			if err := h.d.RegistrationFlowPersister().CreateRegistrationFlow(ctx, f); err != nil {
				response.Fail(c, err)
				return
			}
		}

		// Strategy dispatch: iterate all registered strategies.
		var ident *identity.Identity
		for _, s := range h.d.RegistrationStrategies() {
			if string(s.ID()) != req.Method {
				continue
			}
			var err error
			ident, err = s.Register(c, f, req.Data)

			// Strategy wrote the response itself (e.g. OIDC returning an auth URL).
			//if errors.Is(err, flow.ErrCompletedByStrategy) {
			//	return
			//}
			//if errors.Is(err, flow.ErrStrategyNotResponsible) {
			//	continue
			//}
			if err != nil {
				f.SetState(flow.StateFailed)
				_ = h.d.RegistrationFlowPersister().UpdateRegistrationFlow(ctx, f)
				response.Fail(c, err)
				return
			}
			f.Active = s.ID()
			break
		}

		if ident == nil {
			response.Fail(c, errors.BadRequest("METHOD_NOT_FOUND", "no strategy matched method: "+req.Method))
			return
		}

		if err := h.d.RegistrationHookExecutor().PostHook(ctx, f, ident); err != nil {
			response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
			return
		}

		sess, err := h.issueSession(ctx, ident, f.Active)
		if err != nil {
			response.Fail(c, errors.InternalServerError("SESSION_ISSUE_FAILED", err.Error()))
			return
		}

		f.SetState(flow.StateCompleted)
		_ = h.d.RegistrationFlowPersister().UpdateRegistrationFlow(ctx, f)

		response.SuccessWithData(c, registrationResponse{
			FlowID:  f.GetID().String(),
			State:   f.GetState(),
			Session: sess,
		})
	}
}

func (h *Handler) GetRegistrationFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow id"))
			return
		}

		f, err := h.d.RegistrationFlowPersister().GetRegistrationFlow(c.Request.Context(), id)
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
