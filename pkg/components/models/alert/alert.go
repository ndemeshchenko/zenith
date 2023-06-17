package alert

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetAll(mongoClient *mongo.Client) ([]Alert, error) {
	// get all alerts from the collection
	// Access the "alerts" collection
	collection := mongoClient.Database("zenith").Collection("alerts")

	// Define the filter to find alerts with "status" equal to "open"
	filter := bson.M{"status": "firing"}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"severityCode", 1}}).SetLimit(100)

	// Execute the find operation
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	//cursor, err = collection.Find(context.TODO(), filter, findOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Iterate through the results
	var alerts []Alert
	for cursor.Next(context.TODO()) {
		var alert Alert
		err := cursor.Decode(&alert)
		if err != nil {
			log.Fatal(err)
		}
		alerts = append(alerts, alert)
	}

	// Handle any errors that occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor
	cursor.Close(context.TODO())

	return alerts, nil
}

func GetOne(mongoClient *mongo.Client, id string) (*Alert, error) {
	collection := mongoClient.Database("zenith").Collection("alerts")

	// Execute the find operation
	var result Alert
	objID, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}

	return &result, nil
}

func FindByFingerprint(mongoClient *mongo.Client, fingerprint string) (*Alert, error) {
	collection := mongoClient.Database("zenith").Collection("alerts")

	// Execute the find operation
	var result Alert
	err := collection.FindOne(context.Background(), bson.M{"fingerprint": fingerprint}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}

	return &result, nil
}

func DeleteOne(mongoClient *mongo.Client, id string) error {
	collection := mongoClient.Database("zenith").Collection("alerts")

	// Execute the find operation
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to retrieve documents: %v", err)
	}

	return nil
}

func UpdateStatus(mongoClient *mongo.Client, id string, action string) error {
	collection := mongoClient.Database("zenith").Collection("alerts")

	statuses := map[string]string{"acknowledge": "acknowledged", "resolve": "resolved"}

	// Execute the find operation
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": statuses[action]}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to retrieve documents: %v", err)
	}

	return nil
}

func SeverityToLevel(severity string) int8 {
	switch severity {
	case "critical":
		return 1
	case "warning":
		return 2
	case "info":
		return 3
	default:
		return 4
	}
}
