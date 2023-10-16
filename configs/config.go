package configs

import (
	"os"
	"strconv"
	"time"
)

type config struct {
	DB             string
	DBPort         string
	DBHost         string
	DBUser         string
	DBPassword     string
	APIAddr        string
	APIPort        string
	Secret         string
	Timeout        time.Duration
	SessionExpires time.Duration
}

var cfg *config

func Load() error {
	cfg = new(config)

	cfg.DB = os.Getenv("DB")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.APIAddr = os.Getenv("API_ADDR")
	cfg.APIPort = os.Getenv("API_PORT")
	cfg.Secret = os.Getenv("SECRET")

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		// default timeout
		cfg.Timeout = time.Duration(10)
	} else {
		cfg.Timeout = time.Duration(timeout)
	}

	sessionExpires, err := strconv.Atoi(os.Getenv("SESSION_EXPIRES"))
	if err != nil {
		// default session expires time
		cfg.SessionExpires = time.Duration(1800)
	} else {
		cfg.SessionExpires = time.Duration(sessionExpires)
	}

	return nil
}

func GetConfigs() *config {
	return cfg
}
