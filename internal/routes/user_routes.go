package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/internal/handlers"
	"github.com/musishere/sportsApp/internal/services"
)

func SetupUserRoutes(api *gin.RouterGroup, userService *services.UserService) {
	api.POST("/signup", handlers.NewUserHandler(userService).RegisterUser)
	// api.POST("/verify-phone-otp", handlers.NewUserHandler(userService).VerifyOtp)
	api.POST("/verify-email-otp", handlers.NewUserHandler(userService).VerifyEmailOtp)
	api.POST("/login", handlers.NewUserHandler(userService).LoginUser)
	api.GET("/get-currentUser", handlers.NewUserHandler(userService).GetCurrentUser)
	api.POST("/logout", handlers.NewUserHandler(userService).LogOutUser)
	api.POST("/oauth-facebook", handlers.NewUserHandler(userService).SignUpOauth2Facebook)
}
