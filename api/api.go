package api

import (
	"cncamp_a01/config"
	"cncamp_a01/handler"
	"cncamp_a01/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
)

type Interface interface {
	Run()
}

func New() Interface {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.NewUserContext())
	app.Use(middleware.NewLogger())

	handlers := handler.New()

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

	// Crypto endpoints.
	// Only authenticated user can call, and handlers follow rate-limiting rules.
	cp := app.Group("/crypto")
	cp.Use(middleware.NewAuthLimiter())
	{
		cp.Get("/:crypto_code", handlers.Crypto().GetByCode())
	}

	return &apiImpl{
		app:     app,
		port:    config.Get().Port(),
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
