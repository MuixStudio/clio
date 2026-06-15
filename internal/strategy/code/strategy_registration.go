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
	regflow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
)

// Compile-time assertion: Strategy implements regflow.Strategy.
var _ regflow.Strategy = (*Strategy)(nil)

// registrationInitBody is the Step-1 payload: client provides an email address.
type registrationInitBody struct {
	Email string `json:"email"`
}

// registrationVerifyBody is the Step-2 payload: client submits the received code.
type registrationVerifyBody struct {
	Code string `json:"code"`
}

// Register implements regflow.Strategy — two-step passwordless registration.
//
// Step 1 — body contains {"email":"u@example.com"}:
//   - Generates and persists a 6-digit code.
//   - Delivers it via Courier (SMTP or console).
//   - Writes {"flow_id":"...","state":"sent_email"} and returns ErrCompletedByStrategy.
//     The handler stops; the flow remains pending in DB.
//
// Step 2 — body contains {"code":"123456"}, flow_id in the outer envelope:
//   - Validates code against the pending flow.
//   - Creates a new Identity whose credential identifier is the email.
//   - Returns the Identity so the handler can run post-hooks and issue a session.
func (s *Strategy) Register(c *gin.Context, f *regflow.Flow, body json.RawMessage) (*identity.Identity, error) {
	ctx := c.Request.Context()

	// Step 2: code present → complete registration.
	var vb registrationVerifyBody
	if err := json.Unmarshal(body, &vb); err == nil && vb.Code != "" {
		return s.completeRegistration(ctx, f, vb.Code)
	}

	// Step 1: email present → send code.
	var ib registrationInitBody
	if err := json.Unmarshal(body, &ib); err != nil || ib.Email == "" {
		// Neither field present — not our request; let the next strategy try.
		return nil, flow.ErrStrategyNotResponsible
	}

	if err := s.createAndSendCode(ctx, f.GetID(), "registration", ib.Email); err != nil {
		return nil, err
	}

	// Write response and signal the handler we're done for this request.
	// The flow stays pending in DB; the client must POST again (with flow_id +
	// code) to complete registration.
	response.SuccessWithData(c, map[string]string{
		"flow_id": f.GetID().String(),
		"state":   "sent_email",
	})
	return nil, flow.ErrCompletedByStrategy
}

// completeRegistration verifies the submitted code and creates the identity.
func (s *Strategy) completeRegistration(ctx context.Context, f *regflow.Flow, code string) (*identity.Identity, error) {
	vc, err := s.d.VerificationCodePersister().FindVerificationCode(ctx, f.GetID(), code)
	if err != nil {
		return nil, errors.New("invalid or expired code")
	}
	if vc.IsExpired() {
		return nil, errors.New("code has expired")
	}
	if vc.IsUsed() {
		return nil, errors.New("code has already been used")
	}

	// Idempotent: if an identity already exists for this email (e.g. double-submit),
	// just mark the code used and return the existing identity.
	existing, _, findErr := s.d.IdentityPersister().FindByCredentialsIdentifier(ctx, identity.CredentialsTypeCode, vc.Address)
	if findErr == nil {
		_ = s.d.VerificationCodePersister().UseVerificationCode(ctx, vc.ID)
		return existing, nil
	}

	// Create a new passwordless identity.
	// The email is stored both as a Trait and as the credential identifier.
	i := identity.New()
	i.Traits["email"] = vc.Address
	i.Credentials[identity.CredentialsTypeCode] = &identity.Credentials{
		Type:        identity.CredentialsTypeCode,
		Identifiers: []string{vc.Address},
		Config:      []byte(`{}`),
	}

	if err := s.d.IdentityPersister().CreateIdentity(ctx, i); err != nil {
		return nil, err
	}

	_ = s.d.VerificationCodePersister().UseVerificationCode(ctx, vc.ID)
	return i, nil
}
