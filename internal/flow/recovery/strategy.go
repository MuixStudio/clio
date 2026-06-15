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

package recovery

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/muixstudio/clio/internal/identity"
)

// Strategy is the interface every recovery strategy must implement.
// Currently only the code strategy implements this.
type Strategy interface {
	// ID returns the credential type this strategy handles (e.g. "code").
	ID() identity.CredentialsType

	// Recover handles both steps of the recovery flow.
	//
	// Step 1 — body contains {"email":"..."}: sends a code, writes the response
	//           itself, returns (nil, ErrCompletedByStrategy).
	//
	// Step 2 — body contains {"code":"..."}: verifies the code and returns the
	//           recovered identity. The handler then issues a session so the user
	//           can reset their password without knowing the old one.
	Recover(c *gin.Context, f *Flow, body json.RawMessage) (*identity.Identity, error)
}

// StrategyProvider is satisfied by the Driver.
type StrategyProvider interface {
	RecoveryStrategies() []Strategy
}
