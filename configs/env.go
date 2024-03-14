package configs

import (
	"log"
	"os"
	"strconv"

	env "github.com/golkhandani/shopWise/shared"
	"github.com/joho/godotenv"
)

func SetupEnv() env.Environement {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DATABASE_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		port = 3000
		log.Fatal("Cannot parse port", err)
	}

	result := env.Environement{
		MongodbURI:   mongoURI,
		DatabaseName: dbName,
		JWTSecret:    jwtSecret,
		Port:         port,
	}
	return result
}

var Env = SetupEnv()
