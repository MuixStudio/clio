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

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"

	regflow "github.com/muixstudio/clio/internal/flow/registration"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/errors"
)

// Compile-time assertion: Strategy implements regflow.Strategy.
var _ regflow.Strategy = (*Strategy)(nil)

type registrationBody struct {
	Data struct {
		Identifier string `json:"identifier" validate:"required,min=4,max=32,identifier"`
		Password   string `json:"password" validate:"required,min=8,max=64,password"`
	} `json:"data"`
}

// Register creates a new password-based identity from identifier and password.
// Returns an error if the identifier is already taken or the password is too short.
func (s *Strategy) Register(c *gin.Context, _ *regflow.Flow, body json.RawMessage) (*identity.Identity, error) {
	ctx := c.Request.Context()
	var b registrationBody
	if err := c.ShouldBindBodyWith(&b, binding.JSON); err != nil {
		return nil, err
	}

	if existing, _, err := s.d.IdentityPersister().FindByCredentialsIdentifier(ctx, identity.CredentialsTypePassword, b.Data.Identifier); err == nil && existing != nil {
		return nil, errors.Conflict("RECORD_ALREADY_EXISTS", "identity already exist")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(b.Data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	cfg, _ := json.Marshal(passwordConfig{HashedPassword: string(hashed)})

	i := identity.New()
	i.Traits["identifier"] = b.Data.Identifier
	i.Credentials[identity.CredentialsTypePassword] = &identity.Credentials{
		Type:        identity.CredentialsTypePassword,
		Identifiers: []string{b.Data.Identifier},
		Config:      cfg,
	}

	if err := s.d.IdentityPersister().CreateIdentity(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}
