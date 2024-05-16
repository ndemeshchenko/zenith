package alert

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Alert struct {
	ID               string            `bson:"_id,omitempty"`
	Resource         string            `bson:"resource,omitempty"`
	Event            string            `bson:"event,omitempty"`
	Environment      string            `bson:"environment,omitempty"`
	Cluster          string            `bson:"cluster,omitempty"`
	SeverityCode     int8              `bson:"severityCode,omitempty"`
	SeverityName     string            `bson:"severityName,omitempty"`
	Correlate        map[string]string `bson:"correlate,omitempty"`
	Status           string            `bson:"status,omitempty"`
	Service          string            `bson:"service,omitempty"`
	Group            string            `bson:"group,omitempty"`
	Value            string            `bson:"value,omitempty"` // TODO: what is this?
	Text             string            `bson:"text,omitempty"`
	Summary          string            `bson:"summary,omitempty"`
	Tags             map[string]string `bson:"tags,omitempty"`
	Attributes       map[string]string `bson:"attributes,omitempty"`
	Origin           string            `bson:"origin,omitempty"`
	Type             string            `bson:"type,omitempty"`
	Fingerprint      string            `bson:"fingerprint,omitempty"`
	GeneratorURL     string            `bson:"generatorURL,omitempty"`
	RunbookURL       string            `bson:"runbookURL,omitempty"`
	CreateTime       time.Time         `bson:"createTime,omitempty"`
	Timeout          time.Time         `bson:"timeout,omitempty"`
	RawData          string            `bson:"rawData,omitempty"`
	DuplicateCount   int               `bson:"duplicateCount,omitempty"`
	Repeat           bool              `bson:"repeat,omitempty"`
	PreviousSeverity string            `bson:"previousSeverity,omitempty"`
	TrendIndication  string            `bson:"trendIndication,omitempty"`
	ReceiveTime      time.Time         `bson:"receiveTime,omitempty"`
	UpdateTime       time.Time         `bson:"updateTime,omitempty"`
	//History          []History //TODO implement history

}

func (a *Alert) Validate() error {
	var err error
	if a.Event == "" {
		err = errors.Join(err, errors.New("event is required"))
	}
	if a.Environment == "" {
		err = errors.Join(err, errors.New("environment is required"))
	}

	if a.SeverityCode == 0 {
		err = errors.Join(err, errors.New("severityCode is required"))
	}

	if a.SeverityName == "" {
		err = errors.Join(err, errors.New("severityName is required"))
	}

	if a.Status == "" {
		err = errors.Join(err, errors.New("status is required"))
	}

	//if a.Origin == "" {
	//	err = errors.Join(err, errors.New("origin is required"))
	//}

	if a.Type == "" {
		err = errors.Join(err, errors.New("type is required"))
	}

	if a.Fingerprint == "" {
		err = errors.Join(err, errors.New("fingerprint is required"))
	}

	if a.CreateTime.IsZero() {
		err = errors.Join(err, errors.New("createTime is required"))
	}

	return err
}

func (a *Alert) Create(c *mongo.Client) (interface{}, error) {
	err := a.Validate()
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to validate alert: %v", err)
	}
	// Access the "alerts" collection
	collection := c.Database("zenith").Collection("alerts")

	// Insert the alert into the collection
	one, err := collection.InsertOne(context.Background(), a)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to insert alert: %v", err)
	}
	return one.InsertedID, nil
}

// Upsert the alert into the collection
func (a *Alert) Upsert(c *mongo.Client) (interface{}, error) {
	err := a.Validate()
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to validate alert: %v", err)
	}
	// Access the "alerts" collection
	collection := c.Database("zenith").Collection("alerts")

	// Upsert the alert into the collection
	dup := a.FindDuplicate(c)
	if dup != nil {
		id, err := dup.UpdateDuplicated(c, a)
		if err != nil {
			return primitive.ObjectID{}, fmt.Errorf("failed to update duplicated alert: %v", err)
		}
		return id, nil
	}

	// Insert the alert into the collection
	one, err := collection.InsertOne(context.Background(), a)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to insert alert: %v", err)
	}
	return one.InsertedID, nil

}

func (a *Alert) FindDuplicate(c *mongo.Client) *Alert {
	// Access the "alerts" collection
	collection := c.Database("zenith").Collection("alerts")
	//Find the alert by fingerprint value
	filter := bson.D{{"fingerprint", a.Fingerprint}}
	var result Alert
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Errorf("failed to find alert: %v", err)
		return nil
	}
	if result.Fingerprint == a.Fingerprint {
		return &result
	}
	return nil
}

func (a *Alert) UpdateDuplicated(c *mongo.Client, wa *Alert) (interface{}, error) {
	var err error
	// Access the "alerts" collection
	collection := c.Database("zenith").Collection("alerts")
	//Find the alert by fingerprint value
	filter := bson.D{{"fingerprint", a.Fingerprint}}
	//UpdateDuplicated the alert
	update := bson.D{
		{"$set", bson.D{
			{"repeat", true},
			{"duplicateCount", a.DuplicateCount + 1},
			{"updateTime", time.Now()},
			{"status", wa.Status},
		}},
	}
	one, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to update alert: %v", err)
	}
	return one.UpsertedID, nil
}

func (a *Alert) Delete(c *mongo.Client) error {
	// Access the "alerts" collection
	collection := c.Database("zenith").Collection("alerts")
	//Find the alert by fingerprint value
	filter := bson.D{{"fingerprint", a.Fingerprint}}
	//Delete the alert
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %v", err)
	}
	return nil
}
