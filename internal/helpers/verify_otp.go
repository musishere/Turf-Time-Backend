package helpers

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/musishere/sportsApp/internal/cache"
)

// Verify OTP
func VerifyOTP(phone, otpInput string) (bool, error) {
	key := fmt.Sprintf("otp:%s", phone)
	storedOTP, err := cache.Rdb.Get(cache.Ctx, key).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("OTP expired or not found")
	} else if err != nil {
		return false, err
	}

	if storedOTP != otpInput {
		return false, fmt.Errorf("OTP does not match")
	}

	// OTP correct → delete
	_ = cache.Rdb.Del(cache.Ctx, key).Err()
	return true, nil
}
