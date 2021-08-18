package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/niftynei/glightning/jrpc2"
	"strconv"
	"strings"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type DiagnosticRpcMethod struct {
	MetricId  int    `json:"metric_id"`
	MetricsId string `json:"metrics_id"`
}

type diagnosticRpcModel struct {
	Metrics map[string]interface{} `json:"metrics"`
}

func (rpc *DiagnosticRpcMethod) Name() string {
	return "diagnostic"
}

func NewMetricPlugin() *DiagnosticRpcMethod {
	return &DiagnosticRpcMethod{MetricId: -1, MetricsId: ""}
}

func (rpc *DiagnosticRpcMethod) New() interface{} {
	return NewMetricPlugin()
}

func (instance *DiagnosticRpcMethod) Call() (jrpc2.Result, error) {
	metricsRequired, err := instance.parsingMetrics()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
		return nil, err
	}
	if len(metricsRequired) == 0 && instance.MetricId > 0 {
		metricsRequired = append(metricsRequired, instance.MetricId)
	}
	model := diagnosticRpcModel{Metrics: make(map[string]interface{})}
	for _, metricId := range metricsRequired {
		key, found := MetricsSupported[metricId]
		if !found {
			return nil, errors.New(fmt.Sprintf("ID metrics %d unknown", metricId))
		}
		result, err := db.GetInstance().GetValue(key)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("DB error for the key %s", key))
			log.GetInstance().Error(fmt.Sprintf("Error is: %s", err))
			return nil, errors.New(fmt.Sprintf("DB error for the metric %s with following motivation %s", key, err))
		}

		var metric interface{}
		err = json.Unmarshal([]byte(result), &metric)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			return nil, err
		}
		model.Metrics[key] = metric
	}
	return model, nil
}

func (instance *DiagnosticRpcMethod) parsingMetrics() ([]int, error) {
	metrics := make([]int, 0)
	tokens := strings.Split(instance.MetricsId, ",")
	for _, value := range tokens {
		trimVal := strings.Trim(value, " ")
		if trimVal == "" {
			continue
		}
		val, err := strconv.Atoi(trimVal)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error %s", err))
			return nil, err
		}
		metrics = append(metrics, val)
	}
	return metrics, nil
}
