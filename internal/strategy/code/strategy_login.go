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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/muixstudio/clio/internal/flow"
	loginflow "github.com/muixstudio/clio/internal/flow/login"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/session"
)

// Compile-time assertion: Strategy implements loginflow.Strategy.
var _ loginflow.Strategy = (*Strategy)(nil)

type loginBody struct {
	Data struct {
		Identifier string `json:"identifier"`
		Code       string `json:"code"`
	} `json:"data"`
}

// Login implements loginflow.Strategy — two-step passwordless login.
//
// Step 1 — body contains {"data":{"identifier":"u@example.com"}}:
//   - Verifies the account exists (prevents sending codes to unknown addresses).
//   - Generates and persists a 6-digit code.
//   - Delivers it via Courier.
//   - Writes {"flow_id":"...","state":"sent_email"} and returns ErrCompletedByStrategy.
//
// Step 2 — body contains {"data":{"code":"123456"}}, flow_id in the outer envelope:
//   - Looks up the code by flow_id + code value.
//   - Finds the identity whose email matches the code's stored address.
//   - Returns the Identity so the handler can issue a session.
//
// The sess parameter (existing session) is accepted but unused — code is a
// first-factor method. An AAL2 strategy (e.g. TOTP) would use it.
func (s *Strategy) Login(c *gin.Context, f *loginflow.Flow, _ *session.Session) (*identity.Identity, error) {
	ctx := c.Request.Context()
	var b loginBody
	if err := c.ShouldBindBodyWith(&b, binding.JSON); err != nil {
		return nil, err
	}

	// Step 2: code present → complete login.
	if b.Data.Code != "" {
		return s.completeLogin(ctx, f, b.Data.Code)
	}

	// Step 1: identifier present → send code.
	if b.Data.Identifier == "" {
		return nil, flow.ErrStrategyNotResponsible
	}

	// Verify the account exists before generating a code.
	if _, _, err := s.d.IdentityPersister().FindByCredentialsIdentifier(ctx, identity.CredentialsTypeCode, b.Data.Identifier); err != nil {
		return nil, errors.NotFound("RECORD_NOT_FOUND", "no account found for this email address").WithCause(err)
	}

	if err := s.createAndSendCode(ctx, f.GetID(), "login", b.Data.Identifier); err != nil {
		return nil, err
	}

	response.SuccessWithData(c, map[string]string{
		"flow_id": f.GetID().String(),
		"state":   "sent_email",
	})
	return nil, flow.ErrCompletedByStrategy
}

// completeLogin verifies the submitted code and returns the matching identity.
func (s *Strategy) completeLogin(ctx context.Context, f *loginflow.Flow, code string) (*identity.Identity, error) {
	vc, err := s.d.VerificationCodePersister().FindVerificationCode(ctx, f.GetID(), code)
	if err != nil {
		return nil, errors.BadRequest("INVALID_CODE", "invalid or expired code").WithCause(err)
	}
	if vc.IsExpired() {
		return nil, errors.Gone("CODE_EXPIRED", "code has expired")
	}
	if vc.IsUsed() {
		return nil, errors.BadRequest("CODE_ALREADY_USED", "code has already been used")
	}

	i, _, err := s.d.IdentityPersister().FindByCredentialsIdentifier(ctx, identity.CredentialsTypeCode, vc.Address)
	if err != nil {
		return nil, err
	}

	_ = s.d.VerificationCodePersister().UseVerificationCode(ctx, vc.ID)
	return i, nil
}
