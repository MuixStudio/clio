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

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/domain/team/entity"
)

type TeamUpdate struct {
	Name        *string
	Description *string
}

type TeamRoleUpdate struct {
	Name        *string
	Description *string
}

type (
	TeamPersister interface {
		CreateTeam(ctx context.Context, t *entity.Team) error
		GetTeam(ctx context.Context, orgId uuid.NullUUID, id uuid.UUID) (*entity.Team, error)
		GetTeams(ctx context.Context, orgId uuid.NullUUID, ids *[]uuid.UUID) (*[]entity.Team, error)
		ListTeamsByUserID(ctx context.Context, orgId uuid.NullUUID, userID uuid.UUID) ([]*entity.Team, error)
		UpdateTeam(ctx context.Context, orgId uuid.NullUUID, id uuid.UUID, update TeamUpdate) error
		DeleteTeam(ctx context.Context, orgId uuid.NullUUID, id uuid.UUID) error
		ListTeams(ctx context.Context, orgId uuid.NullUUID) ([]*entity.Team, error)
	}
	TeamPersisterProvider interface {
		TeamPersister() TeamPersister
	}

	TeamMemberPersister interface {
		AddMember(ctx context.Context, orgId uuid.NullUUID, teamID uuid.UUID, userID uuid.UUID) error
		GetMember(ctx context.Context, orgId uuid.NullUUID, id uuid.UUID) (*entity.TeamMember, error)
		RemoveMember(ctx context.Context, orgId uuid.NullUUID, id uuid.UUID) error
		ListMembers(ctx context.Context, orgId uuid.NullUUID, teamID uuid.UUID) ([]*entity.TeamMember, error)
		//UpdateMemberRole(ctx context.Context, memberID uuid.UUID, roleID uuid.UUID) error
	}
	TeamMemberPersisterProvider interface {
		TeamMemberPersister() TeamMemberPersister
	}

	//TeamRolePersister interface {
	//	CreateRole(ctx context.Context, r *TeamRole) error
	//	GetRole(ctx context.Context, id uuid.UUID) (*TeamRole, error)
	//	UpdateRole(ctx context.Context, id uuid.UUID, update TeamRoleUpdate) error
	//	DeleteRole(ctx context.Context, id uuid.UUID) error
	//	ListRoles(ctx context.Context, teamID uuid.UUID) ([]*TeamRole, error)
	//}
	//TeamRolePersisterProvider interface {
	//	TeamRolePersister() TeamRolePersister
	//}
)
