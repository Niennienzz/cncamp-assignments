package handler

import (
	"github.com/gofiber/fiber/v2"
)

type CryptoPrice interface {
	GetByCode() fiber.Handler
}

type cryptoPriceHandler struct {
	*handler
}

func (h *cryptoPriceHandler) GetByCode() fiber.Handler {
	return nil
}
