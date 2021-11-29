package constant

import (
	"database/sql/driver"
	"errors"
)

var ErrInvalidCryptoEnum = errors.New("invalid value for crypto enum")

type CryptoCodeEnum string

const (
	CryptoADA CryptoCodeEnum = "ADA"
	CryptoBNB CryptoCodeEnum = "BNB"
	CryptoBTC CryptoCodeEnum = "BTC"
	CryptoETH CryptoCodeEnum = "ETH"
)

func (e CryptoCodeEnum) String() string {
	return string(e)
}

func (e CryptoCodeEnum) Valid() error {
	switch e {
	case CryptoADA, CryptoBNB, CryptoBTC, CryptoETH:
		return nil
	default:
		return ErrInvalidCryptoEnum
	}
}

// Value implements driver.Valuer for CryptoCodeEnum.
func (e CryptoCodeEnum) Value() (driver.Value, error) {
	if err := e.Valid(); err != nil {
		return nil, err
	}
	return e.String(), nil
}

// Scan implements sql.Scanner for CryptoCodeEnum.
func (e *CryptoCodeEnum) Scan(val interface{}) error {
	var s string
	switch v := val.(type) {
	case string:
		s = v
	case []uint8:
		s = string(v)
	default:
		return ErrInvalidCryptoEnum
	}
	switch s {
	case CryptoADA.String():
		*e = CryptoADA
	case CryptoBNB.String():
		*e = CryptoBNB
	case CryptoBTC.String():
		*e = CryptoBTC
	case CryptoETH.String():
		*e = CryptoETH
	default:
		return ErrInvalidCryptoEnum
	}
	return nil
}
