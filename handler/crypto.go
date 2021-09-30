package handler

import (
	"cncamp_a01/constant"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Crypto interface {
	GetByCode() fiber.Handler
}

type cryptoHandler struct {
	*handler
}

func (h *cryptoHandler) GetByCode() fiber.Handler {
	type enum = constant.CryptoCodeEnum

	type cryptoDAO struct {
		ID         int     `json:"id" db:"id"`
		CryptoCode enum    `json:"name" db:"crypto_code"`
		Price      float64 `json:"priceUSD" db:"price"`
		UpdatedAt  string  `json:"updatedAtUTC" db:"updated_at"`
	}

	return func(c *fiber.Ctx) error {
		var (
			ctx  = c.UserContext()
			code = c.Params("crypto_code")
		)

		cryptoCode := enum(code)
		if err := cryptoCode.Valid(); err != nil {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, err)
		}

		const cryptoQuery = `SELECT * FROM cryptos WHERE crypto_code=?;`
		var (
			crypto = new(cryptoDAO)
			row    = h.db.QueryRowxContext(ctx, cryptoQuery, cryptoCode.String())
		)
		err := row.StructScan(crypto)
		if errors.Is(err, sql.ErrNoRows) {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("crypto not found"))
		}
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		c.Status(fiber.StatusOK)
		return c.JSON(crypto)
	}
}
