package configs

import "github.com/spf13/viper"

type config struct {
	DB         string
	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	APIAddr    string
	APIPort    string
	Secret     string
}

var cfg *config

func Load() error {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	cfg = new(config)

	cfg.DB = viper.GetString("DB")
	cfg.DBPort = viper.GetString("DB_PORT")
	cfg.DBHost = viper.GetString("DB_HOST")
	cfg.DBUser = viper.GetString("DB_USER")
	cfg.DBPassword = viper.GetString("DB_PASSWORD")
	cfg.APIAddr = viper.GetString("API_ADDR")
	cfg.APIPort = viper.GetString("API_PORT")
	cfg.Secret = viper.GetString("SECRET")

	return nil
}

func GetConfigs() *config {
	return cfg
}
