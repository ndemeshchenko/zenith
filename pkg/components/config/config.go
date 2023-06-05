package config

import "os"

type Config struct {
	MongoHost     string
	MongoPort     string
	MongoDatabase string
	MongoUsername string
	MongoPassword string
	Debug         bool
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

	debug := false
	if os.Getenv("DEBUG") == "true" {
		debug = true
	}

	return &Config{
		MongoHost:     mongoHost,
		MongoPort:     mongoPort,
		MongoDatabase: mongoDatabse,
		MongoUsername: mongoUsername,
		MongoPassword: mongoPassword,
		Debug:         debug,
	}

}
