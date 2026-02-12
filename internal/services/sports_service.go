package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/validators"
	"github.com/musishere/sportsApp/types"
)

type SportsService struct {
	repo     *repositories.SportsRepository
	uploader *helpers.ImageUploader
}

func NewSportsService(repo *repositories.SportsRepository, uploader *helpers.ImageUploader) *SportsService {
	return &SportsService{
		repo:     repo,
		uploader: uploader,
	}
}

func (s *SportsService) CreateNewSport(req types.CreateSportRequest) (*models.Sports, error) {
	if err := validators.ValidateSportInput(req.Name, req.MinPlayers, req.MaxPlayers); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	imageURL, err := s.uploader.UploadFromBytes(ctx, req.FileBytes, req.Filename, "sports")
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	sports := &models.Sports{
		Name:       req.Name,
		MinPlayers: req.MinPlayers,
		MaxPlayers: req.MaxPlayers,
		IconUrl:    imageURL,
	}

	err = s.repo.CreateSport(sports)
	if err != nil {
		return nil, fmt.Errorf("failed to create sport: %w", err)
	}

	return sports, nil
}

func (s *SportsService) GetAllSports() (*[]models.Sports, error) {

	sports := s.repo.GetSports()

	if len(sports) == 0 {
		return nil, errors.New("No sports found")
	}

	return &sports, nil
}

func (s *SportsService) GetSportsByID(id string) (*models.Sports, error) {
	sports, err := s.repo.GetSportsByID(id)
	if err != nil {
		return nil, err
	}

	return &sports, nil
}

// UpdateSports updates only the fields that are non-nil (partial update).
func (s *SportsService) UpdateSports(id string, req types.UpdateSportRequest) (*models.Sports, error) {
	sport, err := s.repo.GetSportsByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		sport.Name = *req.Name
	}
	if req.MinPlayers != nil {
		sport.MinPlayers = *req.MinPlayers
	}
	if req.MaxPlayers != nil {
		sport.MaxPlayers = *req.MaxPlayers
	}
	if err := validators.ValidateSportInput(sport.Name, sport.MinPlayers, sport.MaxPlayers); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSport(&sport); err != nil {
		return nil, fmt.Errorf("failed to update sport: %w", err)
	}
	return &sport, nil
}

func (s *SportsService) DeleteSports(id string) error {
	return s.repo.DeleteSport(id)
}
