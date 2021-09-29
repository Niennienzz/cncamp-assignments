package middleware

import (
	"cncamp_a01/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

func NewUserContext() fiber.Handler {
	cfg := config.Instance()

	return func(c *fiber.Ctx) error {
		// Set version as response header.
		c.Set("Api-Version", cfg.GetVersion())

		var (
			authorizationBytes = c.Request().Header.Peek("Authorization")
			authorization      = string(authorizationBytes)
		)

		// No user ID is set if we cannot decode the token.
		if len(authorization) == 0 {
			return c.Next()
		}

		splits := strings.Split(authorization, " ")
		if len(splits) != 2 {
			return c.Next()
		}

		bearer, tokenString := splits[0], splits[1]
		if bearer != "Bearer" {
			return c.Next()
		}

		// Once we can decode the token, the token must be valid.
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.GetTokenHMACSecret()), nil
		})
		if err != nil {
			log.Error(err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !(ok && token.Valid) {
			log.Error("invalid token")
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Set user ID into user context.
		ctx := c.UserContext()
		ctx = withUserIDString(ctx, claims.Subject)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
