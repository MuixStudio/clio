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

package connector

import (
	"context"

	"github.com/google/uuid"
)

type Update struct {
	Name    *string
	Enabled *bool
	TeamID  *uuid.UUID
}

type ListOptions struct {
	Page, PageSize int
	Enable         *bool
	Type           *string
	TeamID         *uuid.UUID
}

type CountOptions struct {
	Enable *bool
	Type   *string
	TeamID *uuid.UUID
}

type (
	// ConnectorPersister manages connector factories and live connector instances.
	ConnectorPersister interface {
		// Create persists a new ConnectorRecord and returns the live Connector.
		CreateConnector(ctx context.Context, c *Connector) error
		// Get retrieves a live Connector by ID, loading from DB on cache miss.
		GetConnector(ctx context.Context, id uuid.UUID) (*Connector, error)
		// List returns all persisted ConnectorRecords.
		ListConnectors(ctx context.Context, options *ListOptions) ([]*Connector, error)
		// UpdateConnector updates a connector's mutable fields.
		UpdateConnector(ctx context.Context, id uuid.UUID, updates *Update) error
		DeleteConnector(ctx context.Context, id uuid.UUID, teamID uuid.UUID) error

		Count(ctx context.Context, options CountOptions) (int64, error)
	}
	ConnectorPersisterProvider interface {
		ConnectorPersister() ConnectorPersister
	}
)
