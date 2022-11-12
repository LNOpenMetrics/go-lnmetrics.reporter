package plugin

import (
	"encoding/json"
	"fmt"

	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

type MetricOneRpcMethod[T MetricsPluginState] struct{}

func NewMetricPlugin[T MetricsPluginState]() *MetricOneRpcMethod[T] {
	return &MetricOneRpcMethod[T]{}
}

func (instance *MetricOneRpcMethod[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	// FIXME: take variable from the payload.
	metricOne, found := plugin.GetState().GetMetrics()[1]
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

	startPeriod, startFound := payload["start"]
	//endPeriod, endFound := payload["end"]

	if !startFound {
		return nil, fmt.Errorf("methor arg missing: need to specify the start period")
	}

	if startPeriod.(string) == "now" {
		// FIXME: encode the result inside a map
		return result, nil
	}

	if startPeriod.(string) == "last" {
		jsonValue, err := plugin.GetState().GetStorage().LoadLastMetricOne()
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(*jsonValue), &metricOne); err != nil {
			return nil, err
		}

		// FIXME: encode the result inside a map
		return result, nil
	}

	return nil, fmt.Errorf("We don't support the filter operation right now")
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
