package handler

import (
	"cncamp_a01/constant"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type CryptoPrice interface {
	GetByCode() fiber.Handler
}

type cryptoPriceHandler struct {
	*handler
}

func (h *cryptoPriceHandler) GetByCode() fiber.Handler {
	type enum = constant.CryptoCodeEnum

	type cryptoDAO struct {
		ID         int  `json:"id" db:"id"`
		CryptoCode enum `json:"cryptoCode" db:"crypto_code"`
		Price      int  `json:"price" db:"price"`
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
		if err := row.StructScan(crypto); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("crypto not found"))
			} else {
				return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
			}
		}

		c.Status(fiber.StatusOK)
		return c.JSON(crypto)
	}
}
