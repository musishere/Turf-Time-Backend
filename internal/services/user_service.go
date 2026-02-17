package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/validators"
	"github.com/musishere/sportsApp/types"
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

func (s *UserService) Register(req types.RegisterRequest) (*models.User, string, error) {
	if err := validators.ValidateRegisterInput(req.Name, req.Email, req.Password, req.Gender, req.Phone, req.Latitude, req.Longitude); err != nil {
		return nil, "", err
	}

	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	existingPhoneNumber, _ := s.userRepo.GetUserByPhone(req.Phone)
	if existingPhoneNumber != nil {
		return nil, "", errors.New("phone number already registered")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	cnicNumber, err := auth.HashedCnic(req.Cnic)
	if err != nil {
		return nil, "", err
	}

	// 1. Generate OTP
	// otpInt := helpers.GenerateOTP()
	// fmt.Print(otpInt)
	// otpStr := strconv.Itoa(otpInt)

	// OTP flow commented out for testing - create user directly with is_active true

	// 2. Store OTP and phone in Redis
	// if err := helpers.StoreOTP(phone, otpStr); err != nil {
	// 	return nil, "", err
	// }
	// 3. Send OTP to phone
	// log.Printf("OTP sent to %s: %s", phone, otpStr)
	// if _, err := utils.SendExistingOTP(phone, otpStr); err != nil {
	// 	return nil, "", err
	// }

	// Create user with is_active true (OTP verification disabled for testing)
	user := &models.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Role:      "player",
		IsActive:  true, // Set to true for testing without OTP
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
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.locationRepo.CreateLocation(location); err != nil {
		return nil, "", err
	}

	user.Location = *location

	// Generate token immediately (OTP verification disabled for testing)
	token, err := auth.GenerateJWT(user.ID, user.Email, user.Name, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ActivateUserByPhone sets user is_active to true after OTP was verified (call helpers.VerifyOTP first).
func (s *UserService) ActivateUserByPhone(req types.ActivateUserRequest) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByPhone(req.Phone)
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

// ActivateUserByEmail sets user is_active to true after OTP was verified via email (call helpers.VerifyOTP first).
func (s *UserService) ActivateUserByEmail(email string) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("user not found for this email")
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

func (s *UserService) Login(req types.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, "", err
	}

	if !user.IsActive {
		return nil, "", errors.New("please verify your phone with OTP first")
	}

	if err := auth.VerifyPassword(user.Password, req.Password); err != nil {
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

	location.Latitude = req.Latitude
	location.Longitude = req.Longitude

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
