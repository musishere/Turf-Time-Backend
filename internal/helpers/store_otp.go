package helpers

import (
	"fmt"
	"time"

	"github.com/musishere/sportsApp/internal/cache"
)

func StoreOTP(phone string, otp string) error {
	key := fmt.Sprintf("otp:%s", phone)
	return cache.Rdb.Set(cache.Ctx, key, otp, 5*time.Minute).Err()
}
