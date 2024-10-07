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
	EnableTLS          bool
	TLSCertFile        string
	TLSKeyFile         string
	FE_URL             string
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
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	if mongoDatabase == "" {
		mongoDatabase = "zenith"
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	if mongoUsername == "" {
		mongoUsername = "zenith"
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword == "" {
		mongoPassword = "zenith"
	}

	mongoTLS := true
	mongoTLSParam := os.Getenv("MONGO_TLS")
	if mongoTLSParam == "false" {
		mongoTLS = false
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

	// Server TLS Configuration
	enableTLS := false
	enableTLSParam := os.Getenv("ENABLE_TLS")
	if enableTLSParam == "true" {
		enableTLS = true
	}

	tlsCertFile := os.Getenv("TLS_CERT")
	if tlsCertFile == "" {
		tlsCertFile = "server.crt"
	}

	tlsKeyFile := os.Getenv("TLS_KEY")
	if tlsKeyFile == "" {
		tlsKeyFile = "server.key"
	}

	feURL := os.Getenv("FE_URL")
	if feURL == "" {
		feURL = "http://localhost:3000"
	}

	return &Config{
		MongoHost:          mongoHost,
		MongoPort:          mongoPort,
		MongoDatabase:      mongoDatabase,
		MongoUsername:      mongoUsername,
		MongoPassword:      mongoPassword,
		MongoAuthMechanism: mongoAuthMechanism,
		LogLevel:           logLevel,
		AuthToken:          authToken,
		MongoTLS:           mongoTLS,
		EnableTLS:          enableTLS,
		TLSCertFile:        tlsCertFile,
		TLSKeyFile:         tlsKeyFile,
		FE_URL:             feURL,
	}

}
