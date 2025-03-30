package config

import "time"

type Config struct {
	ServerAddress   string
	JWTSecret       string
	TokenExpiration time.Duration
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress:   ":8080",
		JWTSecret:       "your-secret-key", // Use .env in prod
		TokenExpiration: 24 * time.Hour,
		DBHost:          "127.0.0.1",
		DBName:          "chat_app",
		DBPort:          "5432",
		DBUser:          "postgres",
		DBPassword:      "password",
	}
}
