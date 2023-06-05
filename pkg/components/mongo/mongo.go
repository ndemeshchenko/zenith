package mongo

import (
	"context"
	"fmt"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InitDBConnection(config *config.Config) (*mongo.Client, error) {
	mongoDatasource := fmt.Sprintf("mongodb://%s:%s", config.MongoHost, config.MongoPort)
	clientOptions := options.Client().ApplyURI(mongoDatasource).
		SetAuth(options.Credential{
			Username:      config.MongoUsername,
			Password:      config.MongoPassword,
			AuthSource:    config.MongoDatabase,
			AuthMechanism: "SCRAM-SHA-256",
		})

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
