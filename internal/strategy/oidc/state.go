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

package oidc

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/flow"
)

const (
	stateLifespan   = 10 * time.Minute
	cleanupInterval = 5 * time.Minute
)

// stateEntry is the payload stored for each CSRF state token.
// It carries everything HandleCallback needs to reconstruct the flow context
// without hitting the database just to decode the state.
type stateEntry struct {
	FlowID   uuid.UUID
	FlowKind flow.FlowName // flow.LoginFlow or flow.RegistrationFlow
	Provider string        // provider ID, e.g. "github"
	ReturnTo string        // post-login redirect target; "/" if not provided
}

// storedEntry wraps stateEntry with an expiry timestamp for cleanup.
type storedEntry struct {
	stateEntry
	ExpiresAt time.Time
}

// stateStore issues and validates short-lived CSRF state tokens.
// The zero value is not usable; construct via newStateStore.
type stateStore struct {
	m sync.Map // state string → storedEntry
}

func newStateStore() *stateStore { return &stateStore{} }

// Generate mints a new state token that encodes the flow context.
// returnTo is the URL HandleCallback will redirect to after a successful login;
// it defaults to "/" when empty.
// The token is returned and must be sent to the provider as the OAuth2 state param.
func (ss *stateStore) Generate(flowID uuid.UUID, flowKind flow.FlowName, provider, returnTo string) string {
	if returnTo == "" {
		returnTo = "/"
	}
	state := uuid.New().String()
	ss.m.Store(state, storedEntry{
		stateEntry: stateEntry{FlowID: flowID, FlowKind: flowKind, Provider: provider, ReturnTo: returnTo},
		ExpiresAt:  time.Now().Add(stateLifespan),
	})
	return state
}

// Consume validates and atomically deletes the state token.
// Returns the embedded stateEntry on success, or an error if the token is
// unknown, already used, or expired.
func (ss *stateStore) Consume(state string) (stateEntry, error) {
	val, ok := ss.m.LoadAndDelete(state)
	if !ok {
		return stateEntry{}, errors.New("invalid or already-used state token")
	}
	entry := val.(storedEntry)
	if time.Now().After(entry.ExpiresAt) {
		return stateEntry{}, errors.New("state token expired — restart the OAuth flow")
	}
	return entry.stateEntry, nil
}

// startCleanup runs a background goroutine that evicts expired tokens.
// Call once after construction; it runs until the process exits.
func (ss *stateStore) startCleanup() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			ss.m.Range(func(key, val any) bool {
				if now.After(val.(storedEntry).ExpiresAt) {
					ss.m.Delete(key)
				}
				return true
			})
		}
	}()
}
