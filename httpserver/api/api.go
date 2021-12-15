package api

import (
	"cncamp_a01/httpserver/config"
	"cncamp_a01/httpserver/handler"
	"cncamp_a01/httpserver/metrics"
	"cncamp_a01/httpserver/middleware"
	"fmt"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
)

type Interface interface {
	Run()
	Shutdown()
}

func New() Interface {
	// Register Prometheus metrics.
	metrics.Register()

	// Initialize handlers.
	handlers := handler.New()

	// Register Fiber middlewares.
	fiberApp := fiber.New(fiber.Config{ReadTimeout: time.Second * 5})
	fiberApp.Use(recover.New())
	fiberApp.Use(cors.New())
	fiberApp.Use(middleware.Header())
	fiberApp.Use(middleware.UserContext())
	fiberApp.Use(middleware.Logger())

	// The Prometheus metrics endpoint.
	fiberApp.All("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Add random delay in the following endpoints.
	delayed := fiberApp.Use(middleware.RandomDelay())

	// The health check endpoint.
	delayed.All("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// User endpoints.
	user := delayed.Group("/user")
	{
		user.Post("/signup", handlers.User().Signup())
		user.Post("/login", handlers.User().Login())
	}

	// Crypto endpoints.
	// Only authenticated user can call, and handlers follow rate-limiting rules.
	cp := delayed.Group("/crypto")
	cp.Use(middleware.AuthLimiter())
	{
		cp.Get("/:crypto_code", handlers.Crypto().GetByCode())
	}

	return &apiImpl{
		fiberApp: fiberApp,
		port:     config.Instance().GetPort(),
		handler:  handlers,
	}
}

type apiImpl struct {
	fiberApp *fiber.App
	port     int
	handler  handler.Handler
}

func (api *apiImpl) Run() {
	err := api.fiberApp.Listen(fmt.Sprintf(":%d", api.port))
	if err != nil {
		api.Shutdown()
		log.Fatalln(err)
	}
}

func (api *apiImpl) Shutdown() {
	log.Info("api server shutting down...")

	log.Info("waiting for active connections to close...")

	if err := api.fiberApp.Shutdown(); err != nil {
		log.Error(err)
	}
	log.Info("api fiber app shut down")

	if err := api.handler.Shutdown(); err != nil {
		log.Error(err)
	}
	log.Info("api handler shut down")

	log.Info("api server gracefully shut down")
}
