package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/database"
	"github.com/musishere/sportsApp/internal/handlers"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/models"
	"github.com/musishere/sportsApp/internal/repositories"
	"github.com/musishere/sportsApp/internal/services"
	"golang.org/x/time/rate"
)

func StartServer() {
	cfg := config.LoadConfig()

	db := database.ConnectDatabase(cfg)

	if err := db.AutoMigrate(&models.User{}, &models.Location{}, &models.Sports{}); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	userRepo := repositories.NewUserRepository(db)
	locationRepo := repositories.NewLocationRepository(db)
	sportsRepo := repositories.NewSportsRepository(db)

	imageUploader, err := helpers.NewImageUploader()
	if err != nil {
		log.Fatal("Cloudinary init failed:", err)
	}

	userService := services.NewUserService(userRepo, locationRepo, cfg.JWTSecret)
	locationService := services.NewLocationService(locationRepo)
	sportsService := services.NewSportsService(sportsRepo, imageUploader)

	// Global rate limit: 100 req/sec, burst 50 (per server)
	limiter := rate.NewLimiter(rate.Limit(100), 50)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	})

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

	// user routes
	api.POST("/signup", handlers.NewUserHandler(userService).RegisterUser)
	api.POST("/login", handlers.NewUserHandler(userService).LoginUser)
	api.GET("/get-currentUser", handlers.NewUserHandler(userService).GetCurrentUser)
	api.POST("/logout", handlers.NewUserHandler(userService).LogOutUser)

	// sports routes
	api.POST("/sports", handlers.NewSportsHandler(sportsService).RegisterNewSports)
	api.GET("/sports", handlers.NewSportsHandler(sportsService).GetAllRegisteredSports)
	api.GET("/sports/:id", handlers.NewSportsHandler(sportsService).GetRegisteredSportsByID)
	api.PATCH("/sports/:id", handlers.NewSportsHandler(sportsService).UpdateRegisteredSports)
	api.DELETE("/sports/:id", handlers.NewSportsHandler(sportsService).DeleteRegisteredSport)

	// turf routes
}
