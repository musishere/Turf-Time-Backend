package services

import (
	"context"
	"fmt"
	"time"

	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
)

type SportsService struct {
	repo *repositories.SportsRepository
}

func NewSportsService(repo *repositories.SportsRepository) *SportsService {
	return &SportsService{
		repo: repo,
	}
}

// SERVICE - Accept bytes instead
func (s *SportsService) CreateNewSport(
	name string,
	minPlayers, maxPlayers int,
	fileBytes []byte,
	filename string,
) (*models.Sports, error) {

	if err := helpers.ValidateSportInput(name, minPlayers, maxPlayers); err != nil {
		return nil, err
	}

	uploader, err := helpers.NewImageUploader()
	if err != nil {
		return nil, err
	}

	// Upload bytes instead of FileHeader
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // shorter timeout
	defer cancel()

	// Use UploadFromBytes instead (see below)
	imageURL, err := uploader.UploadFromBytes(ctx, fileBytes, filename, "sports")
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	sports := &models.Sports{
		Name:       name,
		MinPlayers: minPlayers,
		MaxPlayers: maxPlayers,
		IconUrl:    imageURL,
	}

	err = s.repo.CreateSport(sports)
	if err != nil {
		return nil, fmt.Errorf("Failed to create", err)
	}

	return sports, nil
}

func (s *SportsService) GetAllSports() (*models.Sports, error) { return &models.Sports{}, nil }

func (s *SportsService) GetSportsById(id string) (*models.Sports, error) {
	return &models.Sports{}, nil
}

func (s *SportsService) UpdateSports(id string) (*models.Sports, error) {
	return &models.Sports{}, nil

}

func (s *SportsService) DeleteSports(id string) (string, error) {
	return "", nil
}
