package handler

import (
	"cncamp_a01/httpserver/config"
	"cncamp_a01/httpserver/constant"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	cryptosCol = "cryptos"
	usersCol   = "users"
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

	dbURI := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		cfg.GetMongoUserName(),
		cfg.GetMongoPassword(),
		cfg.GetMongoURL(),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}

	mongoDB := mongoClient.Database("cncamp")

	h := &handler{
		cfg:         cfg,
		mongoClient: mongoClient,
		mongoDB:     mongoDB,
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
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return h.mongoClient.Disconnect(ctx)
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
