package constant

import (
	"errors"
)

var ErrInvalidEnvEnum = errors.New("invalid value for env enum")

type EnvEnum string

const (
	EnvDev     EnvEnum = "DEV"
	EnvTest    EnvEnum = "TEST"
	EnvStaging EnvEnum = "STAGING"
	EnvProd    EnvEnum = "PROD"
)

func (e EnvEnum) String() string {
	return string(e)
}

func (e EnvEnum) Valid() error {
	switch e {
	case EnvDev, EnvTest, EnvStaging, EnvProd:
		return nil
	default:
		return ErrInvalidEnvEnum
	}
}
