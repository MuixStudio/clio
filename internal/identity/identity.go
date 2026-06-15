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

// Package identity defines the core Identity domain object.
// An Identity represents a person or entity that can authenticate.
package identity

import (
	"time"

	"github.com/google/uuid"
)

type CredentialsType string

func (c CredentialsType) String() string { return string(c) }

const (
	// CredentialsTypePassword is the standard identifier + bcrypt-hashed password strategy.
	CredentialsTypePassword CredentialsType = "password"
	// CredentialsTypeOIDC covers all OAuth2 / OIDC providers (GitHub, Google, …).
	// The provider is encoded in the credential identifier: "<provider>:<subject>".
	CredentialsTypeOIDC CredentialsType = "oidc"
	// CredentialsTypeCode is used by the passwordless email-code strategy.
	// The credential identifier is the email address.
	CredentialsTypeCode CredentialsType = "code"
)

// AuthenticatorAssuranceLevel (AAL) describes how strongly a session was
// authenticated.
//
//   - NoAAL  — session not yet authenticated (transient state)
//   - AAL1   — at least one first-factor credential was verified
//   - AAL2   — a second factor (e.g. TOTP) was also verified
type AuthenticatorAssuranceLevel string

const (
	NoAuthenticatorAssuranceLevel AuthenticatorAssuranceLevel = ""
	AuthenticatorAssuranceLevel1  AuthenticatorAssuranceLevel = "aal1"
	AuthenticatorAssuranceLevel2  AuthenticatorAssuranceLevel = "aal2"
)

// VerifiableAddressStatus tracks where a verifiable address is in its lifecycle.
type VerifiableAddressStatus string

const (
	// VerifiableAddressStatusPending means no code has been sent yet.
	VerifiableAddressStatusPending VerifiableAddressStatus = "pending"
	// VerifiableAddressStatusSent means a code was sent but not yet confirmed.
	VerifiableAddressStatusSent VerifiableAddressStatus = "sent"
	// VerifiableAddressStatusCompleted means the address has been verified.
	VerifiableAddressStatusCompleted VerifiableAddressStatus = "completed"
)

// VerifiableAddress represents a contact address (email/phone) that can be
// cryptographically proven to belong to an identity.
type VerifiableAddress struct {
	ID uuid.UUID `json:"id"`

	// Value is the address itself (e.g. "user@example.com").
	Value string `json:"value"`

	// Via is the delivery channel: "email" (sms would be added later).
	Via string `json:"via"`

	// Verified is true once the owner has proven control of this address.
	Verified bool `json:"verified"`

	// VerifiedAt records when verification completed. Nil until then.
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Status tracks the address through the verification lifecycle.
	Status VerifiableAddressStatus `json:"status"`

	// IdentityID is the owning identity's ID (foreign key).
	IdentityID uuid.UUID `json:"identity_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Credentials holds one authentication method's data attached to an Identity.
// Each strategy (password, oidc, ...) stores its own config here.
type Credentials struct {
	Type        CredentialsType `json:"type"`
	Identifiers []string        `json:"identifiers"` // e.g. ["user@example.com"]
	Config      []byte          `json:"config"`      // strategy-specific JSON (e.g. {"hashed_password":"..."})
}

// Identity represents a person or entity that can authenticate.
// Credentials are keyed by CredentialsType and never exposed over the API.
type Identity struct {
	ID          uuid.UUID                        `json:"id"`
	Traits      map[string]any                   `json:"traits"`
	Credentials map[CredentialsType]*Credentials `json:"-"` // never expose raw credentials over the API
	// VerifiableAddresses is eagerly loaded when available.
	// Empty slice (not nil) when no addresses have been registered.
	VerifiableAddresses []VerifiableAddress `json:"verifiable_addresses,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

// New returns an empty Identity with a fresh ID, ready for credential attachment.
func New() *Identity {
	now := time.Now()
	return &Identity{
		ID:          uuid.New(),
		Traits:      make(map[string]any),
		Credentials: make(map[CredentialsType]*Credentials),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
