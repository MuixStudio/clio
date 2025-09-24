package dao

import (
	"gorm.io/gorm"
)

type (
	ProfileModel interface {
		profileModel
		GetTableName() string
	}

	customProfileModel struct {
		*defaultProfileModel
	}
)

func (c *customProfileModel) GetTableName() string {
	return c.tableName
}

func NewProfileModel(conn *gorm.DB) ProfileModel {
	return &customProfileModel{
		defaultProfileModel: newProfileModel(conn),
	}
}
