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
