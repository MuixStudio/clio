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

// Package flow defines types and interfaces shared by all self-service flows
// (login, registration, recovery, ...).
//
// Design:
//   - Flow is a short-lived, persisted state machine.
//   - Base embeds shared fields; concrete flows embed Base.
//   - ErrStrategyNotResponsible lets a Strategy decline a request so the
//     next one gets a chance.
package flow

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ErrStrategyNotResponsible is returned by a Strategy when it cannot handle
// the current request. The Handler will try the next registered strategy.
var ErrStrategyNotResponsible = errors.New("strategy is not responsible for this request")

// ErrCompletedByStrategy is returned when a strategy has already written the
// HTTP response itself (e.g. OIDC returning an authorization URL) and the
// Handler must not write any further response.
var ErrCompletedByStrategy = errors.New("strategy completed the flow and wrote the response")

// FlowName identifies the kind of self-service flow.
type FlowName string

const (
	LoginFlow        FlowName = "login"
	RegistrationFlow FlowName = "registration"
	LogoutFlow       FlowName = "logout"
	VerificationFlow FlowName = "verification"
	RecoveryFlow     FlowName = "recovery"
)

// State is the current step of the flow state machine.
type State string

const (
	StatePending   State = "pending"
	StateCompleted State = "completed"
	StateFailed    State = "failed"
)

// Flow is the interface every concrete flow must satisfy.
// Handlers and hooks depend only on this interface, keeping them decoupled
// from specific flow types.
type Flow interface {
	GetID() uuid.UUID
	GetState() State
	SetState(State)
	GetFlowName() FlowName
	IsExpired() bool
}

// ExpiredError is returned when the client submits an expired flow.
type ExpiredError struct {
	FlowID    uuid.UUID
	ExpiredAt time.Time
}

func (e *ExpiredError) Error() string {
	return fmt.Sprintf("flow %s expired at %s", e.FlowID, e.ExpiredAt.Format(time.RFC3339))
}

// Base holds fields shared by every concrete Flow.
// Embed Base in a concrete flow struct to satisfy most of the Flow interface.
type Base struct {
	ID        uuid.UUID `json:"id"`
	State     State     `json:"state"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewBase(lifespan time.Duration) Base {
	return Base{
		ID:        uuid.New(),
		State:     StatePending,
		ExpiresAt: time.Now().Add(lifespan),
		CreatedAt: time.Now(),
	}
}

func (b *Base) GetID() uuid.UUID { return b.ID }
func (b *Base) GetState() State  { return b.State }
func (b *Base) SetState(s State) { b.State = s }
func (b *Base) IsExpired() bool  { return time.Now().After(b.ExpiresAt) }
