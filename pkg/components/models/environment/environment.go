package environment

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Environment struct {
	Name  string `json:"Name"`
	Count int    `json:"Count"`
}

func GetAll(mongoClient *mongo.Client) ([]Environment, error) {
	// TODO this is a mock
	// TODO implement fetch from DB
	mockEnv := []Environment{
		{Name: "Development"},
		{Name: "Testing"},
		{Name: "Acceptance"},
		{Name: "Production"},
	}
	type Result struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	collection := mongoClient.Database("zenith").Collection("alerts")
	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"status", "firing"},
		}}},
		{{"$group", bson.D{
			{"_id", "$environment"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
	}
	// Execute the aggregation query
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println("Failed to execute aggregation query:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the results
	var results []Result
	if err := cursor.All(context.Background(), &results); err != nil {
		fmt.Println("Failed to decode aggregation results:", err)
		return nil, err
	}
	for _, result := range results {
		updateCount(&mockEnv, result.ID, result.Count)
	}

	return mockEnv, nil
}

// Update the Count value in the slice of Environment structs based on the Name
func updateCount(envs *[]Environment, name string, count int) {
	for i := range *envs {
		if (*envs)[i].Name == name {
			(*envs)[i].Count = count
			return
		}
	}

	// If no matching Name found, add a new element to the slice
	*envs = append(*envs, Environment{Name: name, Count: count})
}
