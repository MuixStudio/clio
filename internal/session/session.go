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

// Package session manages authenticated sessions.
//
// After a successful login or registration a Session is issued to the client
// as a bearer token (X-Session-Token or Authorization: Bearer <token>).
//
// Key concepts:
//
//   - AMR (Authentication Method References) — ordered list of credential
//     types and AAL levels used during this session. Each strategy appends
//     an entry via CompletedLoginFor / CompletedLoginForWithProvider.
//
//   - AAL (Authenticator Assurance Level) — derived from AMR:
//     aal1 = at least one first-factor was verified
//     aal2 = a second factor was also verified
//     The handler calls SetAuthenticatorAssuranceLevel() after all AMR entries
//     are recorded.
//
//   - AuthenticatedAt — when the last factor was verified (≥ IssuedAt).
//     Useful for privileged-session checks ("re-authenticate within 5 min").
package session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/identity"
)

// AuthenticationMethod records a single factor used during this session.
type AuthenticationMethod struct {
	// Method is the credential type (e.g. "password", "oidc", "code").
	Method identity.CredentialsType `json:"method"`
	// AAL is the assurance level this factor contributed.
	AAL identity.AuthenticatorAssuranceLevel `json:"aal"`
	// CompletedAt is when this factor was successfully verified.
	CompletedAt time.Time `json:"completed_at"`
	// Provider is set for OIDC logins (e.g. "github").
	Provider string `json:"provider,omitempty"`
}

// Session represents an authenticated user session.
type Session struct {
	ID uuid.UUID `json:"id"`

	OrganizationID uuid.NullUUID `json:"organization_id"`

	// Active is false when the session has been revoked (logout) or expired.
	// The row is kept rather than deleted to preserve the audit trail.
	Active bool `json:"active"`

	// RefreshToken is the long-lived opaque token stored in the database.
	// Rotated on every use; the old token is invalidated when a new one is issued.
	RefreshToken string `json:"refresh_token"`

	// AccessToken is a short-lived JWT returned to the client on issue/refresh.
	// Not persisted — generated on demand. omitempty so it is omitted in
	// whoami / list responses where it is not set.
	AccessToken string `json:"access_token,omitempty"`

	// ExpiresAt is when the session ceases to be valid.
	ExpiresAt time.Time `json:"expires_at"`

	// IssuedAt is when the session was first created.
	IssuedAt time.Time `json:"issued_at"`

	// AuthenticatedAt is when the last authentication factor was completed.
	// Updated by SetAuthenticatorAssuranceLevel() after all AMR entries are added.
	AuthenticatedAt time.Time `json:"authenticated_at"`

	// AuthenticatorAssuranceLevel is the highest AAL reached in this session.
	// Computed by SetAuthenticatorAssuranceLevel() from the AMR list.
	AuthenticatorAssuranceLevel identity.AuthenticatorAssuranceLevel `json:"authenticator_assurance_level"`

	// AMR is the ordered list of authentication methods used.
	AMR []AuthenticationMethod `json:"authentication_methods"`

	// IdentityID is a foreign-key helper — used by the persistence layer.
	IdentityID uuid.UUID `json:"identity_id"`

	// Identity is eagerly loaded on fetch so callers don't need a second query.
	Identity *identity.Identity `json:"identity,omitempty"`

	// CreatedAt / UpdatedAt are managed by the persistence layer.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GenerateRefreshToken returns a cryptographically random 64-hex-char opaque token.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func mustGenerateRefreshToken() string {
	token, err := GenerateRefreshToken()
	if err != nil {
		panic(err)
	}
	return token
}

// New creates a new active session for the given identity.
// AMR is not set here; call CompletedLoginFor / CompletedLoginForWithProvider
// and then SetAuthenticatorAssuranceLevel before persisting.
func New(i *identity.Identity, lifespan time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:                          uuid.New(),
		IdentityID:                  i.ID,
		RefreshToken:                mustGenerateRefreshToken(),
		Active:                      true,
		ExpiresAt:                   now.Add(lifespan),
		IssuedAt:                    now,
		AuthenticatedAt:             now,
		AuthenticatorAssuranceLevel: identity.NoAuthenticatorAssuranceLevel,
		AMR:                         []AuthenticationMethod{},
		Identity:                    i,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
}

// IsActive returns true when the session is active and has not expired.
func (s *Session) IsActive() bool {
	return s.Active && time.Now().Before(s.ExpiresAt)
}

// CompletedLoginFor appends an AMR entry for a standard (non-OIDC) method.
// Call once per authentication factor during a login or registration flow.
func (s *Session) CompletedLoginFor(method identity.CredentialsType, aal identity.AuthenticatorAssuranceLevel) {
	s.AMR = append(s.AMR, AuthenticationMethod{
		Method:      method,
		AAL:         aal,
		CompletedAt: time.Now(),
	})
}

// CompletedLoginForWithProvider is like CompletedLoginFor but also records
// the OIDC provider name. Use this in HandleCallback after OIDC login.
func (s *Session) CompletedLoginForWithProvider(method identity.CredentialsType, aal identity.AuthenticatorAssuranceLevel, provider string) {
	s.AMR = append(s.AMR, AuthenticationMethod{
		Method:      method,
		AAL:         aal,
		CompletedAt: time.Now(),
		Provider:    provider,
	})
}

// SetAuthenticatorAssuranceLevel computes the session AAL from its AMR list
// and updates AuthenticatedAt to the most recent CompletedAt timestamp.
//
// Call this after all AMR entries have been added (i.e. after all factors
// have been verified).
func (s *Session) SetAuthenticatorAssuranceLevel() {
	if len(s.AMR) == 0 {
		s.AuthenticatorAssuranceLevel = identity.NoAuthenticatorAssuranceLevel
		return
	}

	// Find the highest AAL and the latest completion timestamp.
	highest := identity.AuthenticatorAssuranceLevel1
	latestAt := time.Time{}
	for _, m := range s.AMR {
		if m.AAL == identity.AuthenticatorAssuranceLevel2 {
			highest = identity.AuthenticatorAssuranceLevel2
		}
		if m.CompletedAt.After(latestAt) {
			latestAt = m.CompletedAt
		}
	}

	s.AuthenticatorAssuranceLevel = highest
	if !latestAt.IsZero() {
		s.AuthenticatedAt = latestAt
	}
}

// AuthenticatedVia returns true if this session was authenticated using the
// given credential type. Useful for AAL upgrade checks.
func (s *Session) AuthenticatedVia(method identity.CredentialsType) bool {
	for _, m := range s.AMR {
		if m.Method == method {
			return true
		}
	}
	return false
}
