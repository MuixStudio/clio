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

package session

import (
	"context"

	"github.com/google/uuid"
)

type (
	SessionPersister interface {
		CreateSession(ctx context.Context, sess *Session) error
		UpdateSession(ctx context.Context, sess *Session) error
		GetSessionByToken(ctx context.Context, token string) (*Session, error)
		RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, oldToken, newToken string) error
		GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
		ListSessionsByIdentityID(ctx context.Context, identityID uuid.UUID) ([]*Session, error)
		RevokeSession(ctx context.Context, id uuid.UUID) error
		RevokeAllSessionsByIdentityID(ctx context.Context, identityID uuid.UUID, exceptID uuid.UUID) error
	}
	SessionPersisterProvider interface {
		SessionPersister() SessionPersister
	}

	SessionTokenExchangeCodePersister interface {
		CreateTokenExchangeCode(ctx context.Context, c *TokenExchangeCode) error
		ConsumeTokenExchangeCode(ctx context.Context, code string) (*TokenExchangeCode, error)
	}
	SessionTokenExchangeCodePersisterProvider interface {
		SessionTokenExchangeCodePersister() SessionTokenExchangeCodePersister
	}
)
