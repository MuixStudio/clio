package dao

import (
	"gorm.io/gorm"
)

type (
	UserGroupModel interface {
		userGroupModel
		GetTableName() string
	}

	customUserGroupModel struct {
		*defaultUserGroupModel
	}
)

func (c *customUserGroupModel) GetTableName() string {
	return c.tableName
}

func NewUserGroupModel(conn *gorm.DB) UserGroupModel {
	return &customUserGroupModel{
		defaultUserGroupModel: newUserGroupModel(conn),
	}
}
