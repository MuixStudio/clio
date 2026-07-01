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
	alertRepo "github.com/muixstudio/clio/internal/domain/alert/repository"
	channelRepo "github.com/muixstudio/clio/internal/domain/channel/repository"
	connectorRepo "github.com/muixstudio/clio/internal/domain/connector/repository"
	teamRepo "github.com/muixstudio/clio/internal/domain/team/repository"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
	codestrategy "github.com/muixstudio/clio/internal/strategy/code"
)

type Store interface {
	connectorRepo.ConnectorPersister
	alertRepo.AlertPersister
	channelRepo.ChannelPersister
	channelRepo.ChannelMemberPersister
	teamRepo.TeamPersister
	teamRepo.TeamMemberPersister

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
