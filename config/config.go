package config

type Interface interface {
	Env() Env
	Port() int
	Version() string
	SQLiteFileName() string
	PasswordHashSecret() string
	PasswordSaltLen() int
	TokenHMACSecret() string
}

var _ Interface = config{}

type config struct {
	env                Env
	port               int
	version            string
	sqliteFileName     string
	passwordHashSecret string
	passwordSaltLen    int
	tokenHMACSecret    string
}

func (c config) Env() Env {
	return c.env
}

func (c config) Port() int {
	return c.port
}

func (c config) Version() string {
	return c.version
}

func (c config) SQLiteFileName() string {
	return c.sqliteFileName
}

func (c config) PasswordHashSecret() string {
	return c.passwordHashSecret
}

func (c config) PasswordSaltLen() int {
	return c.passwordSaltLen
}

func (c config) TokenHMACSecret() string {
	return c.tokenHMACSecret
}
