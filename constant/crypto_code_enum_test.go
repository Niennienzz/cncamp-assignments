package constant_test

import (
	"cncamp_a01/constant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoCodeEnum(t *testing.T) {
	ada := "ADA"
	assert.Equal(t, constant.CryptoADA.String(), ada)

	invalid := constant.CryptoCodeEnum("INVALID")
	err := invalid.Valid()
	assert.Equal(t, err, constant.ErrInvalidCryptoEnum)

	val, err := constant.CryptoADA.Value()
	assert.Equal(t, err, nil)
	assert.Equal(t, constant.CryptoADA.String(), val.(string))

	enum := new(constant.CryptoCodeEnum)
	err = enum.Scan("BTC")
	assert.Equal(t, err, nil)
	assert.Equal(t, *enum, constant.CryptoBTC)

	enum = new(constant.CryptoCodeEnum)
	err = enum.Scan([]uint8{'E', 'T', 'H'})
	assert.Equal(t, err, nil)
	assert.Equal(t, *enum, constant.CryptoETH)
}
