package config

import "os"

type Config struct {
	MongoURI string
	DBName   string
}

func Load() *Config {
	return &Config{
		MongoURI: os.Getenv("MONGO_URI"),
		DBName:   os.Getenv("DB_NAME"),
	}
}
