package config

import (
	"crypto/rand"
	"fmt"
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
	LogLevel           string
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

	logLevel := "info"
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = "debug"
	}
	if os.Getenv("LOG_LEVEL") == "trace" {
		logLevel = "trace"
	}

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		authToken = fmt.Sprintf("%x", b)
	}

	return &Config{
		MongoHost:          mongoHost,
		MongoPort:          mongoPort,
		MongoDatabase:      mongoDatabse,
		MongoUsername:      mongoUsername,
		MongoPassword:      mongoPassword,
		MongoAuthMechanism: mongoAuthMechanism,
		LogLevel:           logLevel,
		AuthToken:          authToken,
	}

}
