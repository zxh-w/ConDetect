package repo

import (
	"ConDetect/backend/app/model"
	"ConDetect/backend/global"
)

type MaciousRepo struct{}

type IMaciousRepo interface {
	Page(limit, offset int, opts ...DBOption) (int64, []model.Sandfly, error)
	Create(Sandfly *model.Sandfly) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error
	Get(opts ...DBOption) (model.Sandfly, error)
	List(opts ...DBOption) ([]model.Sandfly, error)
}

func NewIMaciousRepo() IMaciousRepo {
	return &MaciousRepo{}
}

func (u *MaciousRepo) Get(opts ...DBOption) (model.Sandfly, error) {
	var Sandfly model.Sandfly
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&Sandfly).Error
	return Sandfly, err
}

func (u *MaciousRepo) List(opts ...DBOption) ([]model.Sandfly, error) {
	var Sandfly []model.Sandfly
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&Sandfly).Error
	return Sandfly, err
}

func (u *MaciousRepo) Page(page, size int, opts ...DBOption) (int64, []model.Sandfly, error) {
	var users []model.Sandfly
	db := global.DB.Model(&model.Sandfly{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *MaciousRepo) Create(Sandfly *model.Sandfly) error {
	return global.DB.Create(Sandfly).Error
}

func (u *MaciousRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&model.Sandfly{}).Where("id = ?", id).Updates(vars).Error
}

func (u *MaciousRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.Sandfly{}).Error
}
