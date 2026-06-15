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

package verification

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/muixstudio/clio/internal/identity"
)

// Strategy is the interface every verification strategy must implement.
// Currently only the code strategy implements this, but the interface
// allows future strategies (e.g. link-based, magic link) to be plugged in.
type Strategy interface {
	// ID returns the credential type this strategy handles (e.g. "code").
	ID() identity.CredentialsType

	// Verify handles both steps of the verification flow:
	//   Step 1 — body contains {"email":"..."}: sends a code, writes response,
	//             returns ErrCompletedByStrategy.
	//   Step 2 — body contains {"code":"..."}: verifies the code, marks the
	//             address verified, returns nil (handler writes completion response).
	Verify(c *gin.Context, f *Flow, body json.RawMessage) error
}

// StrategyProvider is satisfied by the Driver.
type StrategyProvider interface {
	VerificationStrategies() []Strategy
}
