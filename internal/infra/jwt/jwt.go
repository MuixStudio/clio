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

package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/session"
)

// AccessTokenClaims are the JWT claims embedded in every Access Token.
type AccessTokenClaims struct {
	SessionID      uuid.UUID `json:"session_id"`
	IdentityID     uuid.UUID `json:"identity_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AAL            string    `json:"aal"`
	AMR            []string  `json:"amr"`
	jwt.RegisteredClaims
}

// GenerateAccessToken signs a new JWT for the given session.
// The token contains compact session facts for downstream authorization.
func GenerateAccessToken(sess *session.Session, secret string, lifespan time.Duration) (string, error) {
	now := time.Now()
	amr := make([]string, 0, len(sess.AMR))
	for _, m := range sess.AMR {
		amr = append(amr, string(m.Method))
	}
	claims := AccessTokenClaims{
		SessionID:      sess.ID,
		IdentityID:     sess.IdentityID,
		OrganizationID: sess.OrganizationID.UUID,
		AAL:            string(sess.AuthenticatorAssuranceLevel),
		AMR:            amr,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sess.IdentityID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(lifespan)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken validates the JWT signature and expiry, returning the claims.
func ParseAccessToken(tokenStr, secret string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
