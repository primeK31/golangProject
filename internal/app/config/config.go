package config

import (
	"os"
	"time"

	//"log"

	"github.com/joho/godotenv"
)


type Config struct {
	HTTPPort string
	//SQL_DATABASE_URL string
	JWTSecret string
	DB_USER string 
	DB_PASS string
	DB_HOST string
	DB_PORT string 
	DB_NAME string
	SessionDuration time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	/*if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}*/

	duration := os.Getenv("SESSION_DURATION")
    
    session_duration, _ := time.ParseDuration(duration)

	return &Config{
		HTTPPort: os.Getenv("HTTPPort"),
		// SQL_DATABASE_URL: os.Getenv("SQL_DATABASE_URL"),
		JWTSecret:  os.Getenv("SECRET_KEY"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_NAME: os.Getenv("DB_NAME"),
		SessionDuration: session_duration,
	}
}
