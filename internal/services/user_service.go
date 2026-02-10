package services

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/utils"
	"github.com/musishere/sportsApp/internal/validators"
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

func (s *UserService) Register(name, email, password, gender, phone, cnic string, latitude, longitude float64) (*models.User, string, error) {
	if err := validators.ValidateRegisterInput(name, email, password, gender, phone, latitude, longitude); err != nil {
		return nil, "", err
	}

	existingUser, _ := s.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	cnicNumber, err := auth.HashedCnic(cnic)
	if err != nil {
		return nil, "", err
	}

	// 1. Generate OTP
	otpInt := helpers.GenerateOTP()
	otpStr := strconv.Itoa(otpInt)

	// 2. Store OTP and phone in Redis
	if err := helpers.StoreOTP(phone, otpStr); err != nil {
		return nil, "", err
	}

	// 3. Send OTP to phone
	log.Printf("OTP sent to %s: %s", phone, otpStr)
	if _, err := utils.SendExistingOTP(phone, otpStr); err != nil {
		return nil, "", err
	}

	// 4. Create user with is_active false until OTP is verified
	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Phone:     phone,
		Gender:    gender,
		Role:      "player",
		IsActive:  false,
		Cnic:      cnicNumber,
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

	user.Location = *location

	// No token until OTP verified; client must call verify-otp
	return user, "", nil
}

// ActivateUserByPhone sets user is_active to true after OTP was verified (call helpers.VerifyOTP first).
func (s *UserService) ActivateUserByPhone(phone string) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByPhone(phone)
	if err != nil {
		return nil, "", errors.New("user not found for this phone")
	}

	user.IsActive = true
	user.UpdatedAt = time.Now()
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, "", err
	}

	location, _ := s.locationRepo.GetLocationByUserID(user.ID.String())
	if location != nil {
		user.Location = *location
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Name, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(email, password string, latitude, longitude float64) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if !user.IsActive {
		return nil, "", errors.New("please verify your phone with OTP first")
	}

	if err := auth.VerifyPassword(user.Password, password); err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Name, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	location, err := s.locationRepo.GetLocationByUserID(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	location.Latitude = latitude
	location.Longitude = longitude

	if err := s.locationRepo.UpdateLocation(location); err != nil {
		return nil, "", err
	}

	user.Location = *location

	return user, token, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("Please provide an ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	location, err := s.locationRepo.GetLocationByUserID(user.ID.String())
	if err != nil {
		return nil, err
	}

	user.Location = *location

	return user, nil
}
