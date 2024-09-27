package heartbeat

import (
	"context"
	"fmt"
	l "github.com/ndemeshchenko/zenith/pkg/components/logger"
	"log/slog"
	"time"

	zenithmongo "github.com/ndemeshchenko/zenith/pkg/components/models/alert"

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
		time.Sleep(30 * time.Second)

		// read all records from heartbeats collection
		// check whether any has LastReceivedAt older than 5 minutes

		heartbeats := m.fetchHeartbeats()
		l.Logger.Info("Running heartbeats monitor routine")
		l.Logger.Debug("heartbeats", slog.Any("list", heartbeats))
		for _, heartbeatEvent := range heartbeats {
			l.Logger.Debug("heartbeat", slog.Any("event", heartbeatEvent))
			if !heartbeatEvent.LastReceivedAt.Before(time.Now().UTC().Add(-5 * time.Minute)) {
				//TODO find by fingerprint and delete alert
				alert, err := zenithmongo.FindByFingerprint(m.MongoClient, fmt.Sprintf("heartbeat_%s_%s", heartbeatEvent.Environment, heartbeatEvent.Cluster))
				if err != nil {
					l.Logger.Debug("failed to find alert for cluster %s: %v", heartbeatEvent.Cluster, err)
				}
				if alert == nil {
					l.Logger.Debug("not found", slog.String("alert", heartbeatEvent.Cluster))
					continue
				}
				err = alert.Delete(m.MongoClient)
				if err != nil {
					l.Logger.Error("failed to delete alert for cluster %s: %v", heartbeatEvent.Cluster, err)
				}

			} else {
				l.Logger.Debug("heartbeatEvent is older than 5 minutes", slog.Any("event", heartbeatEvent))
				// create alert
				l.Logger.Debug("create alert", slog.String("cluster", heartbeatEvent.Cluster))
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
					l.Logger.Error("failed to create alert for cluster %s: %v", heartbeatEvent.Cluster, err)
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
		l.Logger.Error("Error finding documents: %v", err)
	}

	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			l.Logger.Error("Error closing cursor: %v\n", err)
		}
	}()

	var heartbeats []heartbeat.Heartbeat
	for cursor.Next(context.Background()) {
		var heartbeatEvent heartbeat.Heartbeat
		err := cursor.Decode(&heartbeatEvent)
		if err != nil {
			l.Logger.Error("Error decoding document: %v", err)
		}
		heartbeats = append(heartbeats, heartbeatEvent)
	}

	if err := cursor.Err(); err != nil {
		l.Logger.Error("Error iterating cursor: %v", err)
	}

	err = cursor.Close(context.Background())
	if err != nil {
		l.Logger.Error("Error closing cursor: %v", err)
	}

	return heartbeats
}
