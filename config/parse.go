package config

import (
	"cncamp_a01/constant"
	"flag"
	"sync"
	"time"
)

const (
	DefaultEnv                = "DEV"
	DefaultPort               = 8080
	DefaultVersion            = "0.0.1"
	DefaultSQLiteFileName     = "sqlite.db"
	DefaultPasswordHashSecret = "twice_security"
	DefaultPasswordSaltLen    = 16
	DefaultTokenSecret        = "twice_security"
	DefaultTokenExpireSec     = 62400 // 1 day
	DefaultRateLimit          = 30    // 30 requests per 30 seconds
	DefaultRateLimitWindowSec = 30    // 30 requests per 30 seconds
)

var (
	once   sync.Once
	single *config
)

// Parse returns parsed configurations.
// It only parses once, and returns the singleton afterwards.
func Parse() Interface {
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
	flag.StringVar(&pwdHashSecret, "PWD_HASH_SECRET", DefaultPasswordHashSecret, "password hash secret")
	flag.IntVar(&pwdSaltLen, "PWD_SALT_LEN", DefaultPasswordSaltLen, "password salt length")
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

	single = &config{
		env:                env,
		port:               port,
		version:            version,
		sqliteFileName:     sqliteFileName,
		passwordHashSecret: pwdHashSecret,
		passwordSaltLen:    pwdSaltLen,
		tokenHMACSecret:    tokenHMACSecret,
		tokenExpiration:    tokenExpiration,
		rateLimit:          rateLimit,
		rateLimitWindow:    rateLimitWindow,
	}
}
