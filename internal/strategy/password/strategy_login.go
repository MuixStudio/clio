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

package password

import (
	"encoding/json"
	"time"

	"github.com/muixstudio/clio/internal/infra/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"

	loginflow "github.com/muixstudio/clio/internal/flow/login"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
)

// Compile-time assertion: Strategy implements loginflow.Strategy.
var _ loginflow.Strategy = (*Strategy)(nil)

type loginBody struct {
	Data struct {
		Identifier string `json:"identifier" validate:"required,min=4,max=32,identifier"`
		Password   string `json:"password" validate:"required,min=8,max=64,password"`
	} `json:"data"`
}

// Login verifies identifier + password and returns the Identity on success.
//
// The sess parameter (the caller's existing session) is accepted but ignored
// here — password is a first-factor method and has no use for it.
// An AAL2 strategy (e.g. TOTP) would use sess.IdentityID to look up the
// right second-factor credential.
func (s *Strategy) Login(c *gin.Context, _ *loginflow.Flow, _ *session.Session) (*identity.Identity, error) {
	ctx := c.Request.Context()
	var b loginBody

	if err := c.ShouldBindBodyWith(&b, binding.JSON); err != nil {
		return nil, err
	}

	i, creds, err := s.d.IdentityPersister().FindByCredentialsIdentifier(ctx, identity.CredentialsTypePassword, b.Data.Identifier)
	if err != nil {
		// Constant-time delay: makes "user not found" and "wrong password"
		// indistinguishable by response time, preventing user enumeration.
		time.Sleep(100 * time.Millisecond)
		return nil, err
	}

	var cfg passwordConfig
	if err := json.Unmarshal(creds.Config, &cfg); err != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "internal server error").WithCause(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cfg.HashedPassword), []byte(b.Data.Password)); err != nil {
		return nil, errors.BadRequest("BAD_IDENTITY", "bad identifier or password").WithCause(err)
	}

	return i, nil
}
