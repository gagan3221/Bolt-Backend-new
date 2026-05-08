package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	Port         string
	JWTSecret    string
	AlchemyURL   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		MongoURI:     os.Getenv("MONGO_URI"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		Port:         os.Getenv("PORT"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		AlchemyURL:   os.Getenv("ALCHEMY_URL"),
	}
}