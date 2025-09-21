package models

import (
	"gorm.io/gorm"
)

type (
	UserUserGroupModel interface {
		userGroupModel
		GetTableName() string
	}

	customUserUserGroupModel struct {
		*defaultUserGroupModel
	}
)

func (c *customUserUserGroupModel) GetTableName() string {
	return c.tableName
}

func NewUserUserGroupModel(conn *gorm.DB) UserUserGroupModel {
	return &customUserUserGroupModel{
		defaultUserGroupModel: newUserGroupModel(conn),
	}
}
