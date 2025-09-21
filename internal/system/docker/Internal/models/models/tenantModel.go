package models

import (
	"gorm.io/gorm"
)

type (
	TenantModel interface {
		tenantModel
		GetTableName() string
	}

	customTenantModel struct {
		*defaultTenantModel
	}
)

func (c *customTenantModel) GetTableName() string {
	return c.tableName
}

func NewTenantModel(conn *gorm.DB) TenantModel {
	return &customTenantModel{
		defaultTenantModel: newTenantModel(conn),
	}
}
