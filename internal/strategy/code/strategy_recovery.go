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

package code

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/muixstudio/clio/internal/flow"
	recflow "github.com/muixstudio/clio/internal/flow/recovery"
	"github.com/muixstudio/clio/internal/identity"
)

// Compile-time assertion: Strategy implements recflow.Strategy.
var _ recflow.Strategy = (*Strategy)(nil)

// recoveryInitBody is the Step-1 payload: client provides the account email.
type recoveryInitBody struct {
	Email string `json:"email"`
}

// recoveryVerifyBody is the Step-2 payload: client submits the received code.
type recoveryVerifyBody struct {
	Code string `json:"code"`
}

// Recover implements recflow.Strategy — two-step account recovery via email code.
//
// Step 1 — body contains {"email":"u@example.com"}:
//   - Looks up the identity that owns this email.
//   - If found, generates and sends a 6-digit code via Courier.
//   - If NOT found, silently succeeds (anti-enumeration: attacker can't tell
//     whether the email is registered by observing the response).
//   - Writes {"flow_id":"...","state":"sent_email"} and returns ErrCompletedByStrategy.
//
// Step 2 — body contains {"code":"123456"}, flow_id in the outer envelope:
//   - Looks up the code by flow_id + code value.
//   - Validates it is not expired or used.
//   - Returns the identity so the handler can issue a recovery session.
//
// The recovery session allows the user to change their password without knowing
// the old one — it's the "prize" of proving email ownership.
func (s *Strategy) Recover(c *gin.Context, f *recflow.Flow, body json.RawMessage) (*identity.Identity, error) {
	ctx := c.Request.Context()

	// Step 2: code present → complete recovery.
	var vb recoveryVerifyBody
	if err := json.Unmarshal(body, &vb); err == nil && vb.Code != "" {
		return s.completeRecovery(ctx, f, vb.Code)
	}

	// Step 1: email present → initiate recovery.
	var ib recoveryInitBody
	if err := json.Unmarshal(body, &ib); err != nil || ib.Email == "" {
		return nil, flow.ErrStrategyNotResponsible
	}

	return nil, s.initiateRecovery(c, f, ib.Email)
}

// initiateRecovery handles Step 1: look up the identity, send the code.
func (s *Strategy) initiateRecovery(c *gin.Context, f *recflow.Flow, email string) error {
	ctx := c.Request.Context()

	pool, ok := s.d.(VerificationIdentityPool)
	if !ok {
		return errors.New("code: driver does not implement VerificationIdentityPool")
	}

	// Anti-enumeration: always return the same "sent_email" response regardless
	// of whether the email is registered. Only send the code when it is.
	i, err := pool.FindIdentityByEmail(ctx, email)
	if err == nil {
		// Identity found — send the code.
		if sendErr := s.createAndSendCode(ctx, f.GetID(), "recovery", email); sendErr != nil {
			return sendErr
		}
	}
	// Silently skip sending if identity not found.

	// Store the email on the flow so Step 2 can reference it without the client
	// re-submitting it — same pattern as the verification flow.
	f.Email = email
	f.SetState(recflow.StateSentEmail)

	response.SuccessWithData(c, map[string]any{
		"flow_id": f.GetID().String(),
		"state":   recflow.StateSentEmail,
		"email":   email,
	})

	_ = i // used only for the send decision above
	return flow.ErrCompletedByStrategy
}

// completeRecovery handles Step 2: verify the code and return the identity.
func (s *Strategy) completeRecovery(ctx context.Context, f *recflow.Flow, code string) (*identity.Identity, error) {
	vc, err := s.d.VerificationCodePersister().FindVerificationCode(ctx, f.GetID(), code)
	if err != nil {
		return nil, errors.New("invalid or expired code")
	}
	if vc.IsExpired() {
		return nil, errors.New("code has expired, please request a new one")
	}
	if vc.IsUsed() {
		return nil, errors.New("code has already been used")
	}

	// Consume the code immediately — prevents replay if the handler fails later.
	_ = s.d.VerificationCodePersister().UseVerificationCode(ctx, vc.ID)

	// Look up the identity so the handler can issue a session for it.
	pool, ok := s.d.(VerificationIdentityPool)
	if !ok {
		return nil, errors.New("code: driver does not implement VerificationIdentityPool")
	}
	i, err := pool.FindIdentityByEmail(ctx, vc.Address)
	if err != nil {
		// Identity deleted between Step 1 and Step 2 — extremely rare.
		return nil, errors.New("account not found")
	}

	return i, nil
}
