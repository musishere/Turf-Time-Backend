package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a random 6-digit OTP and returns it as []int
func GenerateOTP() []int {
	rand.Seed(time.Now().UnixNano()) // Seed RNG

	otp := make([]int, 6)
	for i := 0; i < 6; i++ {
		otp[i] = rand.Intn(10) // Random digit 0-9
	}
	fmt.Println(otp)
	return otp
}
