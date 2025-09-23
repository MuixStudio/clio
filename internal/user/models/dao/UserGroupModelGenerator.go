package dao

import (
	"context"

	"gorm.io/gorm"
)

type (
	userGroupModel interface {
		Insert(ctx context.Context, data *UserGroup) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*UserGroup, error)
		FindAll(ctx context.Context, limit int) ([]*UserGroup, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*UserGroup, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *UserGroup) error
	}

	defaultUserGroupModel struct {
		db        *gorm.DB
		tableName string
	}

	UserGroup struct {
		BaseModel
		Name string `gorm:"column:name;not null;type:varchar(30);comment:用户组名称"`
	}
)

func (emp UserGroup) TableName() string {
	return TablePrefix + "_" + "usergroups"
}

func newUserGroupModel(conn *gorm.DB) *defaultUserGroupModel {
	return &defaultUserGroupModel{
		db:        conn,
		tableName: TablePrefix + "_" + "usergroups",
	}
}

func (m *defaultUserGroupModel) Insert(ctx context.Context, data *UserGroup) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultUserGroupModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultUserGroupModel) FindOne(ctx context.Context, id int64) (*UserGroup, error) {
	var result UserGroup
	err := m.db.WithContext(ctx).Model(&UserGroup{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultUserGroupModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*UserGroup, error) {
	var result []*UserGroup
	err := m.db.WithContext(ctx).Model(&UserGroup{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) FindAll(ctx context.Context, limit int) ([]*UserGroup, error) {
	var result []*UserGroup
	err := m.db.WithContext(ctx).Model(&UserGroup{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&UserGroup{}).Count(&result).Error
	return result, err
}

func (m *defaultUserGroupModel) Update(ctx context.Context, data *UserGroup) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultUserGroupModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&UserGroup{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultUserGroupModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*UserGroup, error) {
	var result UserGroup
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultUserGroupModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*UserGroup, error) {
	var result []*UserGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultUserGroupModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*UserGroup, error) {
	var result []*UserGroup
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultUserGroupModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*UserGroup, error) {
	var result []*UserGroup
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
