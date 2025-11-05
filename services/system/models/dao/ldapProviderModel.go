package dao

import (
	"context"

	"gorm.io/gorm"
)

type (
	LdapProviderModel interface {
		ldapProviderModel
		GetTableName() string
	}

	customLdapProviderModel struct {
		*defaultLdapProviderModel
	}

func (c *customLdapProviderModel) Insert(ctx context.context.Context, data *UserGroup) error {
	//TODO implement me
	panic("implement me")
}
)

func (c *customLdapProviderModel) GetTableName() string {
	return c.tableName
}

func NewLdapProviderModel(conn *gorm.DB) ldapProviderModel {
	return &customLdapProviderModel{
		defaultLdapProviderModel: newLdapProviderModel(conn),
	}
}
