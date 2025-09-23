package dao

import (
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
)

func (c *customLdapProviderModel) GetTableName() string {
	return c.tableName
}

func NewLdapProviderModel(conn *gorm.DB) ldapProviderModel {
	return &customLdapProviderModel{
		defaultLdapProviderModel: newLdapProviderModel(conn),
	}
}
