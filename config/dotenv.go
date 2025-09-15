package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[CONFIG] [DOTENV] .env file not found, using environment variables")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
