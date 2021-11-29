package config

import (
	"cncamp_a01/httpserver/constant"
	"time"
)

type Interface interface {
	GetEnv() constant.EnvEnum
	GetPort() int
	GetVersion() string
	GetMongoUserName() string
	GetMongoPassword() string
	GetMongoURL() string
	GetPasswordSaltLen() int
	GetPasswordHashSecret() string
	GetTokenHMACSecret() string
	GetTokenExpiration() time.Duration
	GetRateLimit() int
	GetRateLimitWindow() time.Duration
	GetFetchWindow() time.Duration
}

var _ Interface = config{}

type enum = constant.EnvEnum

type duration = time.Duration

type config struct {
	Env                enum     `required:"true" envconfig:"ENV" default:"DEV"`
	Port               int      `required:"true" envconfig:"PORT" default:"8080"`
	Version            string   `required:"true" envconfig:"VERSION" default:"0.0.1"`
	MongoUserName      string   `required:"true" envconfig:"MONGO_INITDB_ROOT_USERNAME" default:"mongo_user"`
	MongoPassword      string   `required:"true" envconfig:"MONGO_INITDB_ROOT_PASSWORD" default:"mongo_pwd"`
	MongoURL           string   `required:"true" envconfig:"MONGO_URL" default:"localhost:27017"`
	PasswordSaltLen    int      `required:"true" envconfig:"PWD_SALT_LEN" default:"16"`
	PasswordHashSecret string   `required:"true" envconfig:"PWD_HASH_SECRET" default:"twice_security"`
	TokenHMACSecret    string   `required:"true" envconfig:"TOKEN_SECRET" default:"twice_security"`
	TokenExpiration    duration `required:"true" envconfig:"TOKEN_EXPIRE_SEC" default:"62400s"`
	RateLimit          int      `required:"true" envconfig:"RATE_LIMIT" default:"30"`
	RateLimitWindow    duration `required:"true" envconfig:"RATE_LIMIT_WINDOW_SEC" default:"30s"`
	FetchWindow        duration `required:"true" envconfig:"FETCH_WINDOW_SEC" default:"30s"`
}

func (c config) GetEnv() constant.EnvEnum {
	return c.Env
}

func (c config) GetPort() int {
	return c.Port
}

func (c config) GetVersion() string {
	return c.Version
}

func (c config) GetMongoUserName() string {
	return c.MongoUserName
}

func (c config) GetMongoPassword() string {
	return c.MongoPassword
}

func (c config) GetMongoURL() string {
	return c.MongoURL
}

func (c config) GetPasswordSaltLen() int {
	return c.PasswordSaltLen
}

func (c config) GetPasswordHashSecret() string {
	return c.PasswordHashSecret
}

func (c config) GetTokenHMACSecret() string {
	return c.TokenHMACSecret
}

func (c config) GetTokenExpiration() time.Duration {
	return c.TokenExpiration
}

func (c config) GetRateLimit() int {
	return c.RateLimit
}

func (c config) GetRateLimitWindow() time.Duration {
	return c.RateLimitWindow
}

func (c config) GetFetchWindow() time.Duration {
	return c.FetchWindow
}
