package dao

import (
	"context"

	"gorm.io/gorm"
)

type (
	userGroupModel interface {
		Create(ctx context.Context, data *UserGroup) error
		CreateInBatches(ctx context.Context, data []*UserGroup, batchSize int) (err error, RowsAffected int64)

		Update(ctx context.Context, data *UserGroup) error
		UpdateInBatches(ctx context.Context, data []*UserGroup) error

		Find(ctx context.Context, data *UserGroup, offset int, limit int) ([]*UserGroup, error)
		Count(ctx context.Context, data *UserGroup) (int64, error)

		Delete(ctx context.Context, data *UserGroup) error
		DeleteInBatches(ctx context.Context, data []*UserGroup) error
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

func (d defaultUserGroupModel) Create(ctx context.Context, data *UserGroup) error {
	return d.db.WithContext(ctx).Create(data).Error
}

func (d defaultUserGroupModel) CreateInBatches(ctx context.Context, data []*UserGroup, batchSize int) (err error, RowsAffected int64) {
	res := d.db.WithContext(ctx).CreateInBatches(data, batchSize)
	return res.Error, res.RowsAffected
}

func (d defaultUserGroupModel) Update(ctx context.Context, data *UserGroup) error {
	return d.db.WithContext(ctx).Updates(data).Error
}

func (d defaultUserGroupModel) UpdateInBatches(ctx context.Context, data []*UserGroup) error {
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

func (d defaultUserGroupModel) Find(ctx context.Context, data *UserGroup, offset int, limit int) ([]*UserGroup, error) {
	records := make([]*UserGroup, 0)
	res := d.db.WithContext(ctx).Where(data).Order("created_at desc").Offset(offset).Limit(limit).Find(&records)
	return records, res.Error
}

func (d defaultUserGroupModel) Count(ctx context.Context, data *UserGroup) (int64, error) {
	var count int64
	res := d.db.WithContext(ctx).Model(&UserGroup{}).Where(data).Count(&count)
	return count, res.Error
}

func (d defaultUserGroupModel) Delete(ctx context.Context, data *UserGroup) error {
	return d.db.WithContext(ctx).Delete(data).Error
}

func (d defaultUserGroupModel) DeleteInBatches(ctx context.Context, data []*UserGroup) error {
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
