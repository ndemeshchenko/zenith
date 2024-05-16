package heartbeat

import (
	"context"
	"fmt"
	zenithmongo "github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	"log"
	"time"

	"github.com/ndemeshchenko/zenith/pkg/components/models/heartbeat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Monitor struct with mongoclient
type Monitor struct {
	MongoClient *mongo.Client
}

func NewMonitor(mongoClient *mongo.Client) *Monitor {
	return &Monitor{
		MongoClient: mongoClient,
	}
}

// Run function to start heartbeats monitor routine
func (m *Monitor) Run() {
	for {
		// print date in YYYY-MM-DD HH:MM:SS format
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(30 * time.Second)

		// read all records from heartbeats collection
		// check whether any has LastReceivedAt older than 5 minutes

		heartbeats := m.fetchHeartbeats()
		log.Println("heartbeats: ", heartbeats)
		for _, heartbeatEvent := range heartbeats {
			log.Printf("heartbeatEvent: %+v", heartbeatEvent)
			if !heartbeatEvent.LastReceivedAt.Before(time.Now().Add(-5 * time.Minute)) {
				//TODO find by fingerprint and delete alert
				log.Printf("delete alert for cluster %s", heartbeatEvent.Cluster)
				alert, err := zenithmongo.FindByFingerprint(m.MongoClient, fmt.Sprintf("heartbeat_%s_%s", heartbeatEvent.Environment, heartbeatEvent.Cluster))
				if err != nil {
					log.Printf("failed to find alert for cluster %s: %v", heartbeatEvent.Cluster, err)
				}
				if alert == nil {
					log.Printf("alert for cluster %s not found", heartbeatEvent.Cluster)
					continue
				}
				err = alert.Delete(m.MongoClient)
				if err != nil {
					log.Printf("failed to delete alert for cluster %s: %v", heartbeatEvent.Cluster, err)
				}

			} else {
				log.Println("heartbeatEvent is older than 5 minutes: ", heartbeatEvent)
				// create alert
				log.Printf("create alert for cluster %s", heartbeatEvent.Cluster)
				alert := zenithmongo.Alert{
					Resource:     heartbeatEvent.Cluster,
					Event:        "Heartbeat",
					Environment:  heartbeatEvent.Environment,
					Cluster:      heartbeatEvent.Cluster,
					SeverityCode: 1,
					SeverityName: "Critical",
					Correlate:    map[string]string{},
					Status:       "firing",
					Type:         "heartbeatEvent",
					Fingerprint:  fmt.Sprintf("heartbeat_%s_%s", heartbeatEvent.Environment, heartbeatEvent.Cluster),
					CreateTime:   time.Now(),
				}
				_, err := alert.Upsert(m.MongoClient)
				if err != nil {
					log.Printf("failed to create alert for cluster %s: %v", heartbeatEvent.Cluster, err)
				}
			}
		}
	}
}

func (m *Monitor) fetchHeartbeats() []heartbeat.Heartbeat {
	// read all records from heartbeats collection
	// check whether any has LastReceivedAt older than 5 minutes
	//collection := c.Database("zenith").Collection("alerts")
	collection := m.MongoClient.Database("zenith").Collection("heartbeats")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Printf("Error finding documents: %v", err)
	}

	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}()

	var heartbeats []heartbeat.Heartbeat
	for cursor.Next(context.Background()) {
		var heartbeatEvent heartbeat.Heartbeat
		err := cursor.Decode(&heartbeatEvent)
		if err != nil {
			log.Printf("Error decoding document: %v", err)
		}
		heartbeats = append(heartbeats, heartbeatEvent)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error iterating cursor: %v", err)
	}

	err = cursor.Close(context.Background())
	if err != nil {
		log.Printf("Error closing cursor: %v", err)
	}

	return heartbeats
}
