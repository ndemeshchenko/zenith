package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// mongodbContainer represents the mongodb container type used in the module
type mongodbContainer struct {
	testcontainers.Container
}

// StartContainer creates an instance of the mongodb container type
func StartContainer(ctx context.Context) (*mongodbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mongodbContainer{Container: container}, nil
}

func TestContainerMongoClient(ctx context.Context, container *mongodbContainer) (*mongo.Client, error) {
	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		return nil, fmt.Errorf("error creating mongo client: %w", err)
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo: %w", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging mongo: %w", err)
	}

	return mongoClient, nil
}
