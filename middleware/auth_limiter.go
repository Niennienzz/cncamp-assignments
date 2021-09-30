package middleware

import (
	"cncamp_a01/config"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authLimiter struct {
	limit     int
	window    time.Duration
	mutex     sync.Mutex
	userTimes map[string][]time.Time
}

func (limiter *authLimiter) shouldLimit(userID string) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	var (
		now      = time.Now()
		times    = limiter.userTimes[userID]
		past     = now.Add(-limiter.window)
		newTimes []time.Time
	)

	for _, v := range times {
		if v.After(past) {
			newTimes = append(newTimes, v)
		}
	}

	if len(newTimes) >= limiter.limit {
		limiter.userTimes[userID] = newTimes
		return false
	} else {
		newTimes = append(newTimes, now)
		limiter.userTimes[userID] = newTimes
		return true
	}
}

var (
	once    sync.Once
	limiter *authLimiter
)

// AuthLimiter rate-limits requests based on the user token.
// The same user cannot make more requests than the limit during a given window.
func AuthLimiter() fiber.Handler {
	cfg := config.Instance()

	if limiter == nil {
		once.Do(func() {
			limiter = &authLimiter{
				limit:     cfg.GetRateLimit(),
				window:    cfg.GetRateLimitWindow(),
				userTimes: make(map[string][]time.Time),
				mutex:     sync.Mutex{},
			}
		})
	}

	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		userID, ok := UserIDStringFromContext(ctx)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if ok := limiter.shouldLimit(userID); !ok {
			return c.SendStatus(fiber.StatusTooManyRequests)
		}

		return c.Next()
	}
}
