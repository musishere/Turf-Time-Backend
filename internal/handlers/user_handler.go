package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Name      string  `json:"name" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type CurrentUserResponse struct {
	User interface{} `json:"user"`
}

type RegisterResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type LoginRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	user, token, err := h.userService.Register(
		req.Name,
		req.Email,
		req.Password,
		req.Latitude,
		req.Longitude,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			statusCode = http.StatusConflict
		} else if err.Error() == "name, email, and password are required" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})

	}

	c.SetCookie("Jwt-Token", token, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, RegisterResponse{
		User:  user,
		Token: token,
	})
}

func (h *UserHandler) LoginUser(c *gin.Context) {

	var requestBody LoginRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user, token, err := h.userService.Login(requestBody.Email, requestBody.Password, requestBody.Latitude, requestBody.Longitude)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username or password"})
		return
	}
	c.SetCookie("Jwt-Token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, LoginResponse{
		User:  user,
		Token: token,
	})

}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {

	tokernString, err := c.Cookie("Jwt-Token")
	if err != nil {
		c.JSON(http.StatusBadRequest, "No Token recieved")
		return
	}

	secret := config.LoadConfig()

	existingUser, err := auth.VerifyJWT(tokernString, secret.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid Token")
		return
	}

	user, err := h.userService.GetByID(existingUser.UserID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, CurrentUserResponse{User: user})

}
