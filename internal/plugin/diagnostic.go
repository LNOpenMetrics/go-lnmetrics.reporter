package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/niftynei/glightning/jrpc2"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type DiagnosticRpcMethod struct {
	MetricId int `json:"metric_id"`
}

func (rpc *DiagnosticRpcMethod) Name() string {
	return "diagnostic"
}

func NewMetricPlugin() *DiagnosticRpcMethod {
	return &DiagnosticRpcMethod{}
}

func (rpc *DiagnosticRpcMethod) New() interface{} {
	return NewMetricPlugin()
}

func (instance *DiagnosticRpcMethod) Call() (jrpc2.Result, error) {
	key, found := MetricsSupported[instance.MetricId]
	if !found {
		return nil, errors.New(fmt.Sprintf("ID metrics unknown"))
	}
	result, err := db.GetInstance().GetValue(key)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("DB error for the key %s", key))
		log.GetInstance().Error(fmt.Sprintf("Error is: %s", err))
		return nil, errors.New(fmt.Sprintf("DB error for the metric %s with following motivation %s", key, err))
	}

	log.GetInstance().Debug(fmt.Sprintf("Result in the map %s", result))
	var metricOne interface{}
	err = json.Unmarshal([]byte(result), &metricOne)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, err
	}
	return metricOne, nil
}
