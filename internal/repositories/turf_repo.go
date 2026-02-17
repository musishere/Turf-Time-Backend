package repositories

import (
	"github.com/musishere/sportsApp/internal/models"
	"gorm.io/gorm"
)

type TurfRepostitory struct {
	db *gorm.DB
}

func NewTurfRepository(db *gorm.DB) *TurfRepostitory {
	return &TurfRepostitory{
		db: db,
	}
}

func (r *TurfRepostitory) Create(turf *models.Turf) error {
	if err := r.db.Create(turf).Error; err != nil {
		return err
	}
	return r.db.Preload("Owner").Preload("Owner.Location").First(turf, turf.ID).Error
}

func (r *TurfRepostitory) GetAllTurfsRepo(page, pageSize int) ([]models.Turf, int64, error) {
	var turfs []models.Turf
	var total int64

	if err := r.db.Model(&models.Turf{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	if err := r.db.Preload("Owner").Preload("Owner.Location").Offset(offset).Limit(pageSize).Find(&turfs).Error; err != nil {
		return nil, 0, err
	}

	return turfs, total, nil
}

func (r *TurfRepostitory) GetTurfByID(id string) (*models.Turf, error) {
	var turf models.Turf
	if err := r.db.Preload("Owner").Preload("Owner.Location").First(&turf, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &turf, nil
}

func (r *TurfRepostitory) UpdateTurf(turf *models.Turf) error {
	return r.db.Save(turf).Error
}

func (r *TurfRepostitory) DeleteTurf(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Turf{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
