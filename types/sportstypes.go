package types

import "github.com/musishere/sportsApp/internal/models"

// CreateSportRequest contains all parameters needed for creating a new sport
type CreateSportRequest struct {
	Name       string
	MinPlayers int
	MaxPlayers int
	FileBytes  []byte
	Filename   string
}

// UpdateSportRequest contains optional fields for updating a sport
type UpdateSportRequest struct {
	Name       *string `json:"name" form:"name"`
	MinPlayers *int    `json:"minPlayers" form:"minPlayers"`
	MaxPlayers *int    `json:"maxPlayers" form:"maxPlayers"`
}

// SportsResponse contains the response from sport operations
type SportsResponse struct {
	Sport   *models.Sports `json:"sport"`
	Message string         `json:"message,omitempty"`
}

// GetAllSportsResponse contains the response for getting all sports
type GetAllSportsResponse struct {
	Sport   *[]models.Sports `json:"sport"`
	Message string           `json:"message"`
}
