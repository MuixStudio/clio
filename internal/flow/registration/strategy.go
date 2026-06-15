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

package registration

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/muixstudio/clio/internal/identity"
)

// Strategy is the plugin interface for a registration method (password, oidc, ...).
//
// c is passed so that strategies like OIDC can write an authorization
// URL directly and signal completion via flow.ErrCompletedByStrategy.
//
// Adding a new auth method means implementing this interface and registering it
// in the Driver — no changes to the Handler are needed.
type Strategy interface {
	// ID returns the credential type this strategy handles (e.g. "password", "oidc").
	ID() identity.CredentialsType

	// Register creates the Identity and attaches credentials.
	// body is the raw JSON from the request's "data" field.
	// Return flow.ErrStrategyNotResponsible to skip to the next strategy.
	// Return flow.ErrCompletedByStrategy if the response has already been written.
	Register(c *gin.Context, f *Flow, body json.RawMessage) (*identity.Identity, error)
}

// StrategyProvider is implemented by the Driver to supply the active strategies.
type StrategyProvider interface {
	RegistrationStrategies() []Strategy
}
