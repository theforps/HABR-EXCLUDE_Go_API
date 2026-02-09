package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func RateLimiter(delay time.Duration) fiber.Handler {
	limiter := make(chan struct{}, 1)
	limiter <- struct{}{}

	return func(c fiber.Ctx) error {
		<-limiter

		go func() {
			time.Sleep(delay)
			limiter <- struct{}{}
		}()

		return c.Next()
	}
}
