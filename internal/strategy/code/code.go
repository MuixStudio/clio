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

package code

import (
	"time"

	"github.com/google/uuid"
)

// VerificationCode is a one-time, expiring token sent to a user's email address.
//
// It bridges two separate HTTP requests — the initiation step (send code) and
// the completion step (submit code) — by tying the code to the flow that
// created it.
type VerificationCode struct {
	ID       uuid.UUID
	FlowID   uuid.UUID
	FlowType string    // "login" | "registration"
	Address  string    // email address the code was sent to
	Code     string    // 6-digit numeric string, e.g. "047382"
	ExpiresAt time.Time
	UsedAt   *time.Time // nil until the code is consumed
}

// IsExpired reports whether the code's TTL has passed.
func (c *VerificationCode) IsExpired() bool { return time.Now().After(c.ExpiresAt) }

// IsUsed reports whether the code has already been consumed.
func (c *VerificationCode) IsUsed() bool { return c.UsedAt != nil }
