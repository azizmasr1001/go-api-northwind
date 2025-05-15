package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver  string
	DBSource  string
	Port      string
	JWTSecret string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return Config{
		DBDriver:  viper.GetString("DB_DRIVER"),
		DBSource:  viper.GetString("DB_SOURCE"),
		Port:      viper.GetString("PORT"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}
}
