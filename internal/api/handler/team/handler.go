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

package team

import (
	"github.com/muixstudio/clio/internal/alert"
	"github.com/muixstudio/clio/internal/connector"
	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/logger"
	"github.com/muixstudio/clio/internal/team"
)

type dependencies interface {
	connector.ConnectorPersisterProvider
	alert.AlertPersisterProvider
	connector.ConnectorProviderProvider
	team.TeamPersisterProvider
	team.TeamMemberPersisterProvider

	logger.Logger
	config.Provider
}

// TeamHandler handles connector management endpoints.
type TeamHandler struct {
	d dependencies
}

// NewTeamHandler returns a TeamHandler.
func NewTeamHandler(d dependencies) *TeamHandler {
	return &TeamHandler{d: d}
}
