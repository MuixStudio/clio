package dao

import (
	"gorm.io/gorm"
)

type (
	TokenModel interface {
		tokenModel
		GetTableName() string
	}

	customTokenModel struct {
		*defaultTokenModel
	}
)

func (c *customTokenModel) GetTableName() string {
	return c.tableName
}

func NewTokenModel(db *gorm.DB) TokenModel {
	return &customTokenModel{
		defaultTokenModel: newTokenModel(db),
	}
}
