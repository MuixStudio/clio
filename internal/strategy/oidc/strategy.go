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

// Package oidc implements a generic OAuth2 / OIDC strategy.
//
// File layout:
//
//   - provider.go          — Provider interface + normalized Claims struct
//   - state.go             — CSRF state store; each token carries flow context
//   - strategy.go          — Strategy struct, core wiring, HandleCallback, find-or-create
//   - strategy_login.go    — login.Strategy implementation (Login method)
//   - strategy_registration.go — registration.Strategy implementation (Register method)
//   - provider_github.go   — GitHub OAuth2 implementation of Provider
//
// Flow (API mode):
//
//  1. Client POSTs {"method":"oidc","data":{"provider":"github"}} to
//     /self-service/login or /self-service/registration.
//
//  2. Login() / Register() generates a state token (flow ID + kind + provider),
//     writes {"flow_id":"...","authorization_url":"..."}, returns ErrCompletedByStrategy.
//
//  3. Client opens the authorization URL. User authorizes at the provider.
//
//  4. Provider redirects to GET /auth/callback?code=xxx&state=xxx.
//     HandleCallback() validates state, exchanges code, finds-or-creates the
//     identity, runs flow hooks, issues session, returns {"session":{...}}.
package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/flow"
	loginflow "github.com/muixstudio/clio/internal/flow/login"
	regflow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/jwt"
	"github.com/muixstudio/clio/internal/session"
)

// dependencies is satisfied by the Driver.
type dependencies interface {
	config.Provider
	identity.IdentityPersisterProvider
	session.SessionPersisterProvider
	loginflow.FlowPersistenceProvider
	regflow.FlowPersistenceProvider

	RegistrationHookExecutor() *regflow.HookExecutor
	LoginHookExecutor() *loginflow.HookExecutor
}

// Strategy manages all OAuth2 / OIDC providers.
// Its Login and Register methods live in strategy_login.go and
// strategy_registration.go respectively.
type Strategy struct {
	d         dependencies
	providers map[string]Provider // keyed by Provider.ID()
	states    *stateStore
}

// New creates a Strategy and starts the background state-cleanup goroutine.
func New(d dependencies) *Strategy {
	s := &Strategy{
		d:         d,
		providers: make(map[string]Provider),
		states:    newStateStore(),
	}
	s.states.startCleanup()
	return s
}

// ID returns the credential type this strategy handles.
// Clients select it via {"method":"oidc",...}.
func (s *Strategy) ID() identity.CredentialsType { return identity.CredentialsTypeOIDC }

// AddProvider registers a provider. Call once per provider during setup.
func (s *Strategy) AddProvider(p Provider) { s.providers[p.ID()] = p }

// RegisterRoutes wires up the shared OAuth2 callback endpoint.
func (s *Strategy) RegisterRoutes(router *gin.Engine) {
	router.GET("/auth/callback", s.HandleCallback())
}

// =========================================================================
// Callback — shared by login and registration flows
// =========================================================================

// HandleCallback is the OAuth2 redirect target for every configured provider.
//
// Steps:
//  1. Validate CSRF state → recover flow context (ID, kind, provider)
//  2. Exchange code → Claims
//  3. Find-or-create identity
//  4. Issue session
//  5. Run post-hooks for the originating flow kind (login or registration)
//  6. Redirect to frontend with one-time exchange code
func (s *Strategy) HandleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 1. Validate CSRF state — one-time, expiring token.
		entry, err := s.states.Consume(c.Query("state"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_STATE", err.Error()))
			return
		}

		code := c.Query("code")
		if code == "" {
			response.Fail(c, errors.BadRequest("MISSING_CODE", "missing code parameter"))
			return
		}

		p, ok := s.providers[entry.Provider]
		if !ok {
			response.Fail(c, errors.BadRequest("UNKNOWN_PROVIDER", "oidc: unknown provider "+entry.Provider))
			return
		}

		// 2. Exchange code → normalized Claims (all provider-specific calls happen inside).
		claims, err := p.Exchange(ctx, code)
		if err != nil {
			response.Fail(c, errors.BadGateway("EXCHANGE_FAILED", err.Error()))
			return
		}

		// 3. Find-or-create — identical regardless of login vs registration.
		//    New user → create identity; existing user → return it.
		i, err := s.findOrCreate(ctx, entry.Provider, claims)
		if err != nil {
			response.Fail(c, errors.InternalServerError("IDENTITY_ERROR", err.Error()))
			return
		}

		// 4. Issue session with generic OIDC AMR, then enrich with provider info.
		sess, err := s.issueSession(ctx, i, identity.CredentialsTypeOIDC)
		if err != nil {
			response.Fail(c, errors.InternalServerError("SESSION_ISSUE_FAILED", err.Error()))
			return
		}
		// Replace generic AMR entry with provider-aware one.
		sess.AMR = nil
		sess.CompletedLoginForWithProvider(identity.CredentialsTypeOIDC, identity.AuthenticatorAssuranceLevel1, entry.Provider)
		sess.SetAuthenticatorAssuranceLevel()
		_ = s.d.SessionPersister().UpdateSession(ctx, sess)

		// 5. Finalize the originating flow and run its post-hook chain.
		switch entry.FlowKind {
		case flow.LoginFlow:
			if err := s.processLoginCallback(c, entry.FlowID, i, sess); err != nil {
				return
			}
		case flow.RegistrationFlow:
			if err := s.processRegistrationCallback(c, entry.FlowID, i, sess); err != nil {
				return
			}
		default:
			response.Fail(c, errors.InternalServerError("UNKNOWN_FLOW_KIND", fmt.Sprintf("unknown flow kind: %s", entry.FlowKind)))
			return
		}

		// 6. Write session cookies and redirect to the post-login destination.
		accessMaxAge := int(s.d.Config().Session.AccessTokenLifespan.Seconds())
		sessionMaxAge := int(s.d.Config().Session.Lifespan.Seconds())
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("session_id", sess.ID.String(), sessionMaxAge, "/", "", false, true)
		c.SetCookie("access_token", sess.AccessToken, accessMaxAge, "/", "", false, true)
		c.SetCookie("refresh_token", sess.RefreshToken, sessionMaxAge, "/", "", false, true)
		c.Redirect(http.StatusFound, "http://localhost:3000")
	}
}

func (s *Strategy) processLoginCallback(c *gin.Context, flowID uuid.UUID, i *identity.Identity, sess *session.Session) error {
	ctx := c.Request.Context()
	lf, err := s.d.LoginFlowPersister().GetLoginFlow(ctx, flowID)
	if err != nil {
		response.Fail(c, errors.InternalServerError("FLOW_NOT_FOUND", fmt.Sprintf("login flow not found: %s", err)))
		return err
	}
	if lf.IsExpired() {
		response.Fail(c, errors.Gone("FLOW_EXPIRED", "login flow expired"))
		return fmt.Errorf("login flow expired")
	}
	lf.Active = s.ID()
	lf.SetState(flow.StateCompleted)
	_ = s.d.LoginFlowPersister().UpdateLoginFlow(ctx, lf)
	if err := s.d.LoginHookExecutor().PostHook(ctx, lf, i, sess); err != nil {
		response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
		return err
	}
	return nil
}

func (s *Strategy) processRegistrationCallback(c *gin.Context, flowID uuid.UUID, i *identity.Identity, sess *session.Session) error {
	ctx := c.Request.Context()
	rf, err := s.d.RegistrationFlowPersister().GetRegistrationFlow(ctx, flowID)
	if err != nil {
		response.Fail(c, errors.InternalServerError("FLOW_NOT_FOUND", fmt.Sprintf("registration flow not found: %s", err)))
		return err
	}
	if rf.IsExpired() {
		response.Fail(c, errors.Gone("FLOW_EXPIRED", "registration flow expired"))
		return fmt.Errorf("registration flow expired")
	}
	rf.Active = s.ID()
	rf.SetState(flow.StateCompleted)
	_ = s.d.RegistrationFlowPersister().UpdateRegistrationFlow(ctx, rf)
	if err := s.d.RegistrationHookExecutor().PostHook(ctx, rf, i); err != nil {
		response.Fail(c, errors.InternalServerError("POST_HOOK_FAILED", err.Error()))
		return err
	}
	return nil
}

// =========================================================================
// Find-or-create identity (shared by login and registration callbacks)
// =========================================================================

// oidcCredConfig is persisted as Credentials.Config for any OIDC provider.
type oidcCredConfig struct {
	Provider string `json:"provider"` // e.g. "github"
	Subject  string `json:"subject"`  // provider's stable, immutable user ID
}

// findOrCreate looks up the identity by "<providerID>:<subject>".
// If none exists, a new identity is created from the normalized Claims.
// A single callback handles both first-time sign-up and subsequent logins transparently.
func (s *Strategy) findOrCreate(ctx context.Context, providerID string, claims *Claims) (*identity.Identity, error) {
	pool := s.d.IdentityPersister()

	// "<provider>:<subject>" is globally unique across providers.
	identifier := providerID + ":" + claims.Subject

	existing, _, err := pool.FindByCredentialsIdentifier(ctx, identity.CredentialsTypeOIDC, identifier)
	if err == nil {
		return existing, nil
	}

	cfg, _ := json.Marshal(oidcCredConfig{Provider: providerID, Subject: claims.Subject})

	i := identity.New()
	if claims.Email != "" {
		i.Traits["email"] = claims.Email
	}
	if claims.Name != "" {
		i.Traits["name"] = claims.Name
	}
	if claims.Nickname != "" {
		i.Traits["nickname"] = claims.Nickname
	}
	if claims.Picture != "" {
		i.Traits["picture"] = claims.Picture
	}
	i.Credentials[identity.CredentialsTypeOIDC] = &identity.Credentials{
		Type:        identity.CredentialsTypeOIDC,
		Identifiers: []string{identifier},
		Config:      cfg,
	}

	if err := pool.CreateIdentity(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Strategy) issueSession(ctx context.Context, i *identity.Identity, method identity.CredentialsType) (*session.Session, error) {
	sess := session.New(i, s.d.Config().Session.Lifespan)
	sess.CompletedLoginFor(method, identity.AuthenticatorAssuranceLevel1)
	sess.SetAuthenticatorAssuranceLevel()
	if err := s.d.SessionPersister().CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	at, err := jwt.GenerateAccessToken(sess, s.d.Config().Session.JWTSecret, s.d.Config().Session.AccessTokenLifespan)
	if err != nil {
		return nil, err
	}
	sess.AccessToken = at
	return sess, nil
}
