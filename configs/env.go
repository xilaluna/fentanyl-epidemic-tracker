package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadMongoEnv() string{
	godotenv.Load()
	return os.Getenv("MONGO_URL")
}