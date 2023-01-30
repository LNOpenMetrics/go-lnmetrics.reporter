package plugin

import (
	"fmt"

	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

type RawLocalScoreRPC[T MetricsPluginState] struct{}

func NewRawLocalScoreRPC[T MetricsPluginState]() *RawLocalScoreRPC[T] {
	return &RawLocalScoreRPC[T]{}
}

func (instance *RawLocalScoreRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	metric, found := plugin.GetState().GetMetrics()[RawLocalScoreID]
	if !found {
		return nil, fmt.Errorf("Metric with id %d not found", 1)
	}
	return metric.ToMap()
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
