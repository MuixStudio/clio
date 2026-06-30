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

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Severity string
type Status string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"

	StatusFiring   Status = "firing"
	StatusResolved Status = "resolved"
)

type NormalizedAlert struct {
	// 身份
	ID          uuid.UUID
	Fingerprint string // 统一计算
	ConnectorID uuid.UUID

	TeamID uuid.UUID
	// 状态
	Status   Status   // 统一为 firing / resolved
	Severity Severity // 统一为 critical / high / medium / low

	// 标签（用于计算 fingerprint 和路由）
	Labels map[string]string

	// 描述（不参与 fingerprint）
	Annotations map[string]string

	// 时间
	StartsAt time.Time
	EndsAt   *time.Time // nil = 仍在 firing

	// 原始数据（存入 alert_events.raw_payload）
	RawPayload []byte

	UpdatedAt time.Time
	CreatedAt time.Time
}
