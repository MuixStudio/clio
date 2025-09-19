package models

import (
	"context"

	"gorm.io/gorm"
)

type (
	userGroupModel interface {
		Insert(ctx context.Context, data *userGroup) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*userGroup, error)
		FindAll(ctx context.Context, limit int) ([]*userGroup, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*userGroup, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *userGroup) error
	}

	defaultUserGroupModel struct {
		db        *gorm.DB
		tableName string
	}

	userGroup struct {
		BaseModel
		Name       uint32 `gorm:"column:name;type:varchar(30);comment:用户组名称"`
		TenantName uint32 `gorm:"column:tenant_name;type:varchar(30);not null;unique;comment:租户名(唯一)"`
	}
)

func (emp userGroup) TableName() string {
	return TablePrefix + "_" + "user_groups"
}

func newUserGroupModel(conn *gorm.DB) *defaultUserGroupModel {
	return &defaultUserGroupModel{
		db:        conn,
		tableName: TablePrefix + "_" + "user_groups",
	}
}

func (m *defaultUserGroupModel) Insert(ctx context.Context, data *userGroup) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultUserGroupModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultUserGroupModel) FindOne(ctx context.Context, id int64) (*userGroup, error) {
	var result userGroup
	err := m.db.WithContext(ctx).Model(&userGroup{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultUserGroupModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*userGroup, error) {
	var result []*userGroup
	err := m.db.WithContext(ctx).Model(&userGroup{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) FindAll(ctx context.Context, limit int) ([]*userGroup, error) {
	var result []*userGroup
	err := m.db.WithContext(ctx).Model(&userGroup{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&userGroup{}).Count(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) Update(ctx context.Context, data *userGroup) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultUserGroupModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&userGroup{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultUserGroupModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*userGroup, error) {
	var result userGroup
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultUserGroupModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*userGroup, error) {
	var result []*userGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultUserGroupModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*userGroup, error) {
	var result []*userGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultUserGroupModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*userGroup, error) {
	var result []*userGroup
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
