package validators

import (
	"errors"
	"strings"
)

const (
	minPasswordLength = 6
	minLatitude       = -90
	maxLatitude       = 90
	minLongitude      = -180
	maxLongitude      = 180
)

// ValidationError represents a validation failure so handlers can return 400.
type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "validation failed"
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ValidateRegisterInput validates user registration input.
// Name, email, password, gender, and phone are required. Password must be at least 6 characters.
// Latitude/longitude must be within valid geographic ranges.
// Returns a *ValidationError so callers can use errors.As to respond with 400.
func ValidateRegisterInput(name, email, password, gender, phone string, latitude, longitude float64) error {
	if strings.TrimSpace(name) == "" {
		return &ValidationError{Err: errors.New("name is required")}
	}
	if strings.TrimSpace(email) == "" {
		return &ValidationError{Err: errors.New("email is required")}
	}
	if strings.TrimSpace(password) == "" {
		return &ValidationError{Err: errors.New("password is required")}
	}
	if len(password) < minPasswordLength {
		return &ValidationError{Err: errors.New("password must be at least 6 characters")}
	}
	if strings.TrimSpace(gender) == "" {
		return &ValidationError{Err: errors.New("gender is required")}
	}
	if strings.TrimSpace(phone) == "" {
		return &ValidationError{Err: errors.New("phone is required")}
	}
	if latitude < minLatitude || latitude > maxLatitude {
		return &ValidationError{Err: errors.New("latitude must be between -90 and 90")}
	}
	if longitude < minLongitude || longitude > maxLongitude {
		return &ValidationError{Err: errors.New("longitude must be between -180 and 180")}
	}
	return nil
}
