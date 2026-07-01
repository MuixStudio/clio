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

type Connector struct {
	ID uuid.UUID `json:"id"`

	Name  string
	Type  string
	Token string

	OrganizationID uuid.NullUUID
	TeamID         uuid.UUID

	Labels  map[string]string
	Enabled bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Option func(*Connector)

func WithTeamID(id uuid.UUID) Option {
	return func(c *Connector) { c.TeamID = id }
}

func WithLabels(labels map[string]string) Option {
	return func(c *Connector) { c.Labels = labels }
}

func NewConnector(name, tp string, opts ...Option) *Connector {
	c := &Connector{
		Name:    name,
		Type:    tp,
		Token:   uuid.NewString(),
		Enabled: true,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Connector) SetLabels(labels map[string]string) {
	c.Labels = labels
}

func (c *Connector) Disable() {
	c.Enabled = false
}

func (c *Connector) Enable() {
	c.Enabled = true
}

func (c *Connector) RefreshToken() {
	c.Token = uuid.NewString()
}
