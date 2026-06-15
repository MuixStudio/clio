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
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/response"
)

type verificationRequest struct {
	Method string          `json:"method"`
	FlowID string          `json:"flow_id,omitempty"`
	Data   json.RawMessage `json:"data"`
}

type verificationResponse struct {
	FlowID string     `json:"flow_id"`
	State  flow.State `json:"state"`
	Email  string     `json:"email,omitempty"`
}

func (h *Handler) Verification() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req verificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errors.BadRequest("INVALID_BODY", err.Error()))
			return
		}

		var f *verificationFlow.Flow
		if req.FlowID != "" {
			// Step 2: resume an existing flow.
			id, err := uuid.Parse(req.FlowID)
			if err != nil {
				response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow_id"))
				return
			}
			existing, err := h.d.VerificationFlowPersister().GetVerificationFlow(ctx, id)
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
			f = verificationFlow.New(h.d.Config().SelfService.Flows.Verification.Lifespan)
			if err := h.d.VerificationHookExecutor().PreHook(ctx, f); err != nil {
				response.Fail(c, errors.BadRequest("PRE_HOOK_FAILED", err.Error()))
				return
			}
			if err := h.d.VerificationFlowPersister().CreateVerificationFlow(ctx, f); err != nil {
				response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
				return
			}
		}

		// Strategy dispatch — identical pattern to login/registration.
		var strategyErr error
		for _, s := range h.d.VerificationStrategies() {
			if string(s.ID()) != req.Method {
				continue
			}
			strategyErr = s.Verify(c, f, req.Data)

			if errors.Is(strategyErr, flow.ErrCompletedByStrategy) {
				// Strategy wrote the response (Step 1: sent code + flow_id).
				_ = h.d.VerificationFlowPersister().UpdateVerificationFlow(ctx, f)
				return
			}
			if errors.Is(strategyErr, flow.ErrStrategyNotResponsible) {
				continue
			}
			if strategyErr != nil {
				f.SetState(flow.StateFailed)
				_ = h.d.VerificationFlowPersister().UpdateVerificationFlow(ctx, f)
				response.Fail(c, errors.BadRequest("STRATEGY_ERROR", strategyErr.Error()))
				return
			}
			break
		}

		if strategyErr != nil && !errors.Is(strategyErr, flow.ErrCompletedByStrategy) {
			response.Fail(c, errors.BadRequest("METHOD_NOT_FOUND", "no strategy matched method: "+req.Method))
			return
		}

		// Step 2 succeeded: run post-hooks and return completion response.
		if err := h.d.VerificationHookExecutor().PostHook(ctx, f); err != nil {
			response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
			return
		}

		f.SetState(verificationFlow.StatePassedChallenge)
		_ = h.d.VerificationFlowPersister().UpdateVerificationFlow(ctx, f)

		response.SuccessWithData(c, verificationResponse{
			FlowID: f.GetID().String(),
			State:  f.GetState(),
			Email:  f.Email,
		})
	}
}

func (h *Handler) GetVerificationFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid flow id"))
			return
		}

		f, err := h.d.VerificationFlowPersister().GetVerificationFlow(c.Request.Context(), id)
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
