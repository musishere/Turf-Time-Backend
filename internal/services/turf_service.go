package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/validators"
)

type TurfService struct {
	repo     *repositories.TurfRepostitory
	uploader *helpers.ImageUploader
}

func NewTurfService(repo *repositories.TurfRepostitory, uploader *helpers.ImageUploader) *TurfService {
	return &TurfService{
		repo:     repo,
		uploader: uploader,
	}
}

// CreateTurf creates a turf with 3 required images (uploaded to Cloudinary).
func (s *TurfService) CreateTurf(
	name string,
	startTime, endTime int,
	status string,
	noOfFields int,
	address string,
	ownerID uuid.UUID,
	img1, img2, img3 []byte,
	filename1, filename2, filename3 string,
) (*models.Turf, error) {
	if status == "" {
		status = "active"
	}
	if err := validators.ValidateTurfInput(name, startTime, endTime, status, noOfFields, address); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	upload := func(data []byte, filename string) (string, error) {
		return s.uploader.UploadFromBytes(ctx, data, filename, "turfs")
	}

	var url1, url2, url3 string
	var err1, err2, err3 error
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { defer wg.Done(); url1, err1 = upload(img1, filename1) }()
	go func() { defer wg.Done(); url2, err2 = upload(img2, filename2) }()
	go func() { defer wg.Done(); url3, err3 = upload(img3, filename3) }()
	wg.Wait()

	if err1 != nil {
		return nil, fmt.Errorf("upload image 1: %w", err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("upload image 2: %w", err2)
	}
	if err3 != nil {
		return nil, fmt.Errorf("upload image 3: %w", err3)
	}

	turf := &models.Turf{
		Name:       name,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     status,
		NoOfFields: noOfFields,
		Address:    address,
		TurfImages: []string{url1, url2, url3},
		OwnerID:    ownerID,
	}
	if err := s.repo.Create(turf); err != nil {
		return nil, fmt.Errorf("failed to create turf: %w", err)
	}
	return turf, nil
}

func (r *TurfService) GetTurfByID() (*models.Turf, error)  { return &models.Turf{}, nil }
func (r *TurfService) GetAllTurf() (*[]models.Turf, error) { return &[]models.Turf{}, nil }
func (r *TurfService) UpdateTurf() (*models.Turf, error)   { return &models.Turf{}, nil }
func (r *TurfService) DeleteTurf() (string, error)         { return "", nil }
