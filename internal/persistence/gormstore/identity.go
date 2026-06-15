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
	"encoding/json"

	"github.com/muixstudio/clio/internal/persistence/gormstore/model"

	"github.com/muixstudio/clio/internal/infra/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/muixstudio/clio/internal/identity"
)

func (gs *GormStore) CreateIdentity(ctx context.Context, i *identity.Identity) error {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return err
	}

	m, err := toIdentityModel(i, orgID)
	if err != nil {
		return err
	}
	creds, identifiers := toCredentialModels(i, orgID)

	err = gs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		if len(creds) > 0 {
			if err := tx.Create(&creds).Error; err != nil {
				return err
			}
		}
		if len(identifiers) > 0 {
			if err := tx.Create(&identifiers).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.Conflict("RECORD_ALREADY_EXISTS", "identity already exist").WithCause(err)
	}
	return err
}

func (gs *GormStore) FindByCredentialsIdentifier(ctx context.Context, ct identity.CredentialsType, identifier string) (*identity.Identity, *identity.Credentials, error) {
	orgID, err := gs.orgIDFromCtx(ctx)
	if err != nil {
		return nil, nil, err
	}

	var ci model.CredentialIdentifier
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"type":            string(ct),
			"identifier":      identifier,
		}).
		First(&ci).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errors.NotFound(
			"RECORD_NOT_FOUND",
			"identifier not found",
		).WithCause(err).WithMetadata(map[string]string{
			"identifier":       identifier,
			"credentials_type": string(ct),
		})
	}
	if err != nil {
		return nil, nil, err
	}

	var cm model.Credential
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"id":              ci.CredentialID,
		}).
		First(&cm).Error
	if err != nil {
		return nil, nil, err
	}

	var m model.Identity
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"id":              cm.IdentityID,
		}).
		First(&m).Error
	if err != nil {
		return nil, nil, err
	}

	var identifiers []model.CredentialIdentifier
	err = gs.db.WithContext(ctx).
		Where(map[string]any{
			"organization_id": orgID,
			"credential_id":   cm.ID,
		}).
		Find(&identifiers).Error
	if err != nil {
		return nil, nil, err
	}

	i, err := fromIdentityModel(&m)
	if err != nil {
		return nil, nil, err
	}

	creds := toCredentials(ct, &cm, identifiers)
	i.Credentials[ct] = creds
	return i, creds, nil
}

func toIdentityModel(i *identity.Identity, orgID uuid.NullUUID) (*model.Identity, error) {
	traits, err := json.Marshal(i.Traits)
	if err != nil {
		return nil, err
	}

	return &model.Identity{
		ID:             i.ID,
		OrganizationID: orgID,
		Traits:         string(traits),
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}, nil
}

// toCredentialModels flattens an identity's credentials into rows for the
// identity_credentials and identity_credential_identifiers tables.
func toCredentialModels(i *identity.Identity, orgID uuid.NullUUID) ([]model.Credential, []model.CredentialIdentifier) {
	var creds []model.Credential
	var identifiers []model.CredentialIdentifier

	for ct, c := range i.Credentials {
		credID := uuid.New()
		creds = append(creds, model.Credential{
			ID:             credID,
			OrganizationID: orgID,
			IdentityID:     i.ID,
			Type:           string(ct),
			Config:         string(c.Config),
		})
		for _, id := range c.Identifiers {
			identifiers = append(identifiers, model.CredentialIdentifier{
				ID:             uuid.New(),
				OrganizationID: orgID,
				CredentialID:   credID,
				Type:           string(ct),
				Identifier:     id,
			})
		}
	}

	return creds, identifiers
}

// toCredentials assembles the domain Credentials value for a single
// credential row and its identifiers.
func toCredentials(ct identity.CredentialsType, cm *model.Credential, identifiers []model.CredentialIdentifier) *identity.Credentials {
	ids := make([]string, 0, len(identifiers))
	for _, idm := range identifiers {
		ids = append(ids, idm.Identifier)
	}
	return &identity.Credentials{
		Type:        ct,
		Identifiers: ids,
		Config:      []byte(cm.Config),
	}
}

func fromIdentityModel(m *model.Identity) (*identity.Identity, error) {
	var traits map[string]any
	if err := json.Unmarshal([]byte(m.Traits), &traits); err != nil {
		return nil, err
	}

	return &identity.Identity{
		ID:          m.ID,
		Traits:      traits,
		Credentials: make(map[identity.CredentialsType]*identity.Credentials),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}
