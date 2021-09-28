package config

import (
	"cncamp_a01/constant"
	"flag"
	"sync"
)

const (
	DefaultEnv                = "DEV"
	DefaultPort               = 8080
	DefaultVersion            = "0.0.1"
	DefaultSQLiteFileName     = "sqlite.db"
	DefaultPasswordHashSecret = "twice_security"
	DefaultPasswordSaltLen    = 16
	DefaultTokenSecret        = "twice_security"
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
		envStr          string
		port            int
		version         string
		sqliteFileName  string
		pwdHashSecret   string
		pwdSaltLen      int
		tokenHMACSecret string
	)

	flag.StringVar(&envStr, "ENV", DefaultEnv, "environment of the server")
	flag.IntVar(&port, "PORT", DefaultPort, "port number of the server")
	flag.StringVar(&version, "VERSION", DefaultVersion, "version number of the api server")
	flag.StringVar(&sqliteFileName, "SQLITE_FILE", DefaultSQLiteFileName, "sqlite file name")
	flag.StringVar(&pwdHashSecret, "PWD_HASH_SECRET", DefaultPasswordHashSecret, "password hash secret")
	flag.IntVar(&pwdSaltLen, "PWD_SALT_LEN", DefaultPasswordSaltLen, "password salt length")
	flag.StringVar(&tokenHMACSecret, "TOKEN_SECRET", DefaultTokenSecret, "token secret")

	env := constant.EnvEnum(envStr)
	if err := env.Valid(); err != nil {
		panic(err)
	}

	single = &config{
		env:                env,
		port:               port,
		version:            version,
		sqliteFileName:     sqliteFileName,
		passwordHashSecret: pwdHashSecret,
		passwordSaltLen:    pwdSaltLen,
		tokenHMACSecret:    tokenHMACSecret,
	}
}