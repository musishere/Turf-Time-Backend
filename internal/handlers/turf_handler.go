package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musishere/sportsApp/internal/services"
)

type TurfHandler struct {
	turfService services.TurfService
}

func NewTurfHandler(turfService *services.TurfService) *TurfHandler {
	return &TurfHandler{
		turfService: *turfService,
	}
}

func (h *TurfHandler) RegisterTurf(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20) // 32MB

	name := c.PostForm("name")
	startTime, _ := strconv.Atoi(c.PostForm("startTime"))
	endTime, _ := strconv.Atoi(c.PostForm("endTime"))
	status := c.PostForm("status")
	noOfFields, _ := strconv.Atoi(c.PostForm("noOfFields"))
	address := c.PostForm("address")
	ownerID, err := uuid.Parse(c.PostForm("ownerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ownerId is required and must be a valid UUID"})
		return
	}

	readFile := func(field string) ([]byte, string, error) {
		fileHeader, err := c.FormFile(field)
		if err != nil {
			return nil, "", err
		}
		file, err := fileHeader.Open()
		if err != nil {
			return nil, "", err
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, "", err
		}
		return data, fileHeader.Filename, nil
	}

	img1, fn1, err := readFile("image1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image1 is required (form-data key 'image1', type: File)", "details": err.Error()})
		return
	}
	img2, fn2, err := readFile("image2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image2 is required (form-data key 'image2', type: File)", "details": err.Error()})
		return
	}
	img3, fn3, err := readFile("image3")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image3 is required (form-data key 'image3', type: File)", "details": err.Error()})
		return
	}

	turf, err := h.turfService.CreateTurf(name, startTime, endTime, status, noOfFields, address, ownerID, img1, img2, img3, fn1, fn2, fn3)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, turf)
}

func (h *TurfHandler) GetRegisteredTurfs(c *gin.Context)    {}
func (h *TurfHandler) GetRegisteredTurfByID(c *gin.Context) {}
func (h *TurfHandler) UpdateRegisteredTurf(c *gin.Context)  {}
func (h *TurfHandler) DeleteRegisteredTurf(c *gin.Context)  {}
