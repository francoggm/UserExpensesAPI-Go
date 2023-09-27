package configs

import "github.com/spf13/viper"

type config struct {
	DB string
	DBPort string
	DBHost string
	DBUser string
	DBPassword string
	APIPort string
	Secret string
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
	cfg.DBPort = viper.GetString("DBPORT")
	cfg.DBHost = viper.GetString("DBHOST")
	cfg.DBUser = viper.GetString("DBUSER")
	cfg.DBPassword = viper.GetString("DBPASSWORD")
	cfg.APIPort = viper.GetString("APIPORT")
	cfg.Secret = viper.GetString("SECRET")

	return nil
}

func GetConfigs() *config {
	return cfg
}