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

// Package password implements the password strategy for both login and registration.
//
// File layout:
//
//   - strategy.go              — IdentityPool, dependencies, Strategy struct, New(), ID(), passwordConfig
//   - strategy_login.go        — login.Strategy implementation (Login method)
//   - strategy_registration.go — registration.Strategy implementation (Register method)
package password

import (
	"github.com/muixstudio/clio/internal/identity"
)

type dependencies interface {
	identity.IdentityPersisterProvider
}

// Strategy handles password-based registration and login.
type Strategy struct{ d dependencies }

func New(d dependencies) *Strategy { return &Strategy{d: d} }

func (s *Strategy) ID() identity.CredentialsType { return identity.CredentialsTypePassword }

// passwordConfig is stored as Credentials.Config — never sent over the wire.
// Shared by Register (writes it) and Login (reads it).
type passwordConfig struct {
	HashedPassword string `json:"hashed_password"`
}
