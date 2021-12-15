package config

import (
	"cncamp_a01/httpserver/constant"
	"sync"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultPasswordHashSecret = "twice_security" // Nobody uses this in prod.
	DefaultTokenSecret        = "twice_security" // Nobody uses this in prod.
)

var (
	once   sync.Once
	single *config
)

// Instance returns parsed environment configurations.
// The function parses environment variables only once, and returns the singleton afterwards.
func Instance() Interface {
	if single == nil {
		once.Do(parse)
	}
	return single
}

func parse() {
	single = new(config)

	if err := envconfig.Process("", single); err != nil {
		log.Fatal(err)
	}

	if err := single.GetEnv().Valid(); err != nil {
		log.Fatal(err)
	}

	if single.GetPasswordHashSecret() == DefaultPasswordHashSecret ||
		single.GetTokenHMACSecret() == DefaultTokenSecret {
		if single.GetEnv() == constant.EnvDev {
			log.Warn("should override secrets, ok for dev")
		} else {
			log.Fatal("invalid secrets")
		}
	}
}
