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
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/muixstudio/clio/internal/flow"
	verflow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
)

// VerifiableAddressPool is the narrow persistence interface the code strategy
// needs in order to find and mark verifiable addresses.
type VerifiableAddressPool interface {
	// FindOrCreateVerifiableAddress returns the verifiable address record for
	// the given identity + via + value, creating it (as "pending") if missing.
	FindOrCreateVerifiableAddress(ctx context.Context, identityID uuid.UUID, via, value string) (*identity.VerifiableAddress, error)
	// MarkVerifiableAddressVerified stamps verified=true and verified_at=now.
	MarkVerifiableAddressVerified(ctx context.Context, id uuid.UUID) error
}

// VerificationIdentityPool extends IdentityPool with an email-based lookup
// that searches across all credential types.
// Used during verification initiation to find whose address to verify.
type VerificationIdentityPool interface {
	FindIdentityByEmail(ctx context.Context, email string) (*identity.Identity, error)
}

// Compile-time assertion: Strategy implements verflow.Strategy.
var _ verflow.Strategy = (*Strategy)(nil)

// verifyInitBody is the Step-1 payload: the client provides the email to verify.
type verifyInitBody struct {
	Email string `json:"email"`
}

// verifyCompleteBody is the Step-2 payload: the client submits the received code.
type verifyCompleteBody struct {
	Code string `json:"code"`
}

// Verify implements verflow.Strategy — two-step email address verification.
//
// Step 1 — body contains {"email":"u@example.com"}:
//   - Looks up the identity that owns this email (silently succeeds if not found
//     to prevent email enumeration).
//   - Creates/updates the VerifiableAddress record with status "sent".
//   - Generates and sends a 6-digit code via Courier.
//   - Writes {"flow_id":"...","state":"sent_email"} and returns ErrCompletedByStrategy.
//
// Step 2 — body contains {"code":"123456"}, flow_id in the outer envelope:
//   - Looks up the code by flow_id + code value.
//   - Marks the VerifiableAddress as verified.
//   - Returns nil so the handler can write the completion response.
func (s *Strategy) Verify(c *gin.Context, f *verflow.Flow, body json.RawMessage) error {
	ctx := c.Request.Context()

	// Step 2: code present → complete verification.
	var cb verifyCompleteBody
	if err := json.Unmarshal(body, &cb); err == nil && cb.Code != "" {
		return s.completeVerification(ctx, f, cb.Code)
	}

	// Step 1: email present → initiate verification.
	var ib verifyInitBody
	if err := json.Unmarshal(body, &ib); err != nil || ib.Email == "" {
		return flow.ErrStrategyNotResponsible
	}

	return s.initiateVerification(c, f, ib.Email)
}

// initiateVerification handles Step 1: find-or-create the verifiable address,
// send a code, and return the flow_id to the client.
func (s *Strategy) initiateVerification(c *gin.Context, f *verflow.Flow, email string) error {
	ctx := c.Request.Context()

	pool, ok := s.d.(VerificationIdentityPool)
	if !ok {
		return errors.New("code: driver does not implement VerificationIdentityPool")
	}
	addrPool, ok := s.d.(VerifiableAddressPool)
	if !ok {
		return errors.New("code: driver does not implement VerifiableAddressPool")
	}

	// Look up which identity owns this email.
	// If no identity is found we still go through the motions (no error returned
	// to the caller) so an attacker cannot enumerate registered emails via
	// response differences.
	i, err := pool.FindIdentityByEmail(ctx, email)

	if err == nil {
		// Identity found — create/update the VerifiableAddress and send the code.
		if _, addrErr := addrPool.FindOrCreateVerifiableAddress(ctx, i.ID, "email", email); addrErr != nil {
			return addrErr
		}

		if sendErr := s.createAndSendCode(ctx, f.GetID(), "verification", email); sendErr != nil {
			return sendErr
		}
	}
	// If identity not found, we silently skip sending (anti-enumeration).
	// The client sees the same "sent_email" response either way.

	// Store the email on the flow so Step 2 can reference it.
	f.Email = email
	f.SetState(verflow.StateSentEmail)

	response.SuccessWithData(c, map[string]any{
		"flow_id": f.GetID().String(),
		"state":   verflow.StateSentEmail,
		"email":   email,
	})
	return flow.ErrCompletedByStrategy
}

// completeVerification handles Step 2: verify the code and mark the address.
func (s *Strategy) completeVerification(ctx context.Context, f *verflow.Flow, code string) error {
	addrPool, ok := s.d.(VerifiableAddressPool)
	if !ok {
		return errors.New("code: driver does not implement VerifiableAddressPool")
	}

	vc, err := s.d.VerificationCodePersister().FindVerificationCode(ctx, f.GetID(), code)
	if err != nil {
		return errors.New("invalid or expired code")
	}
	if vc.IsExpired() {
		return errors.New("code has expired, please request a new one")
	}
	if vc.IsUsed() {
		return errors.New("code has already been used")
	}

	// Consume the code immediately to prevent replay attacks.
	_ = s.d.VerificationCodePersister().UseVerificationCode(ctx, vc.ID)

	// Look up the identity that owns this email and mark the address verified.
	pool, ok := s.d.(VerificationIdentityPool)
	if !ok {
		return errors.New("code: driver does not implement VerificationIdentityPool")
	}
	i, err := pool.FindIdentityByEmail(ctx, vc.Address)
	if err != nil {
		// The identity may have been deleted between Step 1 and Step 2.
		return errors.New("identity not found")
	}

	addr, err := addrPool.FindOrCreateVerifiableAddress(ctx, i.ID, "email", vc.Address)
	if err != nil {
		return err
	}
	if err := addrPool.MarkVerifiableAddressVerified(ctx, addr.ID); err != nil {
		return err
	}

	return nil
}
