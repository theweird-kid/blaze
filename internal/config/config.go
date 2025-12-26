package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	DBName   string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		MongoURI: os.Getenv("MONGO_URI"),
		DBName:   os.Getenv("DB_NAME"),
	}
}
