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

//var AdminEmail = "admin798@gmail.com"
//var AdminPassword = "sequence0"
//
//var JWTSecret = os.Getenv("TOKEN_SECRET_KEY")
//
//func EnvTokenSecretKey() string {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	return os.Getenv("TOKEN_SECRET_KEY")
//}
//
////var JWTSecret = func() string {
////	secret := os.Getenv("TOKEN_SECRET_KEY")
////	if secret == "" {
////		secret = "g3Xy8kPfTr5UzN1JQ0MlP2cWzvQnLdVzRzDOZZK-5yI="
////	}
////	return secret
////}()
//
//func EnvMongoURI() string {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	return os.Getenv("MONGO_URI")
//}
//
//func EnvPort() string {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	return os.Getenv("PORT")
//}
