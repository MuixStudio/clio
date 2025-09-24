package dao

import (
	"context"

	"gorm.io/gorm"
)

const (
	LOCAL = iota + 1
	LDAP
	OAUTH2
)

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

	AuthProvider uint16

	User struct {
		BaseModel
		Name           *string      `gorm:"column:name;type:varchar(30);not null;comment:显示名"`
		UserName       *string      `gorm:"column:username;type:varchar(30);not null;uniqueIndex:username_auth_provider;comment:用户名(唯一)"`
		Password       string       `gorm:"column:password;type:varchar(100);comment:密码(身份提供商不是Local时可以为空)"`
		CountryCode    string       `gorm:"column:country_code;type:varchar(60);comment:手机区号"`
		Phone          string       `gorm:"column:phone;type:varchar(60);comment:手机号"`
		Email          string       `gorm:"column:email;type:varchar(60);comment:邮箱"`
		IsAdmin        bool         `gorm:"column:is_admin;type:boolean;not null;comment:是否是管理员"`
		AuthProvider   AuthProvider `gorm:"column:auth_provider;type:smallint;;not null;uniqueIndex:username_auth_provider;comment:身份提供商"`
		AuthProviderID uint32       `gorm:"column:auth_provider_id;type:int;uniqueIndex:username_auth_provider;comment:身份提供商ID"`
	}
)

var AuthProviderMap = map[uint16]string{
	LOCAL:  "Local",
	LDAP:   "Ldap",
	OAUTH2: "Oauth2",
}

func (ap AuthProvider) Desc() string {
	return AuthProviderMap[uint16(ap)]
}

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
