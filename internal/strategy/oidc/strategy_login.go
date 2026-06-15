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
	"net/http"

	"github.com/muixstudio/clio/internal/infra/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/muixstudio/clio/internal/flow"
	loginflow "github.com/muixstudio/clio/internal/flow/login"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
)

// Compile-time assertion: Strategy implements loginflow.Strategy.
var _ loginflow.Strategy = (*Strategy)(nil)

// loginBody is the "data" payload the client sends when selecting OIDC for login or registration.
type loginBody struct {
	Data struct {
		Provider string `json:"provider"`
	} `json:"data"`
}

// Login implements loginflow.Strategy.
//
// Unlike password login, this method does NOT authenticate the user directly.
// It generates a one-time CSRF state token (embedding the flow ID and kind),
// builds the provider's authorization URL, writes it to the response, and
// returns flow.ErrCompletedByStrategy so the login handler stops processing.
//
// Authentication is completed later in HandleCallback when the provider
// redirects the user back with an authorization code.
func (s *Strategy) Login(c *gin.Context, f *loginflow.Flow, _ *session.Session) (*identity.Identity, error) {
	var req loginBody
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		return nil, err
	}
	b := req.Data

	p, ok := s.providers[b.Provider]
	if !ok {
		return nil, errors.BadRequest("UNKNOWN_PROVIDER", "oidc: unknown provider "+b.Provider)
	}

	state := s.states.Generate(f.GetID(), flow.LoginFlow, b.Provider, c.Query("return_to"))

	c.JSON(http.StatusOK, gin.H{
		"redirect_url": p.AuthCodeURL(state),
	})
	return nil, nil
}
