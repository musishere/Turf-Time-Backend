package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/internal/handlers"
	"github.com/musishere/sportsApp/internal/services"
)

func SetupTurfRoutes(api *gin.RouterGroup, turfService *services.TurfService) {
	api.POST("/turfs", handlers.NewTurfHandler(turfService).RegisterTurf)
}
