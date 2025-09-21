package models

import (
	"context"

	"gorm.io/gorm"
)

const (
	WINDOWSMICROSOFTAD = iota + 1
	OPENLDAP
)

type (
	ldapProviderModel interface {
		Insert(ctx context.Context, data *LdapProvider) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*LdapProvider, error)
		FindAll(ctx context.Context, limit int) ([]*LdapProvider, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*LdapProvider, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *LdapProvider) error
		FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*LdapProvider, error)
	}

	defaultLdapProviderModel struct {
		db        *gorm.DB
		tableName string
	}

	LdapProviderType uint16

	LdapProvider struct {
		BaseModel
		Name             string           `gorm:"column:name;type:varchar(30);not null;comment:Ldap提供商名称"`
		LdapProviderType LdapProviderType `gorm:"column:ldap_provider_type;type:smallint;not null;comment:Ldap提供商类型"`
	}
)

var LdapProviderTypeMap = map[uint16]string{
	WINDOWSMICROSOFTAD: "WindowsMicrosoftAD",
	OPENLDAP:           "OpenLdap",
}

func (ap LdapProviderType) Desc() string {
	return LdapProviderTypeMap[uint16(ap)]
}

func (emp LdapProvider) TableName() string {
	return TablePrefix + "_" + "ldap_providers"
}

func newLdapProviderModel(conn *gorm.DB) *defaultLdapProviderModel {
	return &defaultLdapProviderModel{
		db:        conn,
		tableName: TablePrefix + "_" + "ldap_providers",
	}
}

func (m *defaultLdapProviderModel) Insert(ctx context.Context, data *LdapProvider) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultLdapProviderModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultLdapProviderModel) FindOne(ctx context.Context, id int64) (*LdapProvider, error) {
	var result LdapProvider
	err := m.db.WithContext(ctx).Model(&LdapProvider{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultLdapProviderModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*LdapProvider, error) {
	var result []*LdapProvider
	err := m.db.WithContext(ctx).Model(&LdapProvider{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultLdapProviderModel) FindAll(ctx context.Context, limit int) ([]*LdapProvider, error) {
	var result []*LdapProvider
	err := m.db.WithContext(ctx).Model(&LdapProvider{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultLdapProviderModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&LdapProvider{}).Count(&result).Error
	return result, err
}

func (m *defaultLdapProviderModel) Update(ctx context.Context, data *LdapProvider) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultLdapProviderModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&LdapProvider{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultLdapProviderModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*LdapProvider, error) {
	var result LdapProvider
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultLdapProviderModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*LdapProvider, error) {
	var result []*LdapProvider
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultLdapProviderModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*LdapProvider, error) {
	var result []*LdapProvider
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultLdapProviderModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*LdapProvider, error) {
	var result []*LdapProvider
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
