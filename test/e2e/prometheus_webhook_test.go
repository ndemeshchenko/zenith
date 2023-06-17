package e2e

import (
	"context"
	"encoding/json"
	"github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	prometheusWebhook "github.com/ndemeshchenko/zenith/pkg/components/webhooks/prometheus"
	prometheusTest "github.com/ndemeshchenko/zenith/test/e2e/webhooks/prometheus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"strings"
	"testing"

	e2emongo "github.com/ndemeshchenko/zenith/test/e2e/mongo"
)

func setupMongo(t *testing.T) (*mongo.Client, func(t *testing.T)) {
	ctx := context.Background()
	container, err := e2emongo.StartContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	///
	mongoClient, err := e2emongo.TestContainerMongoClient(ctx, container)

	// return teardown function
	return mongoClient, func(t *testing.T) {
		t.Cleanup(func() {
			if err := container.Terminate(ctx); err != nil {
				t.Fatalf("failed to terminate container: %s", err)
			}
		})

	}
}

func TestCreateValidPromWebHook(t *testing.T) {
	mongoClient, teardown := setupMongo(t)
	// perform assertions
	defer teardown(t)

	jsonAlert, _ := prometheusTest.GeneratePrometheusWebhook()
	jsonString, err := json.Marshal(jsonAlert)

	payloadReader := strings.NewReader(string(jsonString))
	payload := io.NopCloser(payloadReader)
	err = prometheusWebhook.ProcessWebhookAlert(payload, mongoClient)
	assert.Nil(t, err)
}

func TestCreateInvalidPromWebHook(t *testing.T) {
	mongoClient, teardown := setupMongo(t)
	// perform assertions
	defer teardown(t)

	jsonAlert, _ := prometheusTest.GeneratePrometheusWebhook()
	jsonAlert.Alerts = nil
	jsonString, err := json.Marshal(jsonAlert)

	log.Println(string(jsonString))

	payloadReader := strings.NewReader(string(jsonString))
	payload := io.NopCloser(payloadReader)
	err = prometheusWebhook.ProcessWebhookAlert(payload, mongoClient)

	assert.NotNil(t, err)
	assert.Equal(t, "empty alerts", err.Error())
}

func TestUpdateStatus(t *testing.T) {
	mongoClient, teardown := setupMongo(t)
	// perform assertions
	defer teardown(t)

	jsonAlert, _ := prometheusTest.GeneratePrometheusWebhook()
	jsonString, err := json.Marshal(jsonAlert)

	payloadReader := strings.NewReader(string(jsonString))
	payload := io.NopCloser(payloadReader)
	err = prometheusWebhook.ProcessWebhookAlert(payload, mongoClient)

	assert.Nil(t, err)

	alerts, err := alert.GetAll(mongoClient)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(alerts))

	alertx := alerts[0]
	assert.Equal(t, "firing", alertx.Status)

	err = alert.UpdateStatus(mongoClient, alertx.ID, "resolve")
	assert.Nil(t, err)

	alertUpd, err := alert.FindByFingerprint(mongoClient, alertx.Fingerprint)
	assert.NotNil(t, alertUpd)
	assert.Equal(t, "resolved", alertUpd.Status)
}

func TestDeleteAction(t *testing.T) {
	mongoClient, teardown := setupMongo(t)
	// perform assertions
	defer teardown(t)

	jsonAlert, _ := prometheusTest.GeneratePrometheusWebhook()
	jsonString, err := json.Marshal(jsonAlert)

	payloadReader := strings.NewReader(string(jsonString))
	payload := io.NopCloser(payloadReader)
	err = prometheusWebhook.ProcessWebhookAlert(payload, mongoClient)

	assert.Nil(t, err)

	alerts, err := alert.GetAll(mongoClient)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(alerts))

	alertx := alerts[0]
	assert.Equal(t, "firing", alertx.Status)

	err = alert.DeleteOne(mongoClient, alertx.ID)
	assert.Nil(t, err)

	alertUpd, err := alert.FindByFingerprint(mongoClient, alertx.Fingerprint)
	assert.Nil(t, alertUpd)
	//assert.Equal(t, "resolved", alertUpd.Status)
}

func TestDuplicateAlertsPromWebHook(t *testing.T) {
	mongoClient, teardown := setupMongo(t)
	// perform assertions
	defer teardown(t)

	jsonAlert1, _ := prometheusTest.GeneratePrometheusWebhook()
	jsonString1, err := json.Marshal(jsonAlert1)

	log.Println(string(jsonString1))

	payloadReader1 := strings.NewReader(string(jsonString1))
	payload1 := io.NopCloser(payloadReader1)
	err = prometheusWebhook.ProcessWebhookAlert(payload1, mongoClient)

	jsonAlert2 := jsonAlert1
	jsonAlert2.Alerts[0].Labels.Severity = "warning"
	jsonString2, err := json.Marshal(jsonAlert2)

	payloadReader2 := strings.NewReader(string(jsonString2))
	payload2 := io.NopCloser(payloadReader2)
	err = prometheusWebhook.ProcessWebhookAlert(payload2, mongoClient)

	alertRecord, err := alert.FindByFingerprint(mongoClient, jsonAlert1.Alerts[0].Fingerprint)

	assert.NotNil(t, alertRecord)
	assert.Nil(t, err)
	assert.Equal(t, 1, alertRecord.DuplicateCount)
}
