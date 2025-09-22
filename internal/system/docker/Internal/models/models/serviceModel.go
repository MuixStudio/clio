package models

import (
	"gorm.io/gorm"
)

type (
	ServiceModel interface {
		tenantModel
		GetTableName() string
	}

	customServiceModel struct {
		*defaultServiceModel
	}
)

func (c *customServiceModel) GetTableName() string {
	return c.tableName
}

func NewServiceModel(conn *gorm.DB) ServiceModel {
	return &customServiceModel{
		defaultServiceModel: newServiceModel(conn),
	}
}
