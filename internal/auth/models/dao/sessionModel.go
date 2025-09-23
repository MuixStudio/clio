package dao

import (
	"gorm.io/gorm"
)

type (
	SessionModel interface {
		sessionModel
		GetTableName() string
	}

	customSessionModel struct {
		*defaultSessionModel
	}
)

func (c *customSessionModel) GetTableName() string {
	return c.tableName
}

func NewSessionModel(db *gorm.DB) SessionModel {
	return &customSessionModel{
		defaultSessionModel: newSessionModel(db),
	}
}
