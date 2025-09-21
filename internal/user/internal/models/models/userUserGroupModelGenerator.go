package models

import (
	"context"

	"gorm.io/gorm"
)

type (
	userUserGroupModel interface {
		Insert(ctx context.Context, data *userUserGroup) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*userUserGroup, error)
		FindAll(ctx context.Context, limit int) ([]*userUserGroup, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*userUserGroup, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *userUserGroup) error
	}

	defaultUserUserGroupModel struct {
		db        *gorm.DB
		tableName string
	}

	userUserGroup struct {
		BaseModel
		UseID      uint32 `gorm:"type:int;comment:用户ID"`
		UseGroupID uint32 `gorm:"type:int;comment:用户组ID"`
	}
)

func (emp userUserGroup) TableName() string {
	return TablePrefix + "_" + "user_usergroups"
}

func newUserUserGroupModel(conn *gorm.DB) *defaultUserUserGroupModel {
	return &defaultUserUserGroupModel{
		db:        conn,
		tableName: TablePrefix + "_" + "user_usergroups",
	}
}

func (m *defaultUserUserGroupModel) Insert(ctx context.Context, data *userUserGroup) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultUserUserGroupModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultUserUserGroupModel) FindOne(ctx context.Context, id int64) (*userUserGroup, error) {
	var result userUserGroup
	err := m.db.WithContext(ctx).Model(&userUserGroup{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultUserUserGroupModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*userUserGroup, error) {
	var result []*userUserGroup
	err := m.db.WithContext(ctx).Model(&userUserGroup{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserUserGroupModel) FindAll(ctx context.Context, limit int) ([]*userUserGroup, error) {
	var result []*userUserGroup
	err := m.db.WithContext(ctx).Model(&userUserGroup{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserUserGroupModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&userUserGroup{}).Count(&result).Error
	return result, err
}

func (m *defaultUserUserGroupModel) Update(ctx context.Context, data *userUserGroup) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultUserUserGroupModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&userUserGroup{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultUserUserGroupModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*userUserGroup, error) {
	var result userUserGroup
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultUserUserGroupModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*userUserGroup, error) {
	var result []*userUserGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultUserUserGroupModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*userUserGroup, error) {
	var result []*userUserGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultUserUserGroupModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*userUserGroup, error) {
	var result []*userUserGroup
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
