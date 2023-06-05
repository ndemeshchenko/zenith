package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	zenithmongo "github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"time"
)

func ProcessWebhookAlert(payload io.ReadCloser, mongoClient *mongo.Client) error {
	//fmt.Println("ProcessWebhookAlert")
	jsonData, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	var webhookAlertPayload WebhookAlertPayload
	err = json.Unmarshal(jsonData, &webhookAlertPayload)
	if err != nil {
		fmt.Errorf("failed to unmarshal json: %v", err)
		return err
	}

	alert, err := transformWebhookAlert(webhookAlertPayload)
	if err != nil {
		fmt.Errorf("failed to transform alert: %v", err)
		return err
	}

	dst := &bytes.Buffer{}
	_ = json.Compact(dst, jsonData)
	alert.RawData = dst.String()

	dupAlert := alert.FindDuplicate(mongoClient)
	if dupAlert != nil {
		err = dupAlert.UpdateDuplicated(mongoClient, &alert)
		return nil
	}

	alert.Create(mongoClient)
	return nil
}

func transformWebhookAlert(wap WebhookAlertPayload) (zenithmongo.Alert, error) {
	alert := zenithmongo.Alert{
		Event:        wap.Alerts[0].Labels.Alertname,
		Environment:  wap.Alerts[0].Labels.Environment,
		Cluster:      wap.Alerts[0].Labels.Cluster,
		Status:       wap.Alerts[0].Status,
		Type:         "prometheus", //TODO dispatch at middleware level
		SeverityCode: zenithmongo.SeverityToLevel(wap.Alerts[0].Labels.Severity),
		SeverityName: wap.Alerts[0].Labels.Severity,
		Text:         wap.Alerts[0].Annotations.Description,
		Fingerprint:  wap.Alerts[0].Fingerprint,
		CreateTime:   time.Now(),
		GeneratorURL: wap.Alerts[0].GeneratorURL,
		RunbookURL:   wap.Alerts[0].Annotations.Runbook,
		Summary:      wap.Alerts[0].Annotations.Summary,
	}

	return alert, nil
}
