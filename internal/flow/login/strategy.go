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

package login

import (
	"github.com/gin-gonic/gin"

	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
)

// Strategy is the plugin interface for a login method.
//
// c is passed so that strategies like OIDC can write an authorization
// URL directly and signal completion via flow.ErrCompletedByStrategy.
// Simple strategies (password) just ignore c and use c.Request.Context().
//
// The sess parameter carries the caller's existing session when available
// (e.g. for a refresh or AAL-upgrade flow). Password login ignores it;
// a hypothetical TOTP strategy would use it to know whose second factor to verify.
type Strategy interface {
	ID() identity.CredentialsType

	// Login authenticates the request and returns the Identity on success.
	// Return flow.ErrStrategyNotResponsible to pass to the next strategy.
	// Return flow.ErrCompletedByStrategy if the response has already been written.
	Login(c *gin.Context, f *Flow, sess *session.Session) (*identity.Identity, error)
}

// StrategyProvider is implemented by the Driver.
type StrategyProvider interface {
	LoginStrategies() map[identity.CredentialsType]Strategy
}
