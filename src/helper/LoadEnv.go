package helper

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func ReadEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load the env file")
	}
	return os.Getenv(key)
}
