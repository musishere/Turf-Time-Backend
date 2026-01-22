package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
)

type UserService struct {
	userRepo     *repositories.UserRepository
	locationRepo *repositories.LocationRepository
	jwtSecret    string
}

func NewUserService(
	userRepo *repositories.UserRepository,
	locationRepo *repositories.LocationRepository,
	jwtSecret string,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		locationRepo: locationRepo,
		jwtSecret:    jwtSecret,
	}
}

func (s *UserService) Register(name, email, password string, latitude, longitude float64) (*models.User, string, error) {
	if name == "" || email == "" || password == "" {
		return nil, "", errors.New("name, email, and password are required")
	}

	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Role:      "player",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, "", err
	}

	location := &models.Location{
		ID:        uuid.New().String(),
		UserID:    user.ID.String(),
		Latitude:  latitude,
		Longitude: longitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.locationRepo.CreateLocation(location); err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Name, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(email, password string, latitude, longitude float64) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := auth.VerifyPassword(user.Password, password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if latitude != 0 && longitude != 0 {
		existingLocation, err := s.locationRepo.GetLocationByUserID(user.ID.String())
		if err != nil {
			location := &models.Location{
				ID:        uuid.New().String(),
				UserID:    user.ID.String(),
				Latitude:  latitude,
				Longitude: longitude,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := s.locationRepo.CreateLocation(location); err != nil {
				fmt.Printf("Failed to create location: %v\n", err)
				return nil, "", err
			}
		} else {
			existingLocation.Latitude = latitude
			existingLocation.Longitude = longitude
			existingLocation.UpdatedAt = time.Now()
			if err := s.locationRepo.UpdateLocation(existingLocation); err != nil {
				fmt.Printf("Failed to update location: %v\n", err)
				return nil, "", err
			}
		}
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Name, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
