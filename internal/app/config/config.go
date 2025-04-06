package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)


type Config struct {
	HTTPPort string
	SQL_DATABASE_URL string
	JWTSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		HTTPPort: os.Getenv("HTTPPort"),
		SQL_DATABASE_URL: os.Getenv("SQL_DATABASE_URL"),
		JWTSecret:  os.Getenv("SECRET_KEY"),
	}
}
