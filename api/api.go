package api

import (
	"cncamp_a01/config"
	"cncamp_a01/handler"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
)

type Interface interface {
	Run()
}

func New(cfg config.Interface) Interface {
	handlers := handler.New(cfg)

	// TODO: Make logger/limiter/auth middlewares.
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())

	// The health endpoint.
	app.All("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// User endpoints.
	user := app.Group("/user")
	{
		user.Post("/signup", handlers.User().Signup())
		user.Post("/login", handlers.User().Login())
	}

	// Crypto price endpoints.
	cp := app.Group("/crypto_price")
	{
		cp.Get("/:crypto_code", handlers.CryptoPrice().GetByCode())
	}

	return &apiImpl{
		app:     app,
		port:    cfg.Port(),
		handler: handlers,
	}
}

type apiImpl struct {
	app     *fiber.App
	port    int
	handler handler.Interface
}

func (api *apiImpl) Run() {
	defer api.handler.Close()
	err := api.app.Listen(fmt.Sprintf(":%d", api.port))
	if err != nil {
		log.Error(err)
	}
}
