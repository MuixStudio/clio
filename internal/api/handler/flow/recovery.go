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
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/flow"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
	"github.com/muixstudio/clio/internal/session"
)

type recoveryRequest struct {
	Method string          `json:"method"`
	FlowID string          `json:"flow_id,omitempty"`
	Data   json.RawMessage `json:"data"`
}

type recoveryResponse struct {
	FlowID string     `json:"flow_id"`
	State  flow.State `json:"state"`
	Email  string     `json:"email,omitempty"`
	// Session is populated only after Step 2 (code verified, recovery complete).
	// The client uses this session token to change their password via settings.
	Session *session.Session `json:"session,omitempty"`
}

func (h *Handler) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req recoveryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errors.BadRequest("INVALID_BODY", err.Error()))
			return
		}

		var f *recoveryFlow.Flow
		if req.FlowID != "" {
			// Step 2: resume an existing flow.
			id, err := uuid.Parse(req.FlowID)
			if err != nil {
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow_id"))
				return
			}
			existing, err := h.d.RecoveryFlowPersister().GetRecoveryFlow(ctx, id)
			if err != nil {
				response.Fail(c, errors.NotFound("FLOW_NOT_FOUND", err.Error()))
				return
			}
			if existing.IsExpired() {
				response.Fail(c, errors.Gone("FLOW_EXPIRED", "flow has expired, please start over"))
				return
			}
			f = existing
		} else {
			// Step 1: create a fresh flow.
			f = recoveryFlow.New(h.d.Config().SelfService.Flows.Recovery.Lifespan)
			if err := h.d.RecoveryHookExecutor().PreHook(ctx, f); err != nil {
				response.Fail(c, errors.BadRequest("PRE_HOOK_FAILED", err.Error()))
				return
			}
			if err := h.d.RecoveryFlowPersister().CreateRecoveryFlow(ctx, f); err != nil {
				response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
				return
			}
		}

		// Strategy dispatch.
		var ident *identity.Identity
		for _, s := range h.d.RecoveryStrategies() {
			if string(s.ID()) != req.Method {
				continue
			}
			var err error
			ident, err = s.Recover(c, f, req.Data)

			if errors.Is(err, flow.ErrCompletedByStrategy) {
				// Step 1: strategy sent the code and wrote the response.
				_ = h.d.RecoveryFlowPersister().UpdateRecoveryFlow(ctx, f)
				return
			}
			if errors.Is(err, flow.ErrStrategyNotResponsible) {
				continue
			}
			if err != nil {
				f.SetState(flow.StateFailed)
				_ = h.d.RecoveryFlowPersister().UpdateRecoveryFlow(ctx, f)
				response.Fail(c, errors.BadRequest("STRATEGY_ERROR", err.Error()))
				return
			}
			break
		}

		if ident == nil {
			response.Fail(c, errors.BadRequest("METHOD_NOT_FOUND", "no strategy matched method: "+req.Method))
			return
		}

		// Step 2 succeeded: issue a session so the user can change their password.
		sess, err := h.issueSession(ctx, ident, identity.CredentialsTypeCode)
		if err != nil {
			response.Fail(c, errors.InternalServerError("SESSION_ISSUE_FAILED", err.Error()))
			return
		}

		if err := h.d.RecoveryHookExecutor().PostHook(ctx, f, ident, sess); err != nil {
			response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
			return
		}

		f.SetState(recoveryFlow.StatePassedChallenge)
		_ = h.d.RecoveryFlowPersister().UpdateRecoveryFlow(ctx, f)

		response.SuccessWithData(c, recoveryResponse{
			FlowID:  f.GetID().String(),
			State:   f.GetState(),
			Email:   f.Email,
			Session: sess,
		})
	}
}

func (h *Handler) GetRecoveryFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow id"))
			return
		}

		f, err := h.d.RecoveryFlowPersister().GetRecoveryFlow(c.Request.Context(), id)
		if err != nil {
			response.Fail(c, err)
			return
		}

		response.SuccessWithData(c, map[string]any{
			"flow_id": f.GetID().String(),
			"state":   f.GetState(),
			"email":   f.Email,
		})
	}
}
