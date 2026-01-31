package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/services"
	"gorm.io/gorm"
)

type SportsHandler struct {
	SportsService services.SportsService
}

type SportsRequest struct {
	Name       string                `form:"name" binding:"required"`
	IconUrl    *multipart.FileHeader `form:"iconUrl" binding:"required"`
	MinPlayers int                   `form:"minPlayer" binding:"required"`
	MaxPlayers int                   `form:"maxPlayer" binding:"required"`
}

// UpdateSportRequest: all fields optional — send only what you want to change.
type UpdateSportRequest struct {
	Name       *string `json:"name" form:"name"`
	MinPlayers *int    `json:"minPlayers" form:"minPlayers"`
	MaxPlayers *int    `json:"maxPlayers" form:"maxPlayers"`
}

type SportsResponse struct {
	Sport   *models.Sports
	Message string
}

type GetAllSportsResponse struct {
	Sport   *[]models.Sports
	Message string
}

func NewSportsHandler(sportsService *services.SportsService) *SportsHandler {
	return &SportsHandler{
		SportsService: *sportsService,
	}
}

func (s *SportsHandler) RegisterNewSports(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20) // 32MB

	name := c.PostForm("name")
	minPlayers, _ := strconv.Atoi(c.PostForm("minPlayers"))
	maxPlayers, _ := strconv.Atoi(c.PostForm("maxPlayers"))

	// Get file header
	fileHeader, err := c.FormFile("iconUrl")
	if err != nil {
		c.JSON(400, gin.H{"error": "iconUrl file is required"})
		return
	}

	// OPEN FILE HERE (critical fix!)
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot read file"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot read file data"})
		return
	}

	sport, err := s.SportsService.CreateNewSport(name, minPlayers, maxPlayers, fileBytes, fileHeader.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, sport)
}

func (s *SportsHandler) GetAllRegisteredSports(c *gin.Context) {

	sports, err := s.SportsService.GetAllSports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Finding sports",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetAllSportsResponse{
		Message: "Sports Fetched Successfully",
		Sport:   sports,
	})

}

func (s *SportsHandler) GetRegisteredSportsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Empty ID",
		})
		return
	}

	sports, err := s.SportsService.GetSportsById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Not Found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SportsResponse{
		Sport: sports,
	})

}

func (s *SportsHandler) UpdateRegisteredSports(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty ID"})
		return
	}

	var req UpdateSportRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	// At least one field must be sent
	if req.Name == nil && req.MinPlayers == nil && req.MaxPlayers == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "send at least one field to update: name, minPlayers, or maxPlayers"})
		return
	}

	sport, err := s.SportsService.UpdateSports(id, req.Name, req.MinPlayers, req.MaxPlayers)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "details": "sport not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, SportsResponse{Sport: sport, Message: "Sport updated successfully"})
}

func (s *SportsHandler) DeleteRegisterSports(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty ID"})
		return
	}

	err := s.SportsService.DeleteSports(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "details": "sport not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sport", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sport deleted successfully"})
}
