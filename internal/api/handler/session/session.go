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
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/errors"
	"github.com/muixstudio/clio/internal/infra/jwt"
	"github.com/muixstudio/clio/internal/infra/response"
	"github.com/muixstudio/clio/internal/session"
)

func setSessionCookies(c *gin.Context, sess *session.Session) {
	maxAge := int(time.Until(sess.ExpiresAt).Seconds())
	for _, cookie := range []http.Cookie{
		{Name: "sid", Value: sess.ID.String()},
		{Name: "access_token", Value: sess.AccessToken},
		{Name: "refresh_token", Value: sess.RefreshToken},
	} {
		cookie.MaxAge = maxAge
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteLaxMode
		http.SetCookie(c.Writer, &cookie)
	}
}

func (h *Handler) bearerToken(c *gin.Context) string {
	parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return parts[1]
	}
	return ""
}

func (h *Handler) requireSession(c *gin.Context) *session.Session {
	tokenStr := h.bearerToken(c)
	if tokenStr == "" {
		tokenStr = c.GetHeader("X-Access-Token")
	}
	if tokenStr == "" {
		response.Fail(c, errors.Unauthorized("MISSING_TOKEN", "no access token provided — set X-Access-Token or Authorization: Bearer"))
		return nil
	}

	claims, err := jwt.ParseAccessToken(tokenStr, h.d.Config().Session.JWTSecret)
	if err != nil {
		response.Fail(c, errors.Unauthorized("INVALID_TOKEN", "invalid or expired access token"))
		return nil
	}

	sess, err := h.d.SessionPersister().GetSessionByID(c.Request.Context(), claims.SessionID)
	if err != nil || !sess.IsActive() {
		response.Fail(c, errors.Unauthorized("SESSION_REVOKED", "session not found or revoked"))
		return nil
	}
	return sess
}

func (h *Handler) WhoAmI() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := h.requireSession(c)
		if sess == nil {
			return
		}
		response.SuccessWithData(c, sess)
	}
}

// ListSessions returns all active sessions for the authenticated identity.
func (h *Handler) ListSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := h.requireSession(c)
		if sess == nil {
			return
		}
		sessions, err := h.d.SessionPersister().ListSessionsByIdentityID(c.Request.Context(), sess.IdentityID)
		if err != nil {
			response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
			return
		}
		if sessions == nil {
			sessions = []*session.Session{}
		}
		response.SuccessWithData(c, sessions)
	}
}

// RevokeOtherSessions invalidates every session for the caller's identity
// except the one used to make this request.
func (h *Handler) RevokeOtherSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := h.requireSession(c)
		if sess == nil {
			return
		}
		if err := h.d.SessionPersister().RevokeAllSessionsByIdentityID(c.Request.Context(), sess.IdentityID, sess.ID); err != nil {
			response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// RevokeSessionByID invalidates the session with the given ID, provided it
// belongs to the caller's identity.
//
// Returns 403 if the session belongs to a different identity (prevents
// cross-identity revocation via brute-forced UUIDs).
func (h *Handler) RevokeSessionByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		caller := h.requireSession(c)
		if caller == nil {
			return
		}

		targetID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.Fail(c, errors.BadRequest("INVALID_ARGUMENT", "invalid session id"))
			return
		}

		target, err := h.d.SessionPersister().GetSessionByID(c.Request.Context(), targetID)
		if err != nil {
			response.Fail(c, errors.NotFound("SESSION_NOT_FOUND", "session not found"))
			return
		}
		if target.IdentityID != caller.IdentityID {
			response.Fail(c, errors.Forbidden("FORBIDDEN", "session does not belong to you"))
			return
		}

		if err := h.d.SessionPersister().RevokeSession(c.Request.Context(), targetID); err != nil {
			response.Fail(c, errors.InternalServerError("DB_ERROR", err.Error()))
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// TokenExchange redeems a one-time exchange code for the real session token.
//
// Request body:  {"code": "<exchange_code>"}
// Response body: {"message":"ok","data":{"session":{...}}}
func (h *Handler) TokenExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code string `json:"code"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errors.BadRequest("INVALID_BODY", err.Error()))
			return
		}
		if req.Code == "" {
			response.Fail(c, errors.BadRequest("MISSING_CODE", "code is required"))
			return
		}

		ctx := c.Request.Context()

		ec, err := h.d.SessionTokenExchangeCodePersister().ConsumeTokenExchangeCode(ctx, req.Code)
		if err != nil {
			response.Fail(c, errors.Unauthorized("INVALID_CODE", "invalid or expired exchange code"))
			return
		}

		sess, err := h.d.SessionPersister().GetSessionByID(ctx, ec.SessionID)
		if err != nil {
			response.Fail(c, errors.InternalServerError("SESSION_NOT_FOUND", "session not found"))
			return
		}
		if !sess.IsActive() {
			response.Fail(c, errors.Unauthorized("SESSION_INACTIVE", "session is no longer active"))
			return
		}

		at, err := jwt.GenerateAccessToken(sess, h.d.Config().Session.JWTSecret, h.d.Config().Session.AccessTokenLifespan)
		if err != nil {
			response.Fail(c, errors.InternalServerError("TOKEN_GEN_FAILED", err.Error()))
			return
		}
		sess.AccessToken = at
		setSessionCookies(c, sess)
		response.SuccessWithData(c, map[string]any{"session": sess})
	}
}

// Refresh validates and rotates a Refresh Token, returning fresh tokens.
//
// Request:  refresh_token cookie
// Response: {"message":"ok","data":{"session":{...}}}
func (h *Handler) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		rtToken, err := c.Cookie("refresh_token")
		if err != nil {
			response.Fail(c, errors.BadRequest("MISSING_COOKIE", "refresh_token cookie is required"))
			return
		}

		sess, err := h.d.SessionPersister().GetSessionByToken(c.Request.Context(), rtToken)
		if err != nil {
			response.Fail(c, errors.Unauthorized("INVALID_TOKEN", "invalid or expired refresh token"))
			return
		}
		if !sess.IsActive() {
			response.Fail(c, errors.Unauthorized("SESSION_INACTIVE", "session is no longer active"))
			return
		}

		newRefreshToken, err := GenerateRefreshToken()
		if err != nil {
			response.Fail(c, errors.InternalServerError("TOKEN_GEN_FAILED", err.Error()))
			return
		}
		if err := h.d.SessionPersister().RotateRefreshToken(c.Request.Context(), sess.ID, rtToken, newRefreshToken); err != nil {
			response.Fail(c, errors.Unauthorized("INVALID_TOKEN", "invalid or expired refresh token"))
			return
		}

		sess.RefreshToken = newRefreshToken

		at, err := jwt.GenerateAccessToken(sess, h.d.Config().Session.JWTSecret, h.d.Config().Session.AccessTokenLifespan)
		if err != nil {
			response.Fail(c, errors.InternalServerError("TOKEN_GEN_FAILED", err.Error()))
			return
		}
		sess.AccessToken = at
		setSessionCookies(c, sess)
		response.SuccessWithData(c, map[string]any{"session": sess})
	}
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
