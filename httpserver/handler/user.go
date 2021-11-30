package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User interface {
	Signup() fiber.Handler
	Login() fiber.Handler
}

type userHandler struct {
	*handler
}

type userDAO struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Email          string             `json:"email" json:"_id"`
	HashedPassword string             `json:"hashed_password" json:"_id"`
	Salt           string             `json:"salt" json:"_id"`
}

func (h *userHandler) Signup() fiber.Handler {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=64"`
	}

	type response struct {
		ID primitive.ObjectID `json:"id"`
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
		req.Email = strings.ToLower(req.Email)

		count, err := h.mongoDB.Collection(usersCol).CountDocuments(ctx, bson.M{"email": req.Email})
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		if count != 0 {
			return h.sendErrorResponse(c, fiber.StatusBadRequest, errors.New("email already exist"))
		}

		hashedPassword, salt := h.generatePasswordHashAndSalt(req.Password)
		user := userDAO{
			ID:             primitive.NewObjectID(),
			Email:          req.Email,
			HashedPassword: hashedPassword,
			Salt:           salt,
		}

		res, err := h.mongoDB.Collection(usersCol).InsertOne(ctx, user)
		if err != nil {
			return h.sendErrorResponse(c, fiber.StatusInternalServerError, err)
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(response{ID: res.InsertedID.(primitive.ObjectID)})
	}
}

func (h *userHandler) Login() fiber.Handler {
	type request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=64"`
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
		req.Email = strings.ToLower(req.Email)

		user := userDAO{}
		err := h.mongoDB.Collection(usersCol).FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
		if errors.Is(err, mongo.ErrNoDocuments) {
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
