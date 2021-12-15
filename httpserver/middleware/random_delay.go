package middleware

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	delayOnce sync.Once
	delayRand *rand.Rand
)

func RandomDelay() fiber.Handler {
	if delayRand == nil {
		delayOnce.Do(func() {
			src := rand.NewSource(time.Now().UnixNano())
			delayRand = rand.New(src)
		})
	}

	return func(c *fiber.Ctx) error {
		ms := delayRand.Int63n(2000)
		time.Sleep(time.Millisecond * time.Duration(ms))
		return c.Next()
	}
}
