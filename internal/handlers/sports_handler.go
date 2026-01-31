package handlers

import (
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/services"
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

type SportsResponse struct {
	Sport   *models.Sports
	Message string
}

func NewSportsHandler(sportsService *services.SportsService) *SportsHandler {
	return &SportsHandler{
		SportsService: *sportsService,
	}
}

// HANDLER - Read file here
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

	// Read bytes (safer for passing to service)
	fileBytes := make([]byte, fileHeader.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		c.JSON(400, gin.H{"error": "cannot read file data"})
		return
	}

	// Pass bytes instead of FileHeader
	sports, err := s.SportsService.CreateNewSport(name, minPlayers, maxPlayers, fileBytes, fileHeader.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, sports)
}
func (s *SportsHandler) GetAllRegisteredSports(c *gin.Context)  {}
func (s *SportsHandler) GetRegisteredSportsByID(c *gin.Context) {}
func (s *SportsHandler) UpdateRegisteredSports(c *gin.Context)  {}
func (s *SportsHandler) DeleteRegisterSports(c *gin.Context)    {}
