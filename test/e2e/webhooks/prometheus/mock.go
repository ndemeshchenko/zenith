package prometheus

import (
	"github.com/brianvoe/gofakeit/v6"
	prometheus "github.com/ndemeshchenko/zenith/pkg/components/webhooks/prometheus"
	"math/rand"
	"time"
)

func randomEventName() string {
	eventNames := []string{
		"AlertmanagerClusterFailedToSendAlerts",
		"AlertmanagerFailedToSendAlerts",
		"ContainerCpuUsage",
		"ContainerMemoryUsage",
		"ExcessiveContainerRestarts",
		"FailedEvictedPods",
		"HighServerErrorRate",
		"KubeCPUOvercommit",
		"KubeContainerWaiting",
		"KubeDeploymentReplicasMismatch",
		"KubeMemoryOvercommit",
		"KubePodCrashLooping",
		"KubePodNotReady",
		"KubeStatefulSetReplicasMismatch",
		"NGINXTooMany400s",
		"NGINXTooMany500s",
		"ServiceNotReady",
	}

	rand.Seed(time.Now().UnixNano())
	return eventNames[rand.Intn(len(eventNames))]
}

// generate output produced by promethues webhook
func GeneratePrometheusWebhook() (prometheus.WebhookAlertPayload, error) {
	gofakeit.Seed(time.Now().UnixNano())

	var alertsList []prometheus.WebhookAlertPayloadSubInstance
	alertInstance := GenerateWebhookAlertPayloadSubInstance()
	alertsList = append(alertsList, alertInstance)

	payload := prometheus.WebhookAlertPayload{
		Receiver: "prometheus",
		Alerts:   alertsList,
	}
	return payload, nil
}

func GenerateWebhookAlertPayloadSubInstance() prometheus.WebhookAlertPayloadSubInstance {
	gofakeit.Seed(time.Now().UnixNano())

	return prometheus.WebhookAlertPayloadSubInstance{
		Status: "firing",
		Labels: prometheus.WebhookAlertPayloadSubInstanceLabels{
			Alertname:   randomEventName(),
			App:         gofakeit.AppName(),
			Cluster:     gofakeit.AppName(),
			Container:   gofakeit.AppName(),
			Environment: "development",
			Namespace:   gofakeit.Animal(),
			Prometheus:  gofakeit.BeerHop(),
			Severity:    "critical",
		},
		StartsAt:    gofakeit.Date(),
		EndsAt:      gofakeit.Date(),
		Fingerprint: gofakeit.UUID(),
	}
}
