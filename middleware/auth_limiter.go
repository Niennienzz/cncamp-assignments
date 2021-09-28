package middleware

import (
	"cncamp_a01/config"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type authLimiter struct {
	limit       int
	window      time.Duration
	tokenSecret string
	mutex       sync.Mutex
	userTimes   map[int][]time.Time
}

func (limiter *authLimiter) shouldLimit(userID int) bool {
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

// NewAuthLimiter rate-limits requests based on the user token.
// The same user cannot make more requests than the limit during a given window.
func NewAuthLimiter(cfg config.Interface) fiber.Handler {
	if limiter == nil {
		once.Do(func() {
			limiter = &authLimiter{
				limit:       cfg.RateLimit(),
				window:      cfg.RateLimitWindow(),
				tokenSecret: cfg.TokenHMACSecret(),
				userTimes:   make(map[int][]time.Time),
				mutex:       sync.Mutex{},
			}
		})
	}

	return func(c *fiber.Ctx) error {
		var (
			authorizationBytes = c.Request().Header.Peek("Authorization")
			authorization      = string(authorizationBytes)
		)

		if len(authorization) == 0 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		splits := strings.Split(authorization, " ")
		if len(splits) != 2 {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		bearer, tokenString := splits[0], splits[1]
		if bearer != "Bearer" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(limiter.tokenSecret), nil
		})

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !(ok && token.Valid) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.SendString("invalid user id in token")
		}

		if ok := limiter.shouldLimit(userID); !ok {
			return c.SendStatus(fiber.StatusTooManyRequests)
		}

		return c.Next()
	}
}
