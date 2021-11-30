package handler

import (
	"cncamp_a01/httpserver/constant"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Crypto interface {
	GetByCode() fiber.Handler
}

type cryptoHandler struct {
	*handler
}

func (h *cryptoHandler) GetByCode() fiber.Handler {
	type cryptoCodeEnum = constant.CryptoCodeEnum

	type cryptoDAO struct {
		ID         primitive.ObjectID `json:"id" bson:"_id"`
		CryptoCode cryptoCodeEnum     `json:"crypto_code" bson:"crypto_code"`
		Price      float64            `json:"price" bson:"price"`
		UpdatedAt  time.Time          `json:"updatedAt" bson:"updated_at"`
	}

	return func(c *fiber.Ctx) error {
		var (
			ctx  = c.UserContext()
			code = c.Params("crypto_code")
		)

		cryptoCode := cryptoCodeEnum(code)
		if err := cryptoCode.Valid(); err != nil {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, err)
		}

		crypto := cryptoDAO{}
		err := h.mongoDB.Collection(cryptosCol).FindOne(ctx, bson.M{"crypto_code": cryptoCode}).Decode(&crypto)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("crypto not found"))
		}
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		c.Status(fiber.StatusOK)
		return c.JSON(crypto)
	}
}
