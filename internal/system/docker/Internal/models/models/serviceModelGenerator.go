package models

import (
	"context"

	"gorm.io/gorm"
)

type (
	serviceModel interface {
		Insert(ctx context.Context, data *Service) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*Service, error)
		FindAll(ctx context.Context, limit int) ([]*Service, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*Service, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *Service) error
	}

	defaultServiceModel struct {
		db        *gorm.DB
		tableName string
	}

	Service struct {
		BaseModel
		Name      string `gorm:"column:name;type:varchar(30);comment:服务名称"`
		ServiceID uint32 `gorm:"column:service_id;type:varchar(64);comment:服务ID(hashID,唯一)"`
	}
)

func (emp Service) TableName() string {
	return TablePrefix + "_" + "tenants"
}

func newServiceModel(conn *gorm.DB) *defaultServiceModel {
	return &defaultServiceModel{
		db:        conn,
		tableName: TablePrefix + "_" + "tenants",
	}
}

func (m *defaultServiceModel) Insert(ctx context.Context, data *Service) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultServiceModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultServiceModel) FindOne(ctx context.Context, id int64) (*Service, error) {
	var result Service
	err := m.db.WithContext(ctx).Model(&Service{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultServiceModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*Service, error) {
	var result []*Service
	err := m.db.WithContext(ctx).Model(&Service{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultServiceModel) FindAll(ctx context.Context, limit int) ([]*Service, error) {
	var result []*Service
	err := m.db.WithContext(ctx).Model(&Service{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultServiceModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&Service{}).Count(&result).Error
	return result, err
}

func (m *defaultServiceModel) Update(ctx context.Context, data *Service) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultServiceModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Service{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultServiceModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*Service, error) {
	var result Service
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultServiceModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*Service, error) {
	var result []*Service
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultServiceModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*Service, error) {
	var result []*Service
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultServiceModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*Service, error) {
	var result []*Service
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
