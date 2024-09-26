package config

import (
	"fmt"
	"math/rand"
	"os"
)

type Config struct {
	MongoHost          string
	MongoPort          string
	MongoDatabase      string
	MongoUsername      string
	MongoPassword      string
	MongoTLS           bool
	MongoAuthMechanism string
	Debug              bool
	AuthToken          string
}

func InitConfigurations() *Config {
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "localhost"
	}

	mongoPort := os.Getenv("MONGO_PORT")
	if mongoPort == "" {
		mongoPort = "27017"
	}
	mongoDatabse := os.Getenv("MONGO_DATABASE")
	if mongoDatabse == "" {
		mongoDatabse = "zenith"
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	if mongoUsername == "" {
		mongoUsername = "zenith"
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword == "" {
		mongoPassword = "zenith"
	}

	mongoTLS := os.Getenv("MONGO_TLS")
	if mongoTLS == "" {
		mongoTLS = "false"
	}

	mongoAuthMechanism := os.Getenv("MONGO_AUTH_MECHANISM")

	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		b := make([]byte, 32)
		rand.Read(b)
		authToken = fmt.Sprintf("%x", b)
	}

	return &Config{
		MongoHost:          mongoHost,
		MongoPort:          mongoPort,
		MongoDatabase:      mongoDatabse,
		MongoUsername:      mongoUsername,
		MongoPassword:      mongoPassword,
		MongoAuthMechanism: mongoAuthMechanism,
		Debug:              debug,
		AuthToken:          authToken,
	}

}
