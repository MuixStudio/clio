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
	"time"

	"github.com/muixstudio/clio/internal/flow"
	"github.com/muixstudio/clio/internal/identity"
)

// Flow represents a single login attempt.
type Flow struct {
	flow.Base
	// Active records which strategy completed the login.
	Active identity.CredentialsType `json:"active,omitempty"`
	// Refresh is true when the caller explicitly requested a forced re-auth
	// (e.g. before a sensitive operation), even if a session already exists.
	Refresh bool `json:"refresh"`
}

func (Flow) GetFlowName() flow.FlowName { return flow.LoginFlow }

// Compile-time check.
var _ flow.Flow = (*Flow)(nil)

func New(lifespan time.Duration, refresh bool) *Flow {
	return &Flow{
		Base:    flow.NewBase(lifespan),
		Refresh: refresh,
	}
}
