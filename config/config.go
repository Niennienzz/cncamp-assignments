package config

import (
	"cncamp_a01/constant"
	"time"
)

type Interface interface {
	Env() constant.EnvEnum
	Port() int
	Version() string
	SQLiteFileName() string
	PasswordHashSecret() string
	PasswordSaltLen() int
	TokenHMACSecret() string
	TokenExpiration() time.Duration
	RateLimit() int
	RateLimitWindow() time.Duration
}

var _ Interface = config{}

type config struct {
	env                constant.EnvEnum
	port               int
	version            string
	sqliteFileName     string
	passwordHashSecret string
	passwordSaltLen    int
	tokenHMACSecret    string
	tokenExpiration    time.Duration
	rateLimit          int
	rateLimitWindow    time.Duration
}

func (c config) Env() constant.EnvEnum {
	return c.env
}

func (c config) Port() int {
	return c.port
}

func (c config) Version() string {
	return c.version
}

func (c config) SQLiteFileName() string {
	return c.sqliteFileName
}

func (c config) PasswordHashSecret() string {
	return c.passwordHashSecret
}

func (c config) PasswordSaltLen() int {
	return c.passwordSaltLen
}

func (c config) TokenHMACSecret() string {
	return c.tokenHMACSecret
}

func (c config) TokenExpiration() time.Duration {
	return c.tokenExpiration
}

func (c config) RateLimit() int {
	return c.rateLimit
}

func (c config) RateLimitWindow() time.Duration {
	return c.rateLimitWindow
}
