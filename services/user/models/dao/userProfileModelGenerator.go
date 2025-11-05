package dao

import (
	"context"

	"gorm.io/gorm"
)

type (
	profileModel interface {
		Create(ctx context.Context, data *Profile) error
		CreateInBatches(ctx context.Context, data []*Profile, batchSize int) (err error, RowsAffected int64)

		Update(ctx context.Context, data *Profile) error
		UpdateInBatches(ctx context.Context, data []*Profile) error

		Find(ctx context.Context, data *Profile, offset int, limit int) ([]*Profile, error)
		Count(ctx context.Context, data *Profile) (int64, error)

		Delete(ctx context.Context, data *Profile) error
		DeleteInBatches(ctx context.Context, data []*Profile) error
	}

	defaultProfileModel struct {
		db        *gorm.DB
		tableName string
	}

	Profile struct {
		BaseModel
		UseID    uint32 `gorm:"type:int;comment:用户ID"`
		Avatar   string `gorm:"column:avatar;type:varchar(255);comment:头像"`
		Locale   string `gorm:"column:locale;type:varchar(255);comment:地址"`
		Timezone string `gorm:"column:timezone;type:varchar(30);comment:时区"`
	}
)

func (emp Profile) TableName() string {
	return TablePrefix + "_" + "profiles"
}

func newProfileModel(conn *gorm.DB) *defaultProfileModel {
	return &defaultProfileModel{
		db:        conn,
		tableName: TablePrefix + "_" + "profiles",
	}
}

func (d defaultProfileModel) Create(ctx context.Context, data *Profile) error {
	return d.db.WithContext(ctx).Create(data).Error
}

func (d defaultProfileModel) CreateInBatches(ctx context.Context, data []*Profile, batchSize int) (err error, RowsAffected int64) {
	res := d.db.WithContext(ctx).CreateInBatches(data, batchSize)
	return res.Error, res.RowsAffected
}

func (d defaultProfileModel) Update(ctx context.Context, data *Profile) error {
	return d.db.WithContext(ctx).Updates(data).Error
}

func (d defaultProfileModel) UpdateInBatches(ctx context.Context, data []*Profile) error {
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

func (d defaultProfileModel) Find(ctx context.Context, data *Profile, offset int, limit int) ([]*Profile, error) {
	records := make([]*Profile, 0)
	res := d.db.WithContext(ctx).Where(data).Order("created_at desc").Offset(offset).Limit(limit).Find(&records)
	return records, res.Error
}

func (d defaultProfileModel) Count(ctx context.Context, data *Profile) (int64, error) {
	var count int64
	res := d.db.WithContext(ctx).Model(&Profile{}).Where(data).Count(&count)
	return count, res.Error
}

func (d defaultProfileModel) Delete(ctx context.Context, data *Profile) error {
	return d.db.WithContext(ctx).Delete(data).Error
}

func (d defaultProfileModel) DeleteInBatches(ctx context.Context, data []*Profile) error {
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
