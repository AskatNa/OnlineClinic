package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AdminEmail = "admin798@gmail.com"
var AdminPassword = "sequence0"
var JWTSecret = []byte(os.Getenv("TOKEN_SECRET_KEY"))

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_URI")
}

func EnvPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("PORT")
}
