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
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/muixstudio/clio/internal/domain/alert/entity"
)

// transport.go is the message<->domain codec for the dispatch pipeline. It owns
// the topic names, the on-the-wire envelopes, and the encode/decode helpers so
// the bus package can stay a pure transport that carries no domain knowledge.
// Publishers (publisher.go) and consumers (manager.go, handler.go) both speak
// through these helpers.

// Topics the dispatch pipeline publishes to and subscribes from.
const (
	// TopicAlerts carries normalized alert events.
	TopicAlerts = "alerts"
	// TopicRouteReload carries route-reload requests. Route CRUD is performed
	// synchronously against the store by the HTTP handler; only the subsequent
	// reload travels the bus, so the Manager rebuilds the team's Dispatcher from
	// the freshly persisted tree.
	TopicRouteReload = "route_reload"
)

// Message metadata keys. Org/team scope is carried out-of-band so a background
// consumer can rebuild the scoped context route loading needs without trusting
// the (serializable) payload for authorization scope.
const (
	metaOrgID  = "org_id"
	metaTeamID = "team_id"
)

// ReloadCommand is the serializable envelope asking the Manager to rebuild a
// team's Dispatcher from the persisted routing tree. It is the only route
// command on the bus: CRUD is applied directly by the HTTP handler.
type ReloadCommand struct {
	TeamID uuid.UUID `json:"team_id"`
}

// DecodeAlert unmarshals a message payload into a NormalizedAlert.
func DecodeAlert(msg *message.Message) (*entity.NormalizedAlert, error) {
	var a entity.NormalizedAlert
	if err := json.Unmarshal(msg.Payload, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// DecodeReloadCommand unmarshals a message payload into a ReloadCommand.
func DecodeReloadCommand(msg *message.Message) (ReloadCommand, error) {
	var cmd ReloadCommand
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		return ReloadCommand{}, err
	}
	return cmd, nil
}

// ScopeFromMessage extracts the organization and team scope carried in message
// metadata so background consumers can pass it explicitly to scoped operations.
func ScopeFromMessage(msg *message.Message) (orgID uuid.NullUUID, teamID uuid.UUID) {
	if v := msg.Metadata.Get(metaOrgID); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			orgID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}
	if v := msg.Metadata.Get(metaTeamID); v != "" {
		teamID, _ = uuid.Parse(v)
	}
	return orgID, teamID
}

// setScope tags a message with org/team scope. Nil/zero values are omitted so
// the consumer's ScopeFromMessage leaves them unset.
func setScope(msg *message.Message, orgID uuid.NullUUID, teamID uuid.UUID) {
	if orgID.Valid {
		msg.Metadata.Set(metaOrgID, orgID.UUID.String())
	}
	if teamID != uuid.Nil {
		msg.Metadata.Set(metaTeamID, teamID.String())
	}
}
