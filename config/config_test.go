package config_test

import (
	"cncamp_a01/config"
	"cncamp_a01/constant"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var (
		testEnv                = constant.EnvDev
		testPort               = 3000
		testVersion            = "1.0.0"
		testSQLiteFile         = "sqlite_test.db"
		testPasswordSaltLen    = 16
		testPasswordHashSecret = "more_security"
		testTokenSecret        = "more_security"
		testRateLimit          = 30
		testDuration           = time.Second * 100
		testDurationSec        = "100s"
	)

	os.Setenv("ENV", testEnv.String())
	os.Setenv("PORT", fmt.Sprintf("%d", testPort))
	os.Setenv("VERSION", testVersion)
	os.Setenv("SQLITE_FILE", testSQLiteFile)
	os.Setenv("PWD_SALT_LEN", fmt.Sprintf("%d", testPasswordSaltLen))
	os.Setenv("PWD_HASH_SECRET", testPasswordHashSecret)
	os.Setenv("TOKEN_SECRET", testTokenSecret)
	os.Setenv("TOKEN_EXPIRE_SEC", testDurationSec)
	os.Setenv("RATE_LIMIT", fmt.Sprintf("%d", testRateLimit))
	os.Setenv("RATE_LIMIT_WINDOW_SEC", testDurationSec)
	os.Setenv("FETCH_WINDOW_SEC", testDurationSec)

	cfg := config.Instance()

	assert.Equal(t, cfg.GetEnv(), testEnv)
	assert.Equal(t, cfg.GetPort(), testPort)
	assert.Equal(t, cfg.GetVersion(), testVersion)
	assert.Equal(t, cfg.GetSQLiteFileName(), testSQLiteFile)
	assert.Equal(t, cfg.GetPasswordSaltLen(), testPasswordSaltLen)
	assert.Equal(t, cfg.GetPasswordHashSecret(), testPasswordHashSecret)
	assert.Equal(t, cfg.GetRateLimit(), testRateLimit)
	assert.Equal(t, cfg.GetTokenExpiration(), testDuration)
	assert.Equal(t, cfg.GetRateLimitWindow(), testDuration)
	assert.Equal(t, cfg.GetFetchWindow(), testDuration)
}
