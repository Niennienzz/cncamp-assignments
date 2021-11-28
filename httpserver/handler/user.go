package handler

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type User interface {
	Signup() fiber.Handler
	Login() fiber.Handler
}

type userHandler struct {
	*handler
}

func (h *userHandler) Signup() fiber.Handler {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=64"`
	}

	type response struct {
		Created bool `json:"created"`
	}

	return func(c *fiber.Ctx) error {
		var (
			ctx = c.UserContext()
			req = new(request)
		)
		if err := c.BodyParser(req); err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		if err := h.validate.Struct(req); err != nil {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, err)
		}

		tx, err := h.db.Begin()
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		defer tx.Rollback()

		const userQuery = `SELECT COUNT(*) FROM users WHERE email=?;`
		var (
			count int
			row   = tx.QueryRowContext(ctx, userQuery, req.Email)
		)
		if err := row.Scan(&count); err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		if count != 0 {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("email already exist"))
		}

		hashedPassword, salt := h.generatePasswordHashAndSalt(req.Password)

		const userCreation = `INSERT INTO users (email, hashed_password, salt) VALUES (?, ?, ?);`
		if _, err := tx.ExecContext(ctx, userCreation, req.Email, hashedPassword, salt); err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		if err := tx.Commit(); err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(response{Created: true})
	}
}

func (h *userHandler) Login() fiber.Handler {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=64"`
	}

	type userDAO struct {
		ID             int    `db:"id"`
		Email          string `db:"email"`
		HashedPassword string `db:"hashed_password"`
		Salt           string `db:"salt"`
	}

	type response struct {
		AccessToken string `json:"accessToken"`
	}

	return func(c *fiber.Ctx) error {
		var (
			ctx = c.UserContext()
			req = new(request)
		)
		if err := c.BodyParser(req); err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		if err := h.validate.Struct(req); err != nil {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, err)
		}

		const userQuery = `SELECT * FROM users WHERE email=?;`
		var (
			user = new(userDAO)
			row  = h.db.QueryRowxContext(ctx, userQuery, req.Email)
		)
		err := row.StructScan(user)
		if errors.Is(err, sql.ErrNoRows) {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("user not found"))
		}
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		hashedPassword, _ := h.generatePasswordHashAndSalt(req.Password, user.Salt)
		if hashedPassword != user.HashedPassword {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("invalid email or password"))
		}

		accessToken, err := h.newAccessToken(user.ID)
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		c.Status(fiber.StatusOK)
		return c.JSON(response{AccessToken: accessToken})
	}
}
