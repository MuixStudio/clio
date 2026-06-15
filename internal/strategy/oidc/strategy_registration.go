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

package oidc

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/muixstudio/clio/internal/infra/response"

	"github.com/muixstudio/clio/internal/flow"
	regflow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
)

// Compile-time assertion: Strategy implements regflow.Strategy.
var _ regflow.Strategy = (*Strategy)(nil)

// Register implements regflow.Strategy.
//
// Like Login, it does not create an identity immediately. It builds the
// provider's authorization URL and returns flow.ErrCompletedByStrategy.
// The actual identity creation (find-or-create) happens in HandleCallback
// when the provider redirects back with an authorization code.
func (s *Strategy) Register(c *gin.Context, f *regflow.Flow, body json.RawMessage) (*identity.Identity, error) {
	var req loginBody
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}
	b := req.Data

	p, ok := s.providers[b.Provider]
	if !ok {
		return nil, errors.BadRequest("UNKNOWN_PROVIDER", "oidc: unknown provider "+b.Provider)
	}

	state := s.states.Generate(f.GetID(), flow.RegistrationFlow, b.Provider, c.Query("return_to"))

	response.SuccessWithData(c, map[string]string{
		"flow_id":           f.GetID().String(),
		"authorization_url": p.AuthCodeURL(state),
	})
	return nil, flow.ErrCompletedByStrategy
}
