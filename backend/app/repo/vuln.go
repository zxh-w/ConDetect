package repo

import (
	"ConDetect/backend/app/model"
	"ConDetect/backend/global"
)

type VulnRepo struct{}

type IVulnRepo interface {
	Page(limit, offset int, opts ...DBOption) (int64, []model.Trivy, error)
	Create(Trivy *model.Trivy) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error
	Get(opts ...DBOption) (model.Trivy, error)
	List(opts ...DBOption) ([]model.Trivy, error)
}

func NewIVulnRepo() IVulnRepo {
	return &VulnRepo{}
}

func (u *VulnRepo) Get(opts ...DBOption) (model.Trivy, error) {
	var Trivy model.Trivy
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&Trivy).Error
	return Trivy, err
}

func (u *VulnRepo) List(opts ...DBOption) ([]model.Trivy, error) {
	var Trivy []model.Trivy
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&Trivy).Error
	return Trivy, err
}

func (u *VulnRepo) Page(page, size int, opts ...DBOption) (int64, []model.Trivy, error) {
	var users []model.Trivy
	db := global.DB.Model(&model.Trivy{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *VulnRepo) Create(Trivy *model.Trivy) error {
	return global.DB.Create(Trivy).Error
}

func (u *VulnRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&model.Trivy{}).Where("id = ?", id).Updates(vars).Error
}

func (u *VulnRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.Trivy{}).Error
}
