package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	url := "http://localhost:8080/api/webhooks/prometheus"
	token := "SUPERTESTENVTOKEN"
	// define slice of string
	// Define the payload from the example
	payload1 := `{"receiver":"default-route-receiver","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"PrometheusOperatorSyncFailed","cluster":"dev-gen-de02","container":"kube-prometheus-stack","controller":"alertmanager","endpoint":"http","environment":"Development","instance":"192.168.70.202:8080","job":"kube-prometheus-stack-operator","namespace":"monitoring","pod":"kube-prometheus-stack-operator-7545856747-ld4zj","prometheus":"monitoring/kube-prometheus-stack-prometheus","service":"kube-prometheus-stack-operator","severity":"warning","status":"failed"},"annotations":{"description":"Controller alertmanager in monitoring namespace fails to reconcile 1 objects.","runbook_url":"https://runbooks.prometheus-operator.dev/runbooks/prometheus-operator/prometheusoperatorsyncfailed","summary":"Last controller reconciliation failed"},"startsAt":"2023-05-26T07:13:03.279Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://prometheus.dev.company.io/graph?g0.expr=min_over_time%28prometheus_operator_syncs%7Bjob%3D%22kube-prometheus-stack-operator%22%2Cnamespace%3D%22monitoring%22%2Cstatus%3D%22failed%22%7D%5B5m%5D%29+%3E+0&g0.tab=1","fingerprint":"b6f2cc7bceec39bc"}],"groupLabels":{"alertname":"PrometheusOperatorSyncFailed"},"commonLabels":{"alertname":"PrometheusOperatorSyncFailed","cluster":"dev-gen-de02","container":"kube-prometheus-stack","controller":"alertmanager","endpoint":"http","environment":"Development","instance":"192.168.70.202:8080","job":"kube-prometheus-stack-operator","namespace":"monitoring","pod":"kube-prometheus-stack-operator-7545856747-ld4zj","prometheus":"monitoring/kube-prometheus-stack-prometheus","service":"kube-prometheus-stack-operator","severity":"warning","status":"failed"},"commonAnnotations":{"description":"Controller alertmanager in monitoring namespace fails to reconcile 1 objects.","runbook_url":"https://runbooks.prometheus-operator.dev/runbooks/prometheus-operator/prometheusoperatorsyncfailed","summary":"Last controller reconciliation failed"},"externalURL":"https://alertmanager.dev.company.io","version":"4","groupKey":"{}:{alertname=\"PrometheusOperatorSyncFailed\"}","truncatedAlerts":0}`

	payload2 := `{"receiver":"default-route-receiver","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"KubeAggregatedAPIDown","cluster":"dev-gen-de02","environment":"Development","name":"v1beta1.metrics.k8s.io","namespace":"default","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"annotations":{"description":"Kubernetes aggregated API v1beta1.metrics.k8s.io/default has been only 0% available over the last 10m.","runbook_url":"https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeaggregatedapidown","summary":"Kubernetes aggregated API is down."},"startsAt":"2023-05-26T07:08:09.502Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://prometheus.dev.company.io/graph?g0.expr=%281+-+max+by+%28name%2C+namespace%2C+cluster%29+%28avg_over_time%28aggregator_unavailable_apiservice%5B10m%5D%29%29%29+%2A+100+%3C+85&g0.tab=1","fingerprint":"a00efbd2e54f47e4"}],"groupLabels":{"alertname":"KubeAggregatedAPIDown"},"commonLabels":{"alertname":"KubeAggregatedAPIDown","cluster":"dev-gen-de02","environment":"Development","name":"v1beta1.metrics.k8s.io","namespace":"default","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"commonAnnotations":{"description":"Kubernetes aggregated API v1beta1.metrics.k8s.io/default has been only 0% available over the last 10m.","runbook_url":"https://runbooks.prometheus-operator.dev/runbooks/kubernetes/kubeaggregatedapidown","summary":"Kubernetes aggregated API is down."},"externalURL":"https://alertmanager.dev.company.io","version":"4","groupKey":"{}:{alertname=\"KubeAggregatedAPIDown\"}","truncatedAlerts":0}`

	payload3 := `{"receiver":"default-route-receiver","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"ContainerMemoryUsage","app":"internal-ingress-nginx-controller","cluster":"dev-gen-de02","container":"controller","environment":"Development","namespace":"internal-ingress-nginx","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"annotations":{"description":"Memory consumption of container controller of internal-ingress-nginx-controller in namespace internal-ingress-nginx\nis at 144.70% (> 100%) of the Memory requests for more than 10 minutes.","runbook":"https://company.atlassian.net/wiki/spaces/AR/pages/4075323867/Prometheus+alerts#Prometheusalerts-ContainerHighMemoryUsage","summary":"Container High Memory Usage"},"startsAt":"2023-05-26T07:13:10.896Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://prometheus.dev.company.io/graph?g0.expr=%28max+by+%28namespace%2C+app%2C+container%29+%28label_replace%28container_memory_working_set_bytes%7Bcontainer%21%3D%22%22%7D%2C+%22app%22%2C+%22%241%22%2C+%22pod%22%2C+%22%5E%28.%2B%29-%28%5Ba-z0-9%5D%2B-%5Ba-z0-9%5D%2B%7C%5B0-9%5D%2B%29%24%22%29%29+%2F+max+by+%28namespace%2C+app%2C+container%29+%28label_replace%28kube_pod_container_resource_requests%7Bresource%3D%22memory%22%7D%2C+%22app%22%2C+%22%241%22%2C+%22pod%22%2C+%22%5E%28.%2B%29-%28%5Ba-z0-9%5D%2B-%5Ba-z0-9%5D%2B%7C%5B0-9%5D%2B%29%24%22%29%29%29+%2A+100+%3E+100&g0.tab=1","fingerprint":"8b3e50ebf91a85bd"}],"groupLabels":{"alertname":"ContainerMemoryUsage","app":"internal-ingress-nginx-controller"},"commonLabels":{"alertname":"ContainerMemoryUsage","app":"internal-ingress-nginx-controller","cluster":"dev-gen-de02","container":"controller","environment":"Development","namespace":"internal-ingress-nginx","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"commonAnnotations":{"description":"Memory consumption of container controller of internal-ingress-nginx-controller in namespace internal-ingress-nginx\nis at 144.70% (> 100%) of the Memory requests for more than 10 minutes.","runbook":"https://company.atlassian.net/wiki/spaces/AR/pages/4075323867/Prometheus+alerts#Prometheusalerts-ContainerHighMemoryUsage","summary":"Container High Memory Usage"},"externalURL":"https://alertmanager.dev.company.io","version":"4","groupKey":"{}:{alertname=\"ContainerMemoryUsage\", app=\"internal-ingress-nginx-controller\"}","truncatedAlerts":0}`

	payload4 := `{"receiver":"default-route-receiver","status":"firing","alerts":[{"status":"firing","labels":{"alertname":"FailedEvictedPods","app":"discovery","cluster":"acc-gen-au01","environment":"Acceptance","namespace":"gloo-system","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"annotations":{"description":"Namespace gloo-system has 1 evicted pods named discovery in the last 10 minutes likely caused by Out of Memory issues.","runbook":"https://company.atlassian.net/wiki/spaces/AR/pages/4075323867/Prometheus+alerts#Prometheusalerts-FailedEvictedPods","summary":"Failed and Evicted Pods"},"startsAt":"2023-05-25T14:08:28.165Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":"https://prometheus.acc.company.io/graph?g0.expr=sum+by+%28namespace%2C+app%29+%28label_replace%28kube_pod_status_phase%7Bphase%3D%22Failed%22%7D+%3E+0+and+on+%28namespace%2C+pod%29+kube_pod_status_reason%7Breason%3D%22Evicted%22%7D+%3E+0%2C+%22app%22%2C+%22%241%22%2C+%22pod%22%2C+%22%5E%28.%2B%29-%28%5Ba-z0-9%5D%2B-%5Ba-z0-9%5D%2B%7C%5B0-9%5D%2B%29%22%29%29+%3E+0&g0.tab=1","fingerprint":"4d2c80e93b4b8dae"}],"groupLabels":{"alertname":"FailedEvictedPods","app":"discovery"},"commonLabels":{"alertname":"FailedEvictedPods","app":"discovery","cluster":"acc-gen-au01","environment":"Acceptance","namespace":"gloo-system","prometheus":"monitoring/kube-prometheus-stack-prometheus","severity":"warning"},"commonAnnotations":{"description":"Namespace gloo-system has 1 evicted pods named discovery in the last 10 minutes likely caused by Out of Memory issues.","runbook":"https://company.atlassian.net/wiki/spaces/AR/pages/4075323867/Prometheus+alerts#Prometheusalerts-FailedEvictedPods","summary":"Failed and Evicted Pods"},"externalURL":"https://alertmanager.acc.company.io","version":"4","groupKey":"{}:{alertname=\"FailedEvictedPods\", app=\"discovery\"}","truncatedAlerts":0}`

	var alertPayloads = []string{payload1, payload2, payload3, payload4}

	for _, value := range alertPayloads {
		// Make the POST request with bearer auth token SUPERTESTENVTOKEN

		//resp, err := http.Post("http://localhost:8080/api/webhooks/prometheus", "application/json", strings.NewReader(value))

		req, err := http.NewRequest("POST", url, strings.NewReader(value))
		if err != nil {
			fmt.Println("Error creating HTTP request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			return
		}
		defer resp.Body.Close()
	}

}
