package prometheus

import (
	"time"
)

type WebhookAlertPayloadSubInstanceLabels struct {
	Alertname   string `json:"alertname"`
	App         string `json:"app"`
	Cluster     string `json:"cluster"`
	Container   string `json:"container"`
	Environment string `json:"environment"`
	Namespace   string `json:"namespace"`
	Prometheus  string `json:"prometheus"`
	Severity    string `json:"severity"`
}

type WebhookAlertPayloadSubInstance struct {
	Status      string                               `json:"status"`
	Labels      WebhookAlertPayloadSubInstanceLabels `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Runbook     string `json:"runbook"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

type WebhookAlertPayload struct {
	Receiver    string                           `json:"receiver"`
	Alerts      []WebhookAlertPayloadSubInstance `json:"alerts"`
	GroupLabels struct {
		Alertname   string `json:"alertname"`
		App         string `json:"app"`
		Cluster     string `json:"cluster"`
		Container   string `json:"container"`
		Environment string `json:"environment"`
		Namespace   string `json:"namespace"`
		Prometheus  string `json:"prometheus"`
		Severity    string `json:"severity"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname   string `json:"alertname"`
		App         string `json:"app"`
		Cluster     string `json:"cluster"`
		Container   string `json:"container"`
		Environment string `json:"environment"`
		Namespace   string `json:"namespace"`
		Prometheus  string `json:"prometheus"`
		Severity    string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Description string `json:"description"`
		Runbook     string `json:"runbook"`
		Summary     string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}
