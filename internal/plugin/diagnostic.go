package plugin

import (
	"encoding/json"
	"fmt"

	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

type RawLocalScoreRPC[T MetricsPluginState] struct{}

func NewRawLocalScoreRPC[T MetricsPluginState]() *RawLocalScoreRPC[T] {
	return &RawLocalScoreRPC[T]{}
}

func (instance *RawLocalScoreRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	metricOne, found := plugin.GetState().GetMetrics()[MetricOneID]
	var result map[string]any
	if !found {
		return nil, fmt.Errorf("Metric with id %d not found", 1)
	}

	// FIXME: improve the metric API to include the ToMap call
	resultStr, err := json.Marshal(metricOne)
	if err != nil {
		return nil, err
	}

	if err != json.Unmarshal(resultStr, &result) {
		return nil, err
	}

	return result, nil
}

// ForceUpdateRPC enable the force update command
type ForceUpdateRPC[T MetricsPluginState] struct{}

func (instance *ForceUpdateRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	for _, metric := range plugin.GetState().GetMetrics() {
		msg := Msg{
			cmd:    "plugin_rpc_method",
			params: map[string]any{"event": "on_force_update"},
		}
		plugin.GetState().CallUpdateOnMetric(metric, &msg)
	}
	response := map[string]any{
		"result": "force call update on all the metrics succeeded",
	}
	return response, nil
}
