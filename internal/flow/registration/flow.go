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
	"time"

	"github.com/muixstudio/clio/internal/flow"
	"github.com/muixstudio/clio/internal/identity"
)

// Flow is the registration-specific flow.
// It embeds flow.Base and adds fields relevant only to registration.
type Flow struct {
	flow.Base
	// Active records which strategy completed the registration.
	Active identity.CredentialsType `json:"active,omitempty"`
}

func (Flow) GetFlowName() flow.FlowName { return flow.RegistrationFlow }

// Compile-time check: *Flow must satisfy flow.Flow.
// If flow.Flow gains a new method and we forget to implement it here,
// the compiler will catch it immediately.
var _ flow.Flow = (*Flow)(nil)

func New(lifespan time.Duration) *Flow {
	return &Flow{Base: flow.NewBase(lifespan)}
}
