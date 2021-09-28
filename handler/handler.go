package handler

import (
	"cncamp_a01/config"
	"cncamp_a01/constant"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Interface interface {
	User() User
	CryptoPrice() CryptoPrice
	Close()
}

// New creates a handler.Interface backed by the private handler struct.
// This function can panic, since say if a database connection
// cannot be established, it does not make sense to proceed.
func New(cfg config.Interface) Interface {
	db := sqlx.MustOpen("sqlite3", cfg.SQLiteFileName())

	users := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		hashed_password TEXT NOT NULL,
		salt TEXT NOT NULL
	);`
	if _, err := db.ExecContext(context.Background(), users); err != nil {
		panic(err)
	}

	// TODO: Add an HTTP client to periodically fetch new prices.
	cryptos := `
	CREATE TABLE IF NOT EXISTS cryptos (
		id INTEGER PRIMARY KEY,
		crypto_code TEXT NOT NULL UNIQUE,
		price INT NOT NULL
	);`
	if _, err := db.ExecContext(context.Background(), cryptos); err != nil {
		panic(err)
	}

	return &handler{
		cfg:      cfg,
		db:       db,
		validate: validator.New(),
	}
}

type handler struct {
	cfg      config.Interface
	db       *sqlx.DB
	validate *validator.Validate
}

func (h *handler) User() User {
	return &userHandler{h}
}

func (h *handler) CryptoPrice() CryptoPrice {
	return &cryptoPriceHandler{h}
}

func (h *handler) Close() {
	err := h.db.Close()
	if err != nil {
		log.Error(err)
	}
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (h *handler) sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	log.Error(err)
	if status == fiber.StatusInternalServerError && h.cfg.Env() == constant.EnvProd {
		err = nil
	}
	c.Status(status)
	return c.JSON(errorResponse{ErrorMessage: err.Error()})
}
