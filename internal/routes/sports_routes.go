package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/internal/handlers"
	"github.com/musishere/sportsApp/internal/services"
)

func SetupSportsRoutes(api *gin.RouterGroup, sportsService *services.SportsService) {
	api.POST("/sports", handlers.NewSportsHandler(sportsService).RegisterNewSports)
	api.GET("/sports", handlers.NewSportsHandler(sportsService).GetAllRegisteredSports)
	api.GET("/sports/:id", handlers.NewSportsHandler(sportsService).GetRegisteredSportsByID)
	api.PATCH("/sports/:id", handlers.NewSportsHandler(sportsService).UpdateRegisteredSports)
	api.DELETE("/sports/:id", handlers.NewSportsHandler(sportsService).DeleteRegisteredSport)
}
