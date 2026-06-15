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

// Package recovery implements the account recovery self-service flow.
//
// Recovery lets a user regain access to their account when they have lost
// their credentials (e.g. forgotten password). The flow ends with a session
// being issued — the user can then use that session to change their password
// via the settings flow without knowing the old one.
//
// This is the key distinction from verification:
//   - Verification → proves address ownership, marks it verified, no session issued.
//   - Recovery     → proves address ownership, issues a session for credential reset.
//
// State machine:
//
//	choose_method → sent_email → passed_challenge
//	                          ↘ failed (bad/expired code)
//
// Endpoints:
//
//	POST /self-service/recovery           (Step 1 + Step 2)
//	GET  /self-service/recovery/flows?id= (inspect flow)
package recovery

import (
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/flow"
)

// Recovery-specific flow states, extending the base set.
const (
	// StateChooseMethod is the initial state — no email submitted yet.
	StateChooseMethod flow.State = "choose_method"
	// StateSentEmail means a code was sent; waiting for the user to submit it.
	StateSentEmail flow.State = "sent_email"
	// StatePassedChallenge means the code was verified and a session was issued.
	StatePassedChallenge flow.State = "passed_challenge"
)

// Flow is the recovery self-service flow.
type Flow struct {
	flow.Base

	// Email is the address the recovery code was sent to.
	// Set after Step 1 so Step 2 can reference it without client re-submission.
	Email string `json:"email,omitempty"`
}

func (f *Flow) GetFlowName() flow.FlowName { return flow.RecoveryFlow }

var _ flow.Flow = (*Flow)(nil)

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
