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

// TokenExchangeCode is a short-lived, one-time token that a native/mobile app
// can exchange for the real session token.
//
// Flow:
//
//  1. Client initiates a login/registration flow.
//     The handler creates a flow as usual AND generates a TokenExchangeCode tied to
//     the session that will be issued upon completion.
//
//  2. The browser/webview completes the flow and receives the exchange code (not
//     the raw session token). The raw token is never exposed to the browser.
//
//  3. The native app calls POST /sessions/token-exchange with the code.
//     The handler validates the code, fetches the session, and returns the real
//     session token to the native app's secure storage.
//
// This pattern prevents the webview from ever seeing the session token, which is
// important on mobile where the webview may run in an isolated but less trusted context.
package session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const (
	// TokenExchangeCodeTTL is how long an exchange code is valid.
	// Short because the native app should claim it immediately after the browser flow.
	TokenExchangeCodeTTL = 2 * time.Minute
)

// TokenExchangeCode is a one-time code that can be redeemed for a session token.
type TokenExchangeCode struct {
	ID uuid.UUID `json:"id"`

	// Code is the opaque, random string given to the client.
	// Tagged json:"-" so it is never accidentally echoed in logs or responses.
	// It is returned exactly once (at creation) via the InitCode field below.
	Code string `json:"-"`

	// SessionID is the session that will be returned when this code is consumed.
	SessionID uuid.UUID `json:"session_id"`

	// ExpiresAt is when the code becomes invalid (typically now + 2 minutes).
	ExpiresAt time.Time `json:"expires_at"`

	// UsedAt is set when the code is consumed. A non-nil value means the code
	// has already been redeemed and must be rejected on any subsequent attempt.
	UsedAt *time.Time `json:"used_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// NewTokenExchangeCode generates a cryptographically random exchange code
// linked to the given session. Call after the session is created.
func NewTokenExchangeCode(sessionID uuid.UUID) (*TokenExchangeCode, error) {
	b := make([]byte, 16) // 128-bit entropy → 32 hex chars
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	now := time.Now()
	return &TokenExchangeCode{
		ID:        uuid.New(),
		Code:      hex.EncodeToString(b),
		SessionID: sessionID,
		ExpiresAt: now.Add(TokenExchangeCodeTTL),
		CreatedAt: now,
	}, nil
}

// IsExpired returns true if the code's TTL has elapsed.
func (c *TokenExchangeCode) IsExpired() bool { return time.Now().After(c.ExpiresAt) }

// IsUsed returns true if the code has already been redeemed.
func (c *TokenExchangeCode) IsUsed() bool { return c.UsedAt != nil }
