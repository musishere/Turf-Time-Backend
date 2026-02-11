package helpers

import (
	"math/rand"
	"time"
)

// GenerateOTP generates a random 6-digit OTP and returns it as []int
func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())

	// Generate 6-digit number between 100000 and 999999
	otp := rand.Intn(900000) + 100000
	return otp
}
