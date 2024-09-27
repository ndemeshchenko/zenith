package heartbeat

import (
	"context"
	"errors"
	"fmt"
	"github.com/ndemeshchenko/zenith/pkg/components/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Heartbeat represents a heartbeat event
type Heartbeat struct {
	Cluster        string    `json:"Cluster" bson:"cluster"`
	Environment    string    `json:"Environment" bson:"environment"`
	Fingerprint    string    `json:"Fingerprint" bson:"fingerprint"`
	LastReceivedAt time.Time `json:"LastReceivedAt" bson:"lastReceivedAt"`
	HasAlerts      bool      `json:"HasAlerts" bson:"hasAlerts"`
}

// GetAll fetches all heartbeat events from the database
func GetAll(filter bson.M, mongoClient *mongo.Client) ([]Heartbeat, error) {
	// get all alerts from the collection
	// Access the "alerts" collection
	collection := mongoClient.Database("zenith").Collection("heartbeats")

	findOptions := options.Find()
	findOptions.
		SetSort(bson.D{{"lastReceivedAt", 1}}).
		SetLimit(100)

	pipeline := mongo.Pipeline{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "alerts"},
					{"let",
						bson.D{
							{"cluster", "$cluster"},
							{"env", "$env"},
						},
					},
					{"pipeline",
						bson.A{bson.D{{"$match", bson.D{{"$expr", bson.D{{"$and", bson.A{bson.D{{"$eq", bson.A{"$cluster", "$$cluster"}}}, bson.D{{"$eq", bson.A{"$env", "$$env"}}}, bson.D{{"$eq", bson.A{"$event", "Heartbeat"}}}}}}}}}}},
					},
					{"as", "alertsData"},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{{"hasAlerts", bson.D{{"$cond", bson.D{{"if", bson.D{{"$gt", bson.A{bson.D{{"$size", "$alertsData"}}, 0}}}}, {"then", true}, {"else", false}}}}}},
			},
		},
		bson.D{{"$project", bson.D{{"alertsData", 0}}}},
	}

	// Execute the find operation
	//cursor, err := collection.Find(context.Background(), filter, findOptions)
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents: %v", err)
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			logger.Logger.Error("Error closing cursor: %v", err)
		}
	}()

	// Iterate through the results
	var heartbeats []Heartbeat
	for cursor.Next(context.TODO()) {
		var heartbeat Heartbeat
		err := cursor.Decode(&heartbeat)
		if err != nil {
			log.Fatal(err)
		}
		heartbeats = append(heartbeats, heartbeat)
	}

	// Handle any errors that occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor
	_ = cursor.Close(context.Background())

	return heartbeats, nil
}

// Create inserts a new heartbeat event into the database
func (heartbeatEvent *Heartbeat) Create(c *mongo.Client) (interface{}, error) {
	err := heartbeatEvent.Validate()
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to validate alert: %v", err)
	}
	// Access the "events" collection
	collection := c.Database("zenith").Collection("heartbeats")

	filter := bson.D{
		{"cluster", heartbeatEvent.Cluster},
		{"environment", heartbeatEvent.Environment},
		{"fingerprint", heartbeatEvent.Fingerprint},
	}
	update := bson.D{{"$set", bson.D{{"lastReceivedAt", time.Now()}}}}

	opts := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(context.TODO(), filter, update, opts)

	//if result != nil {
	//	logger.Logger.Debug("documents updated", result.ModifiedCount)
	//	logger.Logger.Debug("documents upserted", result.UpsertedCount)
	//}

	// Upsert the heartbeat into the collection
	//one, err := collection.UpdateOne(context.Background(), heartbeatEvent, heartbeatEvent)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to insert heartbeat: %v", err)

	}

	//// Insert event into the heartbeats collection
	//event, err := collection.InsertOne(context.Background(), heartbeatEvent)
	//if err != nil {
	//	return primitive.ObjectID{}, fmt.Errorf("failed to insert heartbeat event: %v", err)
	//}
	return result.UpsertedID, nil
}

// Validate checks if the heartbeat event is valid
func (heartbeatEvent *Heartbeat) Validate() error {
	var err error
	if heartbeatEvent.Cluster == "" {
		err = errors.Join(err, errors.New("cluster name is required"))
	}

	return err
}
