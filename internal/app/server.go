package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/database"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/routes"
	"github.com/musishere/sportsApp/internal/services"
	"golang.org/x/time/rate"
)

func StartServer() {
	cfg := config.LoadConfig()
	db := database.ConnectDatabase(cfg)

	if err := db.AutoMigrate(&models.User{}, &models.Location{}, &models.Sports{}, &models.Turf{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	sportsRepo := repositories.NewSportsRepository(db)
	turfRepo := repositories.NewTurfRepository(db)

	imageUploader, err := helpers.NewImageUploader()
	if err != nil {
		log.Fatal("Cloudinary init failed:", err)
	}

	userService := services.NewUserService(userRepo, locationRepo, cfg.JWTSecret)
	sportsService := services.NewSportsService(sportsRepo, imageUploader)
	turfService := services.NewTurfService(turfRepo, imageUploader)

	// Global rate limiter
	limiter := rate.NewLimiter(rate.Limit(100), 50)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	})

	api := router.Group("/api/v1")
	routes.SetupUserRoutes(api, userService)
	routes.SetupSportsRoutes(api, sportsService)
	routes.SetupTurfRoutes(api, turfService)

	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
