package services

import (
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
)

type TurfService struct {
	repo *repositories.TurfRepostitory
}

func NewTurfService(repo *repositories.TurfRepostitory) *TurfService {
	return &TurfService{
		repo: repo,
	}
}

func (r *TurfService) CreateTurf() (*models.Turf, error)   { return &models.Turf{}, nil }
func (r *TurfService) GetTurfByID() (*models.Turf, error)  { return &models.Turf{}, nil }
func (r *TurfService) GetAllTurf() (*[]models.Turf, error) { return &[]models.Turf{}, nil }
func (r *TurfService) UpdateTurf() (*models.Turf, error)   { return &models.Turf{}, nil }
func (r *TurfService) DeleteTurf() (string, error)         { return "", nil }
