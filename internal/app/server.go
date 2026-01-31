package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/database"
	"github.com/musishere/sportsApp/internal/handlers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/services"
)

func StartServer() {
	cfg := config.LoadConfig()

	db := database.ConnectDatabase(cfg)

	if err := db.AutoMigrate(&models.User{}, &models.Location{}, &models.Sports{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	sportsRepo := repositories.NewSportsRepositry(db)

	userService := services.NewUserService(userRepo, locationRepo, cfg.JWTSecret)
	locationService := services.NewLocationService(locationRepo)
	sportsService := services.NewSportsService(sportsRepo)

	router := gin.Default()

	SetupRoutes(router, userService, locationService, sportsService, cfg.JWTSecret)

	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func SetupRoutes(
	router *gin.Engine,
	userService *services.UserService,
	locationService *services.LocationService,
	sportsService *services.SportsService,
	jwtSecret string,
) {
	api := router.Group("/api/v1")

	//User Routes
	api.POST("/signup", handlers.NewUserHandler(userService).RegisterUser)
	api.POST("/login", handlers.NewUserHandler(userService).LoginUser)
	api.GET("/get-currentUser", handlers.NewUserHandler(userService).GetCurrentUser)
	api.POST("/logout", handlers.NewUserHandler(userService).LogOutUser)

	// sports Routes
	api.POST("/sports", handlers.NewSportsHandler(sportsService).RegisterNewSports)
}
