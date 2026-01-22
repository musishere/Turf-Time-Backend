package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
)

type LocationService struct {
	repo *repositories.LocationRepository
}

func NewLocationService(repo *repositories.LocationRepository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) CreateLocationForUser(userID uuid.UUID, latitude, longitude float64) (*models.Location, error) {
	if latitude < -90 || latitude > 90 {
		return nil, errors.New("invalid latitude")
	}
	if longitude < -180 || longitude > 180 {
		return nil, errors.New("invalid longitude")
	}

	location := &models.Location{
		ID:        uuid.New().String(),
		UserID:    userID.String(),
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateLocation(location); err != nil {
		return nil, err
	}

	return location, nil
}
