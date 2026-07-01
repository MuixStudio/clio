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

package persistence

import (
	"context"

	"github.com/google/uuid"
	repository2 "github.com/muixstudio/clio/internal/domain/alert/repository"
	"github.com/muixstudio/clio/internal/domain/channel/repository"
	repository3 "github.com/muixstudio/clio/internal/domain/connector/repository"
	repository4 "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/infra/zus/dispatch"

	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
	codestrategy "github.com/muixstudio/clio/internal/strategy/code"
)

// AlertRoutePersister manages alert route rows and loads a team's routing tree
// for the dispatcher.
type AlertRoutePersister interface {
	LoadAlertRouteTree(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) (*dispatch.Route, error)
	AddAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in dispatch.AddRouteInput) (uuid.UUID, error)
	UpdateAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, in dispatch.UpdateRouteInput) error
	DeleteAlertRoute(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, routeID uuid.UUID) error
}

type Store interface {
	repository3.ConnectorPersister
	repository2.AlertPersister
	AlertRoutePersister
	repository.ChannelPersister
	repository.ChannelMemberPersister
	repository4.TeamPersister
	repository4.TeamMemberPersister

	session.SessionPersister
	session.SessionTokenExchangeCodePersister

	identity.IdentityPersister
	identity.IdentityVerifiableAddressPersister

	loginFlow.FlowPersister
	logoutLfow.FlowPersister
	recoveryFlow.FlowPersister
	registrationFlow.FlowPersister
	verificationFlow.FlowPersister

	codestrategy.VerificationCodePersister
}
