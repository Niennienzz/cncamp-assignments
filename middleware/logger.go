package middleware

import (
	"cncamp_a01/config"
	"cncamp_a01/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func NewLogger() fiber.Handler {
	cfg := config.Get()
	if cfg.Env() == constant.EnvProd {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		userID, ok := UserIDStringFromContext(ctx)
		if !ok {
			userID = "none"
		}

		var (
			reqID  = uuid.NewString()
			logger = log.WithField("api_version", cfg.Version()).
				WithField("ip", c.IP()).
				WithField("path", c.Path()).
				WithField("request_id", reqID).
				WithField("user_id", userID)
		)

		logger.Info("request started")

		if err := c.Next(); err != nil {
			logger.Error(err)
			return err
		}

		logger.Info("request finished")
		return nil
	}
}