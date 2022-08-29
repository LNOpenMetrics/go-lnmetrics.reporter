// Package plugin Provide dev call a runtime to perform integration testing
package plugin

import (
	"encoding/json"
	"github.com/vincenzopalazzo/glightning/jrpc2"
)

// PluginUploadMetricDev is a dev rpc that will be used
// for dev plus
type PluginUploadMetricDev struct {
	jsonPayload   *string        `json:"-"`
	metricsPlugin *MetricsPlugin `json:"-"`
}

func NewPluginDev(metricsPlugin *MetricsPlugin) *PluginUploadMetricDev {
	return &PluginUploadMetricDev{
		metricsPlugin: metricsPlugin,
	}
}

func (instance PluginUploadMetricDev) Name() string {
	return "lnmetrics-dev-upload"
}

func (instance *PluginUploadMetricDev) New() any {
	return instance
}

func (instance *PluginUploadMetricDev) Call() (jrpc2.Result, error) {
	metrics := MetricOne{}
	if err := json.Unmarshal([]byte(*instance.jsonPayload), &metrics); err != nil {
		return nil, err
	}

	if err := instance.metricsPlugin.RegisterMetrics(1, &metrics); err != nil {
		return nil, err
	}

	response := struct {
		metrics *MetricOne
		result  any
	}{
		metrics: &metrics,
		result:  metrics.DevServerUpload(instance.metricsPlugin.Server, instance.metricsPlugin.Rpc),
	}
	return response, nil
}
