package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musishere/sportsApp/config"
	"github.com/musishere/sportsApp/internal/auth"
	"github.com/musishere/sportsApp/internal/helpers"
	"github.com/musishere/sportsApp/internal/oauth"
	"github.com/musishere/sportsApp/internal/services"
	"github.com/musishere/sportsApp/internal/validators"
	"github.com/musishere/sportsApp/types"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
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

type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req types.RegisterRequest

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

	user, token, err := h.userService.Register(req)

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

	// OTP flow disabled for testing - set cookie and return token
	c.SetCookie("Jwt-Token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"user": SignupUserResponse{Name: user.Name, Email: user.Email, Role: user.Role, IsActive: user.IsActive, Gender: user.Gender, Phone: user.Phone},
	})
}

func (h *UserHandler) VerifyOtp(c *gin.Context) {
	var req types.VerifyOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp are required"})
		return
	}

	ok, err := helpers.VerifyOTP(req.Phone, req.Otp)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired OTP"})
		return
	}

	user, token, err := h.userService.ActivateUserByPhone(types.ActivateUserRequest{Phone: req.Phone})
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

	var requestBody types.LoginRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user, token, err := h.userService.Login(requestBody)

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

func (h *UserHandler) SignUpOauth2Facebook(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
		return
	}

	userInfo, err := oauth.ConnectToFacebook(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return the user info
	c.JSON(http.StatusOK, gin.H{
		"message": "Facebook login successful",
		"user":    userInfo,
	})
}
