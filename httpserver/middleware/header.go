package middleware

import (
	"cncamp_a01/httpserver/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Header() fiber.Handler {
	cfg := config.Instance()

	return func(c *fiber.Ctx) error {
		// Set version as a response header.
		c.Set("Version", cfg.GetVersion())

		// Add request headers into response headers.
		c.Request().Header.VisitAll(func(key, val []byte) {
			newKey := fmt.Sprintf("Req-%s", key)
			c.Response().Header.AddBytesV(newKey, val)
		})

		return c.Next()
	}
}
