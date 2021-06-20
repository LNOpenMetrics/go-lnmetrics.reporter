package plugin

import (
	"errors"
	"fmt"
	"github.com/niftynei/glightning/jrpc2"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
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
	switch instance.MetricId {
	case 1:
		return db.GetInstance().GetValue("metric_one")
	default:
		return nil, errors.New(fmt.Sprintf("ID metrics unknown"))
	}
}
