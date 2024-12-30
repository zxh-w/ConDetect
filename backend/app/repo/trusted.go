package repo

import (
	"ConDetect/backend/app/model"
	"ConDetect/backend/global"
)

type TrustedRepo struct{}

type ITrustedRepo interface {
	Page(limit, offset int, opts ...DBOption) (int64, []model.MeasureRecord, error)
	Create(MeasureRecord *model.MeasureRecord) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error
	Get(opts ...DBOption) (model.MeasureRecord, error)
	List(opts ...DBOption) ([]model.MeasureRecord, error)
}

func NewITrustedRepo() ITrustedRepo {
	return &TrustedRepo{}
}

func (u *TrustedRepo) Get(opts ...DBOption) (model.MeasureRecord, error) {
	var MeasureRecord model.MeasureRecord
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&MeasureRecord).Error
	return MeasureRecord, err
}

func (u *TrustedRepo) List(opts ...DBOption) ([]model.MeasureRecord, error) {
	var MeasureRecord []model.MeasureRecord
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&MeasureRecord).Error
	return MeasureRecord, err
}

func (u *TrustedRepo) Page(page, size int, opts ...DBOption) (int64, []model.MeasureRecord, error) {
	var users []model.MeasureRecord
	db := global.DB.Model(&model.MeasureRecord{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *TrustedRepo) Create(MeasureRecord *model.MeasureRecord) error {
	return global.DB.Create(MeasureRecord).Error
}

func (u *TrustedRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&model.MeasureRecord{}).Where("id = ?", id).Updates(vars).Error
}

func (u *TrustedRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.MeasureRecord{}).Error
}
