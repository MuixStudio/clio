package dao

import (
	"context"

	"gorm.io/gorm"
)

type (
	userUserGroupModel interface {
		Create(ctx context.Context, data *UserUserGroup) error
		CreateInBatches(ctx context.Context, data []*UserUserGroup, batchSize int) (err error, RowsAffected int64)

		Update(ctx context.Context, data *UserUserGroup) error
		UpdateInBatches(ctx context.Context, data []*UserUserGroup) error

		Find(ctx context.Context, data *UserUserGroup, offset int, limit int) ([]*UserUserGroup, error)
		Count(ctx context.Context, data *UserUserGroup) (int64, error)

		Delete(ctx context.Context, data *UserUserGroup) error
		DeleteInBatches(ctx context.Context, data []*UserUserGroup) error
	}

	defaultUserUserGroupModel struct {
		db        *gorm.DB
		tableName string
	}

	UserUserGroup struct {
		BaseModel
		UseID      uint32 `gorm:"column:user_id;type:int;not null;comment:用户ID"`
		UseGroupID uint32 `gorm:"column:usergroup_id;type:int;not null;comment:用户组ID"`
	}
)

func (emp UserUserGroup) TableName() string {
	return TablePrefix + "_" + "user_usergroups"
}

func newUserUserGroupModel(conn *gorm.DB) *defaultUserUserGroupModel {
	return &defaultUserUserGroupModel{
		db:        conn,
		tableName: TablePrefix + "_" + "user_usergroups",
	}
}

func (d defaultUserUserGroupModel) Create(ctx context.Context, data *UserUserGroup) error {
	return d.db.WithContext(ctx).Create(data).Error
}

func (d defaultUserUserGroupModel) CreateInBatches(ctx context.Context, data []*UserUserGroup, batchSize int) (err error, RowsAffected int64) {
	res := d.db.WithContext(ctx).CreateInBatches(data, batchSize)
	return res.Error, res.RowsAffected
}

func (d defaultUserUserGroupModel) Update(ctx context.Context, data *UserUserGroup) error {
	return d.db.WithContext(ctx).Updates(data).Error
}

func (d defaultUserUserGroupModel) UpdateInBatches(ctx context.Context, data []*UserUserGroup) error {
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

func (d defaultUserUserGroupModel) Find(ctx context.Context, data *UserUserGroup, offset int, limit int) ([]*UserUserGroup, error) {
	records := make([]*UserUserGroup, 0)
	res := d.db.WithContext(ctx).Where(data).Order("created_at desc").Offset(offset).Limit(limit).Find(&records)
	return records, res.Error
}

func (d defaultUserUserGroupModel) Count(ctx context.Context, data *UserUserGroup) (int64, error) {
	var count int64
	res := d.db.WithContext(ctx).Model(&UserUserGroup{}).Where(data).Count(&count)
	return count, res.Error
}

func (d defaultUserUserGroupModel) Delete(ctx context.Context, data *UserUserGroup) error {
	return d.db.WithContext(ctx).Delete(data).Error
}

func (d defaultUserUserGroupModel) DeleteInBatches(ctx context.Context, data []*UserUserGroup) error {
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
