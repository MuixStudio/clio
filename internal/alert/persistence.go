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

package alert

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ListOptions struct {
	Page, PageSize int
	Status         *Status
	Severity       *Severity
	TeamID         *uuid.UUID
	ConnectorID    *uuid.UUID
	StartAt        *time.Time // 筛选时间段的开始时间
	EndAt          *time.Time // 筛选时间段的结束时间
}

type (
	AlertPersister interface {
		Save(ctx context.Context, alert *NormalizedAlert) error
		SaveAlerts(ctx context.Context, alerts []*NormalizedAlert) error
		ListAlert(ctx context.Context, opts *ListOptions) (alerts []*NormalizedAlert, total int, err error)
	}
	AlertPersisterProvider interface {
		AlertPersister() AlertPersister
	}
)
