package handlers

import (
	"github.com/gin-gonic/gin"
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

func (h *TurfHandler) RegisterTurf(c *gin.Context) {}

func (h *TurfHandler) GetRegisteredTurfs(c *gin.Context)    {}
func (h *TurfHandler) GetRegisteredTurfByID(c *gin.Context) {}
func (h *TurfHandler) UpdateRegisteredTurf(c *gin.Context)  {}
func (h *TurfHandler) DeleteRegisteredTurf(c *gin.Context)  {}
