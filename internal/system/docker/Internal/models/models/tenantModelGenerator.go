package models

import (
	"context"

	"gorm.io/gorm"
)

const (
	Team = iota + 1
	Enterprise
)

type (
	tenantModel interface {
		Insert(ctx context.Context, data *Tenant) error
		Find(ctx context.Context, data interface{}) (interface{}, error)
		FindOne(ctx context.Context, id int64) (*Tenant, error)
		FindAll(ctx context.Context, limit int) ([]*Tenant, error)
		FindByFields(ctx context.Context, fields interface{}, limit int) ([]*Tenant, error)
		Count(ctx context.Context) (int64, error)
		Update(ctx context.Context, data *Tenant) error
	}

	defaultTenantModel struct {
		db        *gorm.DB
		tableName string
	}

	TenantLevel uint16

	Tenant struct {
		BaseModel
		Name        uint32      `gorm:"column:name;type:varchar(30);comment:租户名称"`
		TenantName  uint32      `gorm:"column:tenant_name;type:varchar(30);not null;unique;comment:租户名(唯一)"`
		TenantLevel TenantLevel `gorm:"column:type;type:int;comment:租户类型"`
	}
)

var TenantLevelMap = map[uint16]string{
	Team:       "team",
	Enterprise: "enterprise",
}

func (tp TenantLevel) Desc() string {
	return TenantLevelMap[uint16(tp)]
}

func (emp Tenant) TableName() string {
	return TablePrefix + "_" + "tenants"
}

func newTenantModel(conn *gorm.DB) *defaultTenantModel {
	return &defaultTenantModel{
		db:        conn,
		tableName: TablePrefix + "_" + "tenants",
	}
}

func (m *defaultTenantModel) Insert(ctx context.Context, data *Tenant) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *defaultTenantModel) Find(ctx context.Context, data interface{}) (interface{}, error) {
	u := data
	err := m.db.WithContext(ctx).Find(u).Error
	return u, err
}

func (m *defaultTenantModel) FindOne(ctx context.Context, id int64) (*Tenant, error) {
	var result Tenant
	err := m.db.WithContext(ctx).Model(&Tenant{}).Where("id = ?", id).First(&result).Error
	return &result, err
}

func (m *defaultTenantModel) FindByFields(ctx context.Context, fields interface{}, limit int) ([]*Tenant, error) {
	var result []*Tenant
	err := m.db.WithContext(ctx).Model(&Tenant{}).Where(fields).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultTenantModel) FindAll(ctx context.Context, limit int) ([]*Tenant, error) {
	var result []*Tenant
	err := m.db.WithContext(ctx).Model(&Tenant{}).Limit(limit).Find(&result).Error
	return result, err
}

func (m *defaultTenantModel) Count(ctx context.Context) (int64, error) {
	var result int64
	err := m.db.WithContext(ctx).Model(&Tenant{}).Count(&result).Error
	return result, err
}

func (m *defaultTenantModel) Update(ctx context.Context, data *Tenant) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *defaultTenantModel) UpdateFields(ctx context.Context, id int64, values map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&Tenant{}).Where("id = ?", id).Updates(values).Error
}

func (m *defaultTenantModel) FindByUserProfileIDAndPassword(ctx context.Context, username string, password string) (*Tenant, error) {
	var result Tenant
	err := m.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &result, err
}

func (m *defaultTenantModel) FindByUserProfileId(ctx context.Context, userId int64, limit int) ([]*Tenant, error) {
	var result []*Tenant
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error

	return result, err
}

func (m *defaultTenantModel) FindByFollowedUserProfileIds(ctx context.Context, userId int64, followedUserProfileIds []int64) ([]*Tenant, error) {
	var result []*Tenant
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserProfileIds).
		Find(&result).Error

	return result, err
}

func (m *defaultTenantModel) FindByFollowedUserProfileId(ctx context.Context, userId int64, limit int) ([]*Tenant, error) {
	var result []*Tenant
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}
