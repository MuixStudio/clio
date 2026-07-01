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
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
	"github.com/muixstudio/clio/internal/infra/orgctx"
)

// publisher.go is the write side of the dispatch pipeline: it adapts a generic
// message.Publisher into the alert/route codecs the consumers expect. It depends
// only on the transport role (message.Publisher), so the bus package stays a
// pure transport that knows nothing about these envelopes.

// AlertPublisher marshals normalized alerts and tags each message with its
// org/team scope. It implements alert.AlertPublisher.
type AlertPublisher struct {
	pub message.Publisher
}

// NewAlertPublisher wraps pub so callers can publish normalized alerts.
func NewAlertPublisher(pub message.Publisher) *AlertPublisher {
	return &AlertPublisher{pub: pub}
}

// PublishAlerts marshals each alert and publishes it to TopicAlerts, tagging
// every message with the organization (from ctx) and the alert's team so the
// consumer can scope route loading. Alerts are published individually so the
// dispatcher groups them just as it would a live stream.
func (p *AlertPublisher) PublishAlerts(ctx context.Context, alerts []*entity.NormalizedAlert) error {
	orgID, _ := orgctx.OrgIDFromCtx(ctx)

	msgs := make([]*message.Message, 0, len(alerts))
	for _, a := range alerts {
		payload, err := json.Marshal(a)
		if err != nil {
			return err
		}
		msg := message.NewMessage(uuid.NewString(), payload)
		setScope(msg, orgID, a.TeamID)
		msgs = append(msgs, msg)
	}
	return p.pub.Publish(TopicAlerts, msgs...)
}

// ReloadPublisher publishes a route-reload request to TopicRouteReload for the
// Manager to apply asynchronously. Route CRUD is performed synchronously by the
// HTTP handler against the store; the handler then publishes a reload so the
// live Dispatcher is rebuilt from the freshly persisted tree. It implements
// RouteReloader.
type ReloadPublisher struct {
	pub message.Publisher
}

// NewReloadPublisher wraps pub so callers can publish route reloads.
func NewReloadPublisher(pub message.Publisher) *ReloadPublisher {
	return &ReloadPublisher{pub: pub}
}

// ReloadAlertRoutes publishes a reload request for a team. It is fire-and-forget:
// the caller does not wait for the Dispatcher to be rebuilt.
func (p *ReloadPublisher) ReloadAlertRoutes(_ context.Context, orgID uuid.NullUUID, teamID uuid.UUID) error {
	payload, err := json.Marshal(ReloadCommand{TeamID: teamID})
	if err != nil {
		return err
	}
	msg := message.NewMessage(uuid.NewString(), payload)
	setScope(msg, orgID, teamID)
	return p.pub.Publish(TopicRouteReload, msg)
}
