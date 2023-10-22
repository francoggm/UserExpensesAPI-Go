package configs

import (
	"time"

	"github.com/spf13/viper"
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
	LogRemoveDays  int
	CookieDomain   string
}

var cfg *config

func init() {
	cfg = new(config)

	viper.SetDefault("TIMEOUT", "10")
	viper.SetDefault("SESSION_EXPIRES", "1800")
	viper.SetDefault("API_PORT", "8080")
}

func Load() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	cfg.DB = viper.GetString("POSTGRES_DB")
	cfg.DBPort = viper.GetString("POSTGRES_PORT")
	cfg.DBHost = viper.GetString("POSTGRES_HOST")
	cfg.DBUser = viper.GetString("POSTGRES_USER")
	cfg.DBPassword = viper.GetString("POSTGRES_PASSWORD")
	cfg.APIAddr = viper.GetString("API_ADDR")
	cfg.APIPort = viper.GetString("API_PORT")
	cfg.Secret = viper.GetString("SECRET")
	cfg.Timeout = time.Duration(viper.GetInt("TIMEOUT"))
	cfg.SessionExpires = time.Duration(viper.GetInt("SESSION_EXPIRES"))
	cfg.LogRemoveDays = viper.GetInt("LOG_REMOVE_DAYS")
	cfg.CookieDomain = viper.GetString("COOKIE_DOMAIN")

	return nil
}

func GetConfigs() *config {
	return cfg
}
