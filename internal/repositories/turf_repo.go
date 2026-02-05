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
	return r.db.Create(turf).Error
}
