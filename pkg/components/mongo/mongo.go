package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ndemeshchenko/zenith/pkg/components/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InitDBConnection(config *config.Config) (*mongo.Client, error) {
	mongoDatasource := fmt.Sprintf("mongodb://%s:%s", config.MongoHost, config.MongoPort)
	fmt.Println("Connecting to MongoDB: ", mongoDatasource)

	clientCredentials := options.Credential{
		Username:   config.MongoUsername,
		Password:   config.MongoPassword,
		AuthSource: config.MongoDatabase,
	}
	if config.MongoAuthMechanism != "" {
		clientCredentials.AuthMechanism = config.MongoAuthMechanism
	}

	clientOptions := options.Client().ApplyURI(mongoDatasource).
		SetAuth(clientCredentials).SetTLSConfig(&tls.Config{})

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
