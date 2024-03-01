package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const errorString = "Package Config Error: "

type Config struct{
	Port string
	DbURL string
	DefaultPageSize int32
	DefaultPage int32
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(errorString+"Error loading .env file: ", err)
		log.Println(errorString+"Using predefined credentials")
		return Config{
			Port: "8000",
			DbURL: "postgresql://postgres:staphone@127.0.0.1:5432/go_auth?sslmode=disable",
			DefaultPageSize: 100,
			DefaultPage: 1,
		}
	}

	return Config{
		Port: os.Getenv("PORT"),
		DbURL: os.Getenv("DB_URL"),
		DefaultPageSize: 100,
		DefaultPage: 1,
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}