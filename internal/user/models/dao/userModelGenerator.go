package dao

import (
	"context"
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

type Identity string

const (
	ADMIN          Identity = "admin"
	ORDINARY       Identity = "ordinary"
	SERVICEACCOUNT Identity = "serviceaccount"
)

var (
	identityToUintMap = map[Identity]uint16{
		ADMIN:          1,
		ORDINARY:       2,
		SERVICEACCOUNT: 3,
	}

	uintToIdentityMap = map[uint16]Identity{
		1: ADMIN,
		2: ORDINARY,
		3: SERVICEACCOUNT,
	}
)

func (l *Identity) Value() (driver.Value, error) {
	if val, ok := identityToUintMap[*l]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("invalid identity: %v", l)
}

func (l *Identity) Scan(value interface{}) error {
	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid identity db value: %v", value)
	}
	if strVal, ok := uintToIdentityMap[uint16(intVal)]; ok {
		*l = strVal
		return nil
	}
	return fmt.Errorf("unknown identity: %d", intVal)
}

type Status string

const (
	ACTIVITY Status = "activity"
	DISABLED Status = "disabled"
	INACTIVE Status = "inactive"
)

var (
	statusToUintMap = map[Status]uint16{
		ACTIVITY: 1,
		DISABLED: 2,
		INACTIVE: 3,
	}

	uintToStatusMap = map[uint16]Status{
		1: ACTIVITY,
		2: DISABLED,
		3: INACTIVE,
	}
)

func (s *Status) Value() (driver.Value, error) {
	if val, ok := statusToUintMap[*s]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("invalid status: %v", s)
}

func (s *Status) Scan(value interface{}) error {
	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid status db value: %v", value)
	}
	if strVal, ok := uintToStatusMap[uint16(intVal)]; ok {
		*s = strVal
		return nil
	}
	return fmt.Errorf("unknown status: %d", intVal)
}

type AuthProvider string

const (
	LOCAL  AuthProvider = "local"
	LDAP   AuthProvider = "ldap"
	OAUTH2 AuthProvider = "oauth2"
)

var (
	authProviderToUintMap = map[AuthProvider]uint16{
		LOCAL:  1,
		LDAP:   2,
		OAUTH2: 3,
	}

	uintToAuthProviderMap = map[uint16]AuthProvider{
		1: LOCAL,
		2: LDAP,
		3: OAUTH2,
	}
)

func (ap *AuthProvider) Value() (driver.Value, error) {
	if val, ok := authProviderToUintMap[*ap]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("invalid auth provider: %v", ap)
}

func (ap *AuthProvider) Scan(value interface{}) error {
	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid auth provider db value: %v", value)
	}
	if strVal, ok := uintToAuthProviderMap[uint16(intVal)]; ok {
		*ap = strVal
		return nil
	}
	return fmt.Errorf("unknown auth provider: %d", intVal)
}

type (
	userModel interface {
		Create(ctx context.Context, data *User) error
		CreateInBatches(ctx context.Context, data []*User, batchSize int) (err error, RowsAffected int64)

		Update(ctx context.Context, data *User) error
		UpdateInBatches(ctx context.Context, data []*User) error

		Find(ctx context.Context, data *User, offset int, limit int) ([]*User, error)
		Count(ctx context.Context, data *User) (int64, error)

		Delete(ctx context.Context, data *User) error
		DeleteInBatches(ctx context.Context, data []*User) error
	}

	defaultUserModel struct {
		db        *gorm.DB
		tableName string
	}

	User struct {
		BaseModel
		Name           *string       `gorm:"column:name;type:varchar(30);not null;comment:显示名"`
		UserName       *string       `gorm:"column:username;type:varchar(30);not null;uniqueIndex:username_auth_provider;comment:用户名(唯一)"`
		Password       *string       `gorm:"column:password;type:varchar(100);comment:密码(身份提供商不是Local时可以为空)"`
		CountryCode    *string       `gorm:"column:country_code;type:varchar(60);comment:手机区号"`
		Phone          *string       `gorm:"column:phone;type:varchar(60);comment:手机号"`
		Email          *string       `gorm:"column:email;type:varchar(60);comment:邮箱"`
		Identity       *Identity     `gorm:"column:identity;type:smallint;not null;comment:用户类型"`
		Status         *Status       `gorm:"column:status;type:smallint;not null;comment:用户状态"`
		AuthProvider   *AuthProvider `gorm:"column:auth_provider;type:smallint;;not null;uniqueIndex:username_auth_provider;comment:身份提供商"`
		AuthProviderID *uint32       `gorm:"column:auth_provider_id;type:int;uniqueIndex:username_auth_provider;comment:身份提供商ID"`
	}
)

func (emp User) TableName() string {
	return TablePrefix + "_" + "users"
}

func newUserModel(conn *gorm.DB) *defaultUserModel {
	return &defaultUserModel{
		db:        conn,
		tableName: TablePrefix + "_" + "users",
	}
}

func (d defaultUserModel) Create(ctx context.Context, data *User) error {
	return d.db.WithContext(ctx).Create(data).Error
}

func (d defaultUserModel) CreateInBatches(ctx context.Context, data []*User, batchSize int) (err error, RowsAffected int64) {
	res := d.db.WithContext(ctx).CreateInBatches(data, batchSize)
	return res.Error, res.RowsAffected
}

func (d defaultUserModel) Update(ctx context.Context, data *User) error {
	return d.db.WithContext(ctx).Updates(data).Error
}

func (d defaultUserModel) UpdateInBatches(ctx context.Context, data []*User) error {
	callFc := func(tx *gorm.DB) error {
		for _, userGroup := range data {
			if err := tx.WithContext(ctx).Updates(userGroup).Error; err != nil {
				return err
			}
		}
		return nil
	}
	return d.db.WithContext(ctx).Transaction(callFc)
}

func (d defaultUserModel) Find(ctx context.Context, data *User, offset int, limit int) ([]*User, error) {
	records := make([]*User, 0)
	res := d.db.WithContext(ctx).Where(data).Order("created_at desc").Offset(offset).Limit(limit).Find(&records)
	return records, res.Error
}

func (d defaultUserModel) Count(ctx context.Context, data *User) (int64, error) {
	var count int64
	res := d.db.WithContext(ctx).Model(&User{}).Where(data).Count(&count)
	return count, res.Error
}

func (d defaultUserModel) Delete(ctx context.Context, data *User) error {
	return d.db.WithContext(ctx).Delete(data).Error
}

func (d defaultUserModel) DeleteInBatches(ctx context.Context, data []*User) error {
	callFc := func(tx *gorm.DB) error {
		for _, userGroup := range data {
			if err := tx.WithContext(ctx).Delete(userGroup).Error; err != nil {
				return err
			}
		}
		return nil
	}
	return d.db.WithContext(ctx).Transaction(callFc)
}
