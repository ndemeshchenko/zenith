package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	l "github.com/ndemeshchenko/zenith/pkg/components/logger"
	zenithmongo "github.com/ndemeshchenko/zenith/pkg/components/models/alert"
	"github.com/ndemeshchenko/zenith/pkg/components/models/heartbeat"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"time"
)

func ProcessWebhookAlert(payload io.ReadCloser, mongoClient *mongo.Client) error {
	jsonData, err := io.ReadAll(payload)
	if err != nil {
		return err
	}

	l.Logger.Debug("received alert: ", string(jsonData))

	var webhookAlertPayload WebhookAlertPayload
	err = json.Unmarshal(jsonData, &webhookAlertPayload)
	if err != nil {
		l.Logger.Error("failed to unmarshal json: %v", err)
		return err
	}

	alert, err := transformWebhookAlert(webhookAlertPayload)
	if err != nil {
		_ = fmt.Errorf("failed to transform alert: %v", err)
		return err
	}

	if alert.Event == "Watchdog" {

		l.Logger.Debug("received watchdog alert for cluster %+v", alert)
		// create heartbeatEvent event
		heartbeatEvent := heartbeat.Heartbeat{
			Cluster:        alert.Cluster,
			Environment:    alert.Environment,
			Fingerprint:    fmt.Sprintf("heartbeat_%s_%s", alert.Environment, alert.Cluster),
			LastReceivedAt: time.Now(),
		}

		_, err = heartbeatEvent.Create(mongoClient)

		return err
	}

	dst := &bytes.Buffer{}
	_ = json.Compact(dst, jsonData)
	alert.RawData = dst.String()

	dupAlert := alert.FindDuplicate(mongoClient)
	if dupAlert != nil {
		id, err := dupAlert.UpdateDuplicated(mongoClient, &alert)
		if err != nil {
			l.Logger.Error("failed to update duplicated alert: ", id, err)
			return err
		}
		return nil
	}

	_, err = alert.Create(mongoClient)
	return err
}

func transformWebhookAlert(wap WebhookAlertPayload) (zenithmongo.Alert, error) {
	if len(wap.Alerts) == 0 {
		return zenithmongo.Alert{}, fmt.Errorf("empty alerts")
	}

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
		UpdateTime:   time.Now(),
		GeneratorURL: wap.Alerts[0].GeneratorURL,
		RunbookURL:   wap.Alerts[0].Annotations.Runbook,
		Summary:      wap.Alerts[0].Annotations.Summary,
	}

	return alert, nil
}
