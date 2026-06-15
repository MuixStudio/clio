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

// Package code implements passwordless email-code login and registration.
//
// File layout:
//
//   - code.go                  — VerificationCode domain object
//   - strategy.go              — Strategy struct, deps interfaces, shared helpers
//   - strategy_registration.go — registration.Strategy implementation (two-step)
//   - strategy_login.go        — login.Strategy implementation (two-step)
//
// Two-step flow (same endpoint, two POST requests):
//
//	Step 1 — initiation
//	  POST /self-service/registration  {"method":"code","data":{"email":"u@example.com"}}
//	  → generates a 6-digit code, persists it, sends it via Courier
//	  → {"flow_id":"<uuid>","state":"sent_email"}
//	  → returns ErrCompletedByStrategy (handler stops, flow stays pending in DB)
//
//	Step 2 — completion
//	  POST /self-service/registration  {"method":"code","flow_id":"<uuid>","data":{"code":"123456"}}
//	  → validates code against the flow, creates identity, issues session
//	  → {"flow_id":"...","state":"completed","session":{...}}
//
// Login follows the identical two-step pattern.
package code

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/courier"
	"github.com/muixstudio/clio/internal/identity"
)

const codeLifespan = 10 * time.Minute

// dependencies is satisfied by the Driver.
type dependencies interface {
	VerificationCodePersisterProvider
	identity.IdentityPersisterProvider
	Courier() courier.Courier
}

// Strategy handles passwordless email-code registration and login.
// Its Register and Login methods live in strategy_registration.go and
// strategy_login.go respectively.
type Strategy struct{ d dependencies }

func New(d dependencies) *Strategy { return &Strategy{d: d} }

func (s *Strategy) ID() identity.CredentialsType { return identity.CredentialsTypeCode }

// =========================================================================
// Shared helpers
// =========================================================================

// generateCode returns a cryptographically random 6-digit string.
func generateCode() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("code: failed to generate random bytes: %w", err)
	}
	// Map 3 bytes (range 0–16 777 215) to 000000–999999.
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1_000_000
	return fmt.Sprintf("%06d", n), nil
}

// createAndSendCode generates a one-time code, persists it, and delivers it
// via the configured Courier. Called by both Register and Login initiation.
func (s *Strategy) createAndSendCode(ctx context.Context, flowID uuid.UUID, flowType, email string) error {
	code, err := generateCode()
	if err != nil {
		return err
	}

	vc := &VerificationCode{
		ID:        uuid.New(),
		FlowID:    flowID,
		FlowType:  flowType,
		Address:   email,
		Code:      code,
		ExpiresAt: time.Now().Add(codeLifespan),
	}
	if err := s.d.VerificationCodePersister().CreateVerificationCode(ctx, vc); err != nil {
		return fmt.Errorf("code: failed to persist verification code: %w", err)
	}
	if err := s.d.Courier().SendVerificationCode(ctx, email, code); err != nil {
		return fmt.Errorf("code: failed to send verification code: %w", err)
	}
	return nil
}
