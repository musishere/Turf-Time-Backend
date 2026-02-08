package repositories

import (
	"github.com/musishere/sportsApp/internal/models"
	"gorm.io/gorm"
)

type SportsRepository struct {
	db *gorm.DB
}

func NewSportsRepository(db *gorm.DB) *SportsRepository {
	return &SportsRepository{
		db: db,
	}
}

func (r *SportsRepository) CreateSport(sports *models.Sports) error {
	return r.db.Create(sports).Error
}

func (r *SportsRepository) GetSports() []models.Sports {
	var sports []models.Sports
	r.db.Find(&sports)
	return sports
}

func (r *SportsRepository) GetSportsByID(id string) (models.Sports, error) {
	var sport models.Sports
	result := r.db.Where("id = ?", id).First(&sport)
	if result.Error != nil {
		return models.Sports{}, result.Error
	}
	return sport, nil
}

func (r *SportsRepository) UpdateSport(sport *models.Sports) error {
	return r.db.Save(sport).Error
}

func (r *SportsRepository) DeleteSport(id string) error {
	result := r.db.Where("id = ?", id).Delete(&models.Sports{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
