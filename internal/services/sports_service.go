package services

import "github.com/musishere/sportsApp/internal/repositories"

type SportsService struct {
	repo *repositories.SportsRepository
}

func NewSportsService(repo *repositories.SportsRepository) *SportsService {
	return &SportsService{
		repo: repo,
	}
}

// Todo
