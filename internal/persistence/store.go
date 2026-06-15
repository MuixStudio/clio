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
	"github.com/muixstudio/clio/internal/alert"
	"github.com/muixstudio/clio/internal/channel"
	"github.com/muixstudio/clio/internal/connector"
	loginFlow "github.com/muixstudio/clio/internal/flow/login"
	logoutLfow "github.com/muixstudio/clio/internal/flow/logout"
	recoveryFlow "github.com/muixstudio/clio/internal/flow/recovery"
	registrationFlow "github.com/muixstudio/clio/internal/flow/registration"
	verificationFlow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/session"
	codestrategy "github.com/muixstudio/clio/internal/strategy/code"
	"github.com/muixstudio/clio/internal/team"
)

type Store interface {
	connector.ConnectorPersister
	alert.AlertPersister
	channel.ChannelPersister
	channel.ChannelMemberPersister
	team.TeamPersister
	team.TeamMemberPersister

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
