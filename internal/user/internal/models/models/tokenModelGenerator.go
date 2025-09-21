package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	tokenModel interface {
		Insert(ctx context.Context, data *Token) error
		FindOne(ctx context.Context, id int64) (*Token, error)
		Update(ctx context.Context, data *Token) error
	}

	defaultTokenModel struct {
		db        *gorm.DB
		tableName string
	}

	Token struct {
		BaseModel
		Name           string    `gorm:"column:name;type:varchar(30);not null;comment:token名称"`
		ExpirationTime time.Time `gorm:"column:expiration_time;comment:过期时间(为空时永不过期)"`
		Scope          time.Time `gorm:"column:scope;unique;comment:权限范围"`
	}
)

func (emp Token) TableName() string {
	return TablePrefix + "_" + "tokens"
}

func newTokenModel(conn *gorm.DB) *defaultTokenModel {
	return &defaultTokenModel{
		db:        conn,
		tableName: TablePrefix + "_" + "tokens",
	}
}

func (m *defaultTokenModel) Insert(ctx context.Context, data *Token) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultTokenModel) FindOne(ctx context.Context, id int64) (*Token, error) {
	var result Token
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultTokenModel) Update(ctx context.Context, data *Token) error {
	return m.db.WithContext(ctx).Save(data).Error
}
