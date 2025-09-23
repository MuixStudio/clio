package dao

import (
	"gorm.io/gorm"
)

type (
	UserUserGroupModel interface {
		userUserGroupModel
		GetTableName() string
	}

	customUserUserGroupModel struct {
		*defaultUserUserGroupModel
	}
)

func (c *customUserUserGroupModel) GetTableName() string {
	return c.tableName
}

func NewUserUserGroupModel(conn *gorm.DB) UserUserGroupModel {
	return &customUserUserGroupModel{
		defaultUserUserGroupModel: newUserUserGroupModel(conn),
	}
}
