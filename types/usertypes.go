package types

// RegisterRequest contains all parameters needed for user registration
type RegisterRequest struct {
	Name      string  `json:"name" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	Gender    string  `json:"gender" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Cnic      string  `json:"cnic"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// LoginRequest contains all parameters needed for user login
type LoginRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=6"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// ActivateUserRequest contains parameters needed to activate a user
type ActivateUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// VerifyOtpRequest contains parameters needed to verify OTP
type VerifyOtpRequest struct {
	Phone string `json:"phone" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

// RegisterResponse contains the response from user registration
type RegisterResponse struct {
	User  interface{}
	Token string
}

// LoginResponse contains the response from user login
type LoginResponse struct {
	User  interface{}
	Token string
}

// ActivateUserResponse contains the response from user activation
type ActivateUserResponse struct {
	User  interface{}
	Token string
}
