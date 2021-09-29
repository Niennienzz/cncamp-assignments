package config

import (
	"cncamp_a01/constant"
	"flag"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultEnv                = "DEV"
	DefaultPort               = 8080
	DefaultVersion            = "0.0.1"
	DefaultSQLiteFileName     = "sqlite.db"
	DefaultPasswordSaltLen    = 16
	DefaultPasswordHashSecret = "twice_security" // Nobody uses this in prod.
	DefaultTokenSecret        = "twice_security" // Nobody uses this in prod.
	DefaultTokenExpireSec     = 62400            // 1 day
	DefaultRateLimit          = 30               // 30 requests per 30 seconds
	DefaultRateLimitWindowSec = 30               // 30 requests per 30 seconds
)

var (
	once   sync.Once
	single *config
)

// Get returns parsed configurations.
// It only parses once, and returns the singleton afterwards.
func Get() Interface {
	if single == nil {
		once.Do(parse)
	}
	return single
}

func parse() {
	var (
		envStr             string
		port               int
		version            string
		sqliteFileName     string
		pwdHashSecret      string
		pwdSaltLen         int
		tokenHMACSecret    string
		tokenExpireSec     int
		rateLimit          int
		rateLimitWindowSec int
	)

	flag.StringVar(&envStr, "ENV", DefaultEnv, "environment of the api server")
	flag.IntVar(&port, "PORT", DefaultPort, "port number of the api server")
	flag.StringVar(&version, "VERSION", DefaultVersion, "version number of the api server")
	flag.StringVar(&sqliteFileName, "SQLITE_FILE", DefaultSQLiteFileName, "sqlite file name")
	flag.IntVar(&pwdSaltLen, "PWD_SALT_LEN", DefaultPasswordSaltLen, "password salt length")
	flag.StringVar(&pwdHashSecret, "PWD_HASH_SECRET", DefaultPasswordHashSecret, "password hash secret")
	flag.StringVar(&tokenHMACSecret, "TOKEN_SECRET", DefaultTokenSecret, "token secret")
	flag.IntVar(&tokenExpireSec, "TOKEN_EXPIRE_SEC", DefaultTokenExpireSec, "token expiration in seconds")
	flag.IntVar(&rateLimit, "RATE_LIMIT", DefaultRateLimit, "rate limit: how many requests per window")
	flag.IntVar(&rateLimitWindowSec, "RATE_LIMIT_WINDOW_SEC", DefaultRateLimitWindowSec, "rate limit window in seconds")

	env := constant.EnvEnum(envStr)
	if err := env.Valid(); err != nil {
		panic(err)
	}

	tokenExpiration := time.Duration(tokenExpireSec) * time.Second

	rateLimitWindow := time.Duration(rateLimitWindowSec) * time.Second

	if pwdHashSecret == DefaultPasswordHashSecret || tokenHMACSecret == DefaultTokenSecret {
		if env == constant.EnvDev {
			log.Warn("should override secrets, ok for dev")
		} else {
			log.Fatal("invalid secrets")
		}
	}

	single = &config{
		env:                env,
		port:               port,
		version:            version,
		sqliteFileName:     sqliteFileName,
		passwordSaltLen:    pwdSaltLen,
		passwordHashSecret: pwdHashSecret,
		tokenHMACSecret:    tokenHMACSecret,
		tokenExpiration:    tokenExpiration,
		rateLimit:          rateLimit,
		rateLimitWindow:    rateLimitWindow,
	}
}
