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

type VisibilityType string

const (
	Public  VisibilityType = "public"
	Private VisibilityType = "private"
)

type Option func(*Channel)

func WithDescription(desc string) Option {
	return func(c *Channel) {
		c.Description = desc
	}
}

type Channel struct {
	ID          uuid.UUID
	Description string
	TeamID      uuid.UUID
	Name        string
	Visibility  VisibilityType
	Enabled     bool
	IsDefault   bool
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func new(teamID uuid.UUID, owner uuid.UUID, name string, visibilityType VisibilityType, opts ...Option) *Channel {
	c := &Channel{
		ID:         uuid.New(),
		TeamID:     teamID,
		Name:       name,
		Visibility: visibilityType,
		Enabled:    true,
		IsDefault:  false,
		CreatedBy:  owner,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

func NewDefaultChannel(teamID uuid.UUID, owner uuid.UUID, opts ...Option) *Channel {
	c := &Channel{
		ID:         uuid.New(),
		TeamID:     teamID,
		Name:       "default",
		Visibility: Public,
		Enabled:    true,
		IsDefault:  true,
		CreatedBy:  owner,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func NewPublicChannel(teamID uuid.UUID, owner uuid.UUID, name string, opts ...Option) *Channel {
	return new(teamID, owner, name, Public, opts...)
}

func NewPrivateChannel(teamID uuid.UUID, owner uuid.UUID, name string, opts ...Option) *Channel {
	return new(teamID, owner, name, Private, opts...)
}
