package models

import (
	"gorm.io/gorm"
)

type (
	UserModel interface {
		userModel
		GetTableName() string
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func (c *customUserModel) GetTableName() string {
	return c.tableName
}

func NewUserModel(conn *gorm.DB) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}
