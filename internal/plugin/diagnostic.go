package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/metrics"
	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

type RawLocalScoreRPC[T MetricsPluginState] struct{}

func NewRawLocalScoreRPC[T MetricsPluginState]() *RawLocalScoreRPC[T] {
	return &RawLocalScoreRPC[T]{}
}

func (instance *RawLocalScoreRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	metricStr, err := plugin.GetState().GetStorage().LoadLastMetricOne()
	if err != nil {
		return nil, fmt.Errorf("Metric with id %d not found", 1)
	}
	// FIXME get the encoder from a plugin with GetEncoder and decode the string into a map
	var metric map[string]any
	if err := json.Unmarshal([]byte(*metricStr), &metric); err != nil {
		return nil, err
	}
	return metric, nil
}

// ForceUpdateRPC enable the force update command
type ForceUpdateRPC[T MetricsPluginState] struct{}

func (instance *ForceUpdateRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	for _, metric := range plugin.GetState().GetMetrics() {
		msg := metrics.NewMsg("dev-collect", map[string]any{"event": "on_force_update"})
		plugin.GetState().CallUpdateOnMetric(metric, &msg)
	}
	response := map[string]any{
		"result": "force call update on all the metrics succeeded",
	}
	return response, nil
}
