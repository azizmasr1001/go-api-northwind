package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBDriver  string
	DBSource  string
	Port      string
	JWTSecret string
	RedisHost string
	RedisPort string
	RedisPass string
}

func LoadConfig() Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")  // current dir (misal saat run dari main.go)
	viper.AddConfigPath("..") // parent dir (saat test di subfolder)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env not found, trying .env.test")
		viper.SetConfigName(".env.test")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..") // <--- ini penting untuk test
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	return Config{
		DBDriver:  viper.GetString("DB_DRIVER"),
		DBSource:  viper.GetString("DB_SOURCE"),
		Port:      viper.GetString("PORT"),
		JWTSecret: viper.GetString("JWT_SECRET"),
		RedisHost: viper.GetString("REDIS_HOST"),
		RedisPort: viper.GetString("REDIS_PORT"),
		RedisPass: viper.GetString("REDIS_PASS"),
	}
}
