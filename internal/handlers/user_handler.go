package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/services"
	"github.com/musishere/sportsApp/internal/validators"
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
	Cnic      string  `json:"cnic"` // optional at signup; can be set later
	Phone     string  `json:"phone" required:"true"`
	Gender    string  `json:"gender" required:"true"`
}

type CurrentUserResponse struct {
	User interface{} `json:"user"`
}

// SignupUserResponse contains only the user fields returned on signup (token is in cookie).
type SignupUserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}

type RegisterResponse struct {
	User SignupUserResponse `json:"user"`
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
		return
	}

	if req.Cnic == "" {
		req.Cnic = " " // store space until set later
	}

	user, _, err := h.userService.Register(
		req.Name,
		req.Email,
		req.Password,
		req.Gender,
		req.Phone,
		req.Cnic,
		req.Latitude,
		req.Longitude,
	)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already registered" {
			statusCode = http.StatusConflict
		} else {
			var ve *validators.ValidationError
			if errors.As(err, &ve) {
				statusCode = http.StatusBadRequest
			}
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// No token until OTP verified; do not set cookie
	c.JSON(http.StatusCreated, gin.H{
		"user":    SignupUserResponse{Name: user.Name, Email: user.Email, Role: user.Role, IsActive: user.IsActive, Gender: user.Gender, Phone: user.Phone},
		"message": "OTP sent to your phone. Please verify to activate your account.",
	})
}

type VerifyOtpRequest struct {
	Phone string `json:"phone" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

func (h *UserHandler) VerifyOtp(c *gin.Context) {
	var req VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp are required"})
		return
	}

	ok, err := helpers.VerifyOTP(req.Phone, req.Otp)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired OTP"})
		return
	}

	user, token, err := h.userService.ActivateUserByPhone(req.Phone)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "user not found for this phone" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("Jwt-Token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"user":    SignupUserResponse{Name: user.Name, Email: user.Email, Role: user.Role, IsActive: user.IsActive, Gender: user.Gender, Phone: user.Phone},
		"message": "Phone verified. Account is now active.",
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

func (h *UserHandler) LogOutUser(c *gin.Context) {
	c.SetCookie("Jwt-Token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out",
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
