package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	sessionModel interface {
		Insert(ctx context.Context, data *Session) error
		FindOne(ctx context.Context, id int64) (*Session, error)
		Update(ctx context.Context, data *Session) error
	}

	defaultSessionModel struct {
		db        *gorm.DB
		tableName string
	}

	Session struct {
		BaseModel
		UseID       uint32    `gorm:"column:user_id;type:int;not null;comment:用户ID"`
		IP          string    `gorm:"column:ip;type:varchar(30);not null;comment:ip"`
		Platform    time.Time `gorm:"column:platform;comment:平台设备信息"`
		LastLoginAt time.Time `gorm:"column:last_login_at;comment:最后登陆时间"`
	}
)

func (emp Session) TableName() string {
	return TablePrefix + "_" + "sessions"
}

func newSessionModel(conn *gorm.DB) *defaultSessionModel {
	return &defaultSessionModel{
		db:        conn,
		tableName: TablePrefix + "_" + "sessions",
	}
}

func (m *defaultSessionModel) Insert(ctx context.Context, data *Session) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultSessionModel) FindOne(ctx context.Context, id int64) (*Session, error) {
	var result Session
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultSessionModel) Update(ctx context.Context, data *Session) error {
	return m.db.WithContext(ctx).Save(data).Error
}
