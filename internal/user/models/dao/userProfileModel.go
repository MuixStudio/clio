package dao

import (
	"gorm.io/gorm"
)

type (
	UserProfileModel interface {
		userProfileModel
		GetTableName() string
	}

	customUserProfileModel struct {
		*defaultUserProfileModel
	}
)

func (c *customUserProfileModel) GetTableName() string {
	return c.tableName
}

func NewUserProfileModel(conn *gorm.DB) UserProfileModel {
	return &customUserProfileModel{
		defaultUserProfileModel: newUserProfileModel(conn),
	}
}
