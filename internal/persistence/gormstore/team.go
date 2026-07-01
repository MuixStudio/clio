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

package gormstore

import (
	"context"
	"time"

	"github.com/muixstudio/clio/internal/domain/team/entity"
	team2 "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/persistence/gormstore/model"
	"gorm.io/datatypes"

	"github.com/google/uuid"
	"github.com/muixstudio/clio/internal/infra/errors"
	"gorm.io/gorm"
)

// ── Team ──────────────────────────────────────────────────────────────────────

func toTeamModel(t *entity.Team) *model.Team {
	return &model.Team{
		ID:             t.ID,
		Name:           t.Name,
		Description:    t.Description,
		OrganizationID: t.OrganizationID,
	}
}

func fromTeamModel(m *model.Team) *entity.Team {
	return &entity.Team{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		OrganizationID: m.OrganizationID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (gs *GormStore) CreateTeam(ctx context.Context, t *entity.Team) error {
	if t.ID == (uuid.UUID{}) {
		t.ID = uuid.New()
	}
	m := toTeamModel(t)

	// Create the team together with its root alert route in one transaction:
	// every team always owns exactly one root route, which can never be deleted.
	err := gs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return tx.Create(rootAlertRoute(t.OrganizationID, t.ID)).Error
	})
	if err != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(err)
	}
	return nil
}

func (gs *GormStore) GetTeam(ctx context.Context, orgID uuid.NullUUID, id uuid.UUID) (*entity.Team, error) {
	var m model.Team
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("RECORD_NOT_FOUND", "team not found: "+id.String())
		}
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return fromTeamModel(&m), nil
}

func (gs *GormStore) GetTeams(ctx context.Context, orgID uuid.NullUUID, ids *[]uuid.UUID) (*[]entity.Team, error) {
	var models []model.Team
	result := gs.db.WithContext(ctx).
		Where("organization_id = ? AND id IN ?", orgID, *ids).
		Find(&models)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	teams := make([]entity.Team, 0, len(models))
	for i := range models {
		teams = append(teams, *fromTeamModel(&models[i]))
	}
	return &teams, nil
}

func (gs *GormStore) ListTeamsByUserID(ctx context.Context, orgID uuid.NullUUID, userID uuid.UUID) ([]*entity.Team, error) {
	var models []model.Team
	result := gs.db.WithContext(ctx).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("teams.organization_id = ? AND team_members.user_id = ?", orgID, userID).
		Find(&models)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	teams := make([]*entity.Team, 0, len(models))
	for i := range models {
		teams = append(teams, fromTeamModel(&models[i]))
	}
	return teams, nil
}

func (gs *GormStore) UpdateTeam(ctx context.Context, orgID uuid.NullUUID, id uuid.UUID, update team2.TeamUpdate) error {
	updates := map[string]any{}
	if update.Name != nil {
		updates["name"] = *update.Name
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}

	result := gs.db.WithContext(ctx).Model(&model.Team{}).
		Where(map[string]any{
			"id":              id,
			"organization_id": orgID,
		}).
		Updates(updates)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "team not found: "+id.String())
	}
	return nil
}

func (gs *GormStore) DeleteTeam(ctx context.Context, orgID uuid.NullUUID, id uuid.UUID) error {
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).Delete(&model.Team{})
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "team not found: "+id.String())
	}
	return nil
}

func (gs *GormStore) ListTeams(ctx context.Context, orgID uuid.NullUUID) ([]*entity.Team, error) {
	var models []model.Team
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"organization_id": orgID,
	}).Find(&models)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	teams := make([]*entity.Team, 0, len(models))
	for i := range models {
		teams = append(teams, fromTeamModel(&models[i]))
	}
	return teams, nil
}

// ── TeamRole ──────────────────────────────────────────────────────────────────

//func toTeamRoleModel(r *team.TeamRole) *teamRoleModel {
//	return &teamRoleModel{
//		ID:           r.ID,
//		TeamID:       r.TeamID,
//		Name:         r.Name,
//		Description:  r.Description,
//		IsSystemRole: r.IsSystemRole,
//	}
//}
//
//func fromTeamRoleModel(m *teamRoleModel) (*team.TeamRole, error) {
//	return &team.TeamRole{
//		ID:           m.ID,
//		TeamID:       m.TeamID,
//		Name:         m.Name,
//		Description:  m.Description,
//		IsSystemRole: m.IsSystemRole,
//		CreatedAt:    m.CreatedAt,
//		UpdatedAt:    m.UpdatedAt,
//	}, nil
//}
//
//func (gs *GormStore) CreateRole(ctx context.Context, r *team.TeamRole) error {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return err
//	}
//
//	if r.ID == (uuid.UUID{}) {
//		r.ID = uuid.New()
//	}
//	m := toTeamRoleModel(r)
//	m.OrganizationID = orgID
//	result := gs.db.WithContext(ctx).Create(m)
//	if result.Error != nil {
//		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	return nil
//}
//
//func (gs *GormStore) GetRole(ctx context.Context, id uuid.UUID) (*team.TeamRole, error) {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	var m teamRoleModel
//	result := gs.db.WithContext(ctx).Where(map[string]any{
//		"id":              id,
//		"organization_id": orgID,
//	}).First(&m)
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return nil, errors.NotFound("RECORD_NOT_FOUND", "team role not found: "+id.String())
//		}
//		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	return fromTeamRoleModel(&m)
//}
//
//func (gs *GormStore) UpdateRole(ctx context.Context, id uuid.UUID, update team.TeamRoleUpdate) error {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return err
//	}
//
//	updates := map[string]any{}
//	if update.Name != nil {
//		updates["name"] = *update.Name
//	}
//	if update.Description != nil {
//		updates["description"] = *update.Description
//	}
//
//	result := gs.db.WithContext(ctx).Model(&teamRoleModel{}).
//		Where(map[string]any{
//			"id":              id,
//			"organization_id": orgID,
//		}).
//		Updates(updates)
//	if result.Error != nil {
//		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	if result.RowsAffected == 0 {
//		return errors.NotFound("RECORD_NOT_FOUND", "team role not found: "+id.String())
//	}
//	return nil
//}
//
//func (gs *GormStore) DeleteRole(ctx context.Context, id uuid.UUID) error {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return err
//	}
//
//	result := gs.db.WithContext(ctx).Where(map[string]any{
//		"id":              id,
//		"organization_id": orgID,
//	}).Delete(&teamRoleModel{})
//	if result.Error != nil {
//		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	if result.RowsAffected == 0 {
//		return errors.NotFound("RECORD_NOT_FOUND", "team role not found: "+id.String())
//	}
//	return nil
//}
//
//func (gs *GormStore) ListRoles(ctx context.Context, teamID uuid.UUID) ([]*team.TeamRole, error) {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	var models []teamRoleModel
//	result := gs.db.WithContext(ctx).
//		Where(map[string]any{
//			"team_id":         teamID,
//			"organization_id": orgID,
//		}).
//		Find(&models)
//	if result.Error != nil {
//		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	roles := make([]*team.TeamRole, 0, len(models))
//	for i := range models {
//		r, err := fromTeamRoleModel(&models[i])
//		if err != nil {
//			return nil, err
//		}
//		roles = append(roles, r)
//	}
//	return roles, nil
//}

// ── TeamMember ────────────────────────────────────────────────────────────────

func fromTeamMemberModel(m *model.TeamMember) (*entity.TeamMember, error) {
	return &entity.TeamMember{
		ID:     m.ID,
		TeamID: m.TeamID,
		UserID: m.UserID,
		Email:  m.Email,
		//RoleID:    m.RoleID,
		JoinedAt:  m.JoinedAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (gs *GormStore) AddMember(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID, userID uuid.UUID) error {
	model := &model.TeamMember{
		ID:       uuid.New(),
		TeamID:   teamID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
	model.OrganizationID = orgID
	result := gs.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return nil
}

func (gs *GormStore) GetMember(ctx context.Context, orgID uuid.NullUUID, id uuid.UUID) (*entity.TeamMember, error) {
	var m model.TeamMember
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).First(&m)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("RECORD_NOT_FOUND", "team member not found: "+id.String())
		}
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	return fromTeamMemberModel(&m)
}

func (gs *GormStore) RemoveMember(ctx context.Context, orgID uuid.NullUUID, id uuid.UUID) error {
	result := gs.db.WithContext(ctx).Where(map[string]any{
		"id":              id,
		"organization_id": orgID,
	}).Delete(&model.TeamMember{})
	if result.Error != nil {
		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("RECORD_NOT_FOUND", "team member not found: "+id.String())
	}
	return nil
}

func (gs *GormStore) ListMembers(ctx context.Context, orgID uuid.NullUUID, teamID uuid.UUID) ([]*entity.TeamMember, error) {
	var models []model.TeamMember
	result := gs.db.WithContext(ctx).
		Where(map[string]any{
			"team_id":         teamID,
			"organization_id": orgID,
		}).
		Find(&models)
	if result.Error != nil {
		return nil, errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
	}
	members := make([]*entity.TeamMember, 0, len(models))
	for i := range models {
		member, err := fromTeamMemberModel(&models[i])
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

//func (gs *GormStore) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, roleID uuid.UUID) error {
//	orgID, err := gs.orgIDFromCtx(ctx)
//	if err != nil {
//		return err
//	}
//
//	result := gs.db.WithContext(ctx).Model(&model.TeamMember{}).
//		Where(map[string]any{
//			"id":              memberID,
//			"organization_id": orgID,
//		}).
//		Update("role_id", roleID)
//	if result.Error != nil {
//		return errors.InternalServerError("INTERNAL_SERVER_ERROR", "persistence: infra has unknow error").WithCause(result.Error)
//	}
//	if result.RowsAffected == 0 {
//		return errors.NotFound("RECORD_NOT_FOUND", "team member not found: "+memberID.String())
//	}
//	return nil
//}

func rootAlertRoute(orgID uuid.NullUUID, teamID uuid.UUID) *model.AlertRoute {
	return &model.AlertRoute{
		ID:             uuid.New(),
		OrganizationID: orgID,
		TeamID:         teamID,
		Name:           "root",
		ParentID:       uuid.NullUUID{}, // NULL => root
		Priority:       0,
		Matchers:       datatypes.JSON("[]"),
		Enabled:        true,
		Continue:       false,
	}
}
