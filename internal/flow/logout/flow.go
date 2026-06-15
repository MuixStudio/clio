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

package logout

import (
	"time"

	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/flow"
)

// Flow represents a single logout attempt.
// It is created when the client initiates logout and updated to completed
// after the session is revoked. Persisting it enables post-hooks (audit log,
// analytics, etc.).
type Flow struct {
	flow.Base
	// IdentityID is the identity being logged out.
	// Available to post-hooks so they know who logged out without having
	// to re-fetch the (now-revoked) session.
	IdentityID uuid.UUID `json:"identity_id"`
}

func (Flow) GetFlowName() flow.FlowName { return flow.LogoutFlow }

// Compile-time check.
var _ flow.Flow = (*Flow)(nil)

func New(lifespan time.Duration, identityID uuid.UUID) *Flow {
	return &Flow{
		Base:       flow.NewBase(lifespan),
		IdentityID: identityID,
	}
}
