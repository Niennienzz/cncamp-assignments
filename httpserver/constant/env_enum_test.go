package constant_test

import (
	"cncamp_a01/httpserver/constant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvEnum(t *testing.T) {
	dev := "DEV"
	assert.Equal(t, constant.EnvDev.String(), dev)

	invalid := constant.EnvEnum("INVALID")
	err := invalid.Valid()
	assert.Equal(t, err, constant.ErrInvalidEnvEnum)
}
