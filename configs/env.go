package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() string{
	godotenv.Load()
	return os.Getenv("MONGO_URL")
}