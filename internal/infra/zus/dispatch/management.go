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

package dispatch

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/alertmanager/pkg/labels"
)

// RouteWriter persists route mutations to a team's routing tree. The HTTP
// handler applies CRUD synchronously through it (then publishes a reload so the
// live dispatcher catches up). *gormstore.GormStore satisfies it.
type RouteWriter interface {
	AddAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in AddRouteInput) (uuid.UUID, error)
	UpdateAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in UpdateRouteInput) error
	DeleteAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, routeID uuid.UUID) error
}

// RouteWriterProvider exposes a RouteWriter from the driver.
type RouteWriterProvider interface {
	AlertRouteWriter() RouteWriter
}

// RouteReloader asks the dispatcher manager to rebuild a team's active
// dispatcher from the persisted routing tree. Implementations publish a reload
// onto the bus; the manager applies it asynchronously. It is the only route
// command that travels the bus — CRUD goes straight to the RouteWriter.
type RouteReloader interface {
	ReloadAlertRoutes(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) error
}

// RouteReloaderProvider exposes a RouteReloader from the driver.
type RouteReloaderProvider interface {
	AlertRouteReloader() RouteReloader
}

// RouteTreeLoader loads a team's effective alert route tree.
type RouteTreeLoader interface {
	LoadAlertRouteTree(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) (*Route, error)
}

// RouteTreeLoaderProvider exposes a RouteTreeLoader from the driver.
type RouteTreeLoaderProvider interface {
	AlertRouteLoader() RouteTreeLoader
}

// UpdateRouteInput describes mutable route fields. Nil fields are left
// unchanged; empty string/slice values are meaningful when the pointer is set.
type UpdateRouteInput struct {
	RouteID uuid.UUID

	Name       *string
	Matchers   *labels.Matchers
	Continue   *bool
	Receiver   *string
	GroupBy    *[]string
	GroupByAll *bool

	GroupWait           *time.Duration
	GroupInterval       *time.Duration
	RepeatInterval      *time.Duration
	MuteTimeIntervals   *[]string
	ActiveTimeIntervals *[]string
}

type AddRouteInput struct {
	// RouteID, when non-nil, is the id to assign the new route. The route-write
	// path generates it client-side before publishing the command so the caller
	// learns the id without waiting for the async apply. Nil lets the store
	// generate one.
	RouteID uuid.UUID

	ParentID uuid.NullUUID
	Name     string
	Matchers labels.Matchers
	Continue bool

	Receiver            string
	GroupBy             []string
	GroupByAll          bool
	GroupWait           *time.Duration
	GroupInterval       *time.Duration
	RepeatInterval      *time.Duration
	MuteTimeIntervals   []string
	ActiveTimeIntervals []string
}
