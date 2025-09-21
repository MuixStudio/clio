package models

import (
	"context"

	"gorm.io/gorm"
)

type (
	userProfileModel interface {
		Insert(ctx context.Context, data *UserProfile) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*UserProfile, error)
		FindAll(ctx context.Context, limit int) ([]*UserProfile, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*UserProfile, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *UserProfile) error
		FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*UserProfile, error)
	}

	defaultUserProfileModel struct {
		db        *gorm.DB
		tableName string
	}

	UserProfile struct {
		BaseModel
		UseID    uint32 `gorm:"type:int;comment:用户ID"`
		Avatar   string `gorm:"column:avatar;type:varchar(255);comment:头像"`
		Locale   string `gorm:"column:locale;type:varchar(255);comment:地址"`
		Timezone string `gorm:"column:timezone;type:varchar(30);comment:时区"`
	}
)

func (emp UserProfile) TableName() string {
	return TablePrefix + "_" + "user_profiles"
}

func newUserProfileModel(conn *gorm.DB) *defaultUserProfileModel {
	return &defaultUserProfileModel{
		db:        conn,
		tableName: TablePrefix + "_" + "user_profiles",
	}
}

func (m *defaultUserProfileModel) Insert(ctx context.Context, data *UserProfile) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultUserProfileModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultUserProfileModel) FindOne(ctx context.Context, id int64) (*UserProfile, error) {
	var result UserProfile
	err := m.db.WithContext(ctx).Model(&UserProfile{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultUserProfileModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*UserProfile, error) {
	var result []*UserProfile
	err := m.db.WithContext(ctx).Model(&UserProfile{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserProfileModel) FindAll(ctx context.Context, limit int) ([]*UserProfile, error) {
	var result []*UserProfile
	err := m.db.WithContext(ctx).Model(&UserProfile{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserProfileModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&UserProfile{}).Count(&result).Error
	return result, err
}

func (m *defaultUserProfileModel) Update(ctx context.Context, data *UserProfile) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultUserProfileModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&UserProfile{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultUserProfileModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*UserProfile, error) {
	var result UserProfile
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultUserProfileModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*UserProfile, error) {
	var result []*UserProfile
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultUserProfileModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*UserProfile, error) {
	var result []*UserProfile
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultUserProfileModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*UserProfile, error) {
	var result []*UserProfile
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
