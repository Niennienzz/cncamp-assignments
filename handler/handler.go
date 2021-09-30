package handler

import (
	"cncamp_a01/config"
	"cncamp_a01/constant"
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type Handler interface {
	User() User
	Crypto() Crypto
	Shutdown() error
}

// New creates a handler.Handler backed by the private handler struct.
// This function can panic, since say if a database connection
// cannot be established, it does not make sense to proceed.
func New() Handler {
	cfg := config.Instance()

	db := sqlx.MustOpen("sqlite3", cfg.GetSQLiteFileName())

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

	cryptos := `
	CREATE TABLE IF NOT EXISTS cryptos (
		id INTEGER PRIMARY KEY,
		crypto_code TEXT NOT NULL UNIQUE,
		price REAL NOT NULL,
		updated_at TEXT NOT NULL
	);`
	if _, err := db.ExecContext(context.Background(), cryptos); err != nil {
		panic(err)
	}

	h := &handler{
		cfg:         cfg,
		db:          db,
		validate:    validator.New(),
		fetchClient: http.DefaultClient,
		fetchDone:   make(chan struct{}),
		fetchTicker: time.NewTicker(cfg.GetFetchWindow()),
	}

	go func() {
		h.fetchAll()
		for {
			select {
			case <-h.fetchDone:
				h.fetchTicker.Stop()
				log.Info("handler.fetchTicker stopped")
				return
			case <-h.fetchTicker.C:
				h.fetchAll()
			}
		}
	}()

	return h
}

type handler struct {
	cfg         config.Interface
	db          *sqlx.DB
	validate    *validator.Validate
	fetchClient *http.Client
	fetchDone   chan struct{}
	fetchTicker *time.Ticker
}

func (h *handler) User() User {
	return &userHandler{h}
}

func (h *handler) Crypto() Crypto {
	return &cryptoHandler{h}
}

func (h *handler) Shutdown() error {
	h.fetchDone <- struct{}{}
	return h.db.Close()
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (h *handler) sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	log.Error(err)
	if status == fiber.StatusInternalServerError && h.cfg.GetEnv() == constant.EnvProd {
		err = nil
	}
	c.Status(status)
	return c.JSON(errorResponse{ErrorMessage: err.Error()})
}
