package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var randomStringRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (h *userHandler) randomString(length int) string {
	runes := make([]rune, length)
	for idx := range runes {
		runes[idx] = randomStringRunes[rand.Intn(len(randomStringRunes))]
	}
	return string(runes)
}

func (h *userHandler) generatePasswordHashAndSalt(rawPassword string, existingSalt ...string) (string, string) {
	salt := h.randomString(h.cfg.GetPasswordSaltLen())
	if len(existingSalt) == 1 {
		if len(existingSalt[0]) == h.cfg.GetPasswordSaltLen() {
			salt = existingSalt[0]
		} else {
			panic("invalid salt length")
		}
	}
	var (
		pass = []byte(rawPassword + salt + h.cfg.GetPasswordHashSecret())
		hash = sha256.Sum256(pass)
	)
	return base64.URLEncoding.EncodeToString(hash[:]), salt
}

func (h *userHandler) newAccessToken(id int) (string, error) {
	var (
		sub = strconv.Itoa(id)
		now = time.Now()
	)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "api",
		ExpiresAt: now.Add(h.cfg.GetTokenExpiration()).Unix(),
		Id:        uuid.NewString(),
		IssuedAt:  now.Unix(),
		Issuer:    "api-server",
		NotBefore: now.Unix(),
		Subject:   sub,
	})
	return token.SignedString([]byte(h.cfg.GetTokenHMACSecret()))
}
