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

// Package verification implements the address verification self-service flow.
//
// A verification flow proves that an identity owns a contact address
// (email or phone number). The flow has two steps:
//
//  1. Initiate — client submits an email address; server sends a one-time code.
//  2. Complete — client submits the code; server marks the address as verified.
//
// The flow does NOT require the user to be authenticated. Anyone can initiate
// verification for an address. To prevent email enumeration, the server always
// returns the same response regardless of whether the address is registered.
//
// State machine:
//
//	choose_method → sent_email → passed_challenge
//	                          ↘ failed (bad/expired code)
package verification

import (
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/flow"
)

// Verification-specific flow states, extending the base set.
const (
	// StateChooseMethod is the initial state — no email submitted yet.
	StateChooseMethod flow.State = "choose_method"
	// StateSentEmail means a code was sent to the address. Waiting for submission.
	StateSentEmail flow.State = "sent_email"
	// StatePassedChallenge means the code was verified and the address is now marked verified.
	StatePassedChallenge flow.State = "passed_challenge"
)

// Flow is the verification self-service flow.
// Embeds flow.Base for ID, State, ExpiresAt, CreatedAt.
type Flow struct {
	flow.Base

	// Email is the address being verified. Set after Step 1 (initiate).
	// Stored so Step 2 (complete) can look up the right VerifiableAddress
	// without requiring the client to re-submit the email.
	Email string `json:"email,omitempty"`
}

// GetFlowName satisfies flow.Flow.
func (f *Flow) GetFlowName() flow.FlowName { return flow.VerificationFlow }

// Compile-time assertion.
var _ flow.Flow = (*Flow)(nil)

// New creates a fresh verification flow.
func New(lifespan time.Duration) *Flow {
	now := time.Now()
	return &Flow{
		Base: flow.Base{
			ID:        uuid.New(),
			State:     StateChooseMethod,
			ExpiresAt: now.Add(lifespan),
			CreatedAt: now,
		},
	}
}
