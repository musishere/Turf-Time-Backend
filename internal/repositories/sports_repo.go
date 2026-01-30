package repositories

import (
	"github.com/musishere/sportsApp/internal/models"
	"gorm.io/gorm"
)

type SportsRepository struct {
	db *gorm.DB
}

func NewSportsRepositry(db *gorm.DB) *SportsRepository {
	return &SportsRepository{
		db: db,
	}
}

func (r *SportsRepository) CreateSport(sports *models.Sports) error {
	return r.db.Create(sports).Error
}
