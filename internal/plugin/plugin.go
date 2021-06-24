package plugin

import (
	"errors"
	"fmt"
	"time"

	"github.com/niftynei/glightning/glightning"
	"github.com/robfig/cron/v3"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type MetricsPlugin struct {
	Plugin  *glightning.Plugin
	Metrics map[int]Metric
	Rpc     *glightning.Lightning
	Cron    *cron.Cron
}

func (plugin *MetricsPlugin) HendlerRPCMessage(event *glightning.RpcCommandEvent) error {
	command := event.Cmd
	switch command.MethodName {
	case "stop":
		// Share to all the metrics, so we need a global method that iterate over the metrics map
		params := make(map[string]interface{})
		params["timestamp"] = time.Now()
		msg := Msg{"stop", params}
		for _, metric := range plugin.Metrics {
			go plugin.callOnStopOnMetrics(metric, &msg)
		}
		log.GetInstance().Info("Close command received")
	default:
		return nil
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMetrics(id int, metric Metric) error {
	_, ok := plugin.Metrics[id]
	if ok {
		//TODO add more information in the error message
		log.GetInstance().Error(fmt.Sprintf("Metrics with is %d already registered."))
		return errors.New(fmt.Sprintf("Metrics with is %d already registered."))
	}
	plugin.Metrics[id] = metric
	return nil
}

func (plugin *MetricsPlugin) RegisterMethods() {
	method := NewMetricPlugin()
	rpcMethod := glightning.NewRpcMethod(method, "Show diagnostic node")
	rpcMethod.LongDesc = "Show the diagnostic data of the lightning network node"
	rpcMethod.Category = "metrics"
	plugin.Plugin.RegisterMethod(rpcMethod)
}

func (instance *MetricsPlugin) callUpdateOnMetric(metric Metric, msg *Msg) {
	metric.UpdateWithMsg(msg, instance.Rpc)
}

func (instance *MetricsPlugin) callOnStopOnMetrics(metric Metric, msg *Msg) {
	err := metric.OnClose(msg, instance.Rpc)
	if err != nil {
		log.GetInstance().Error(err)
	}
}

func (instance *MetricsPlugin) callUpdateOnMetricNoMsg(metric Metric) {
	log.GetInstance().Debug("Calling Update on metrics")
	err := metric.Update(instance.Rpc)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

func (instance *MetricsPlugin) RegisterRecurrentEvt(after string) {
	instance.Cron = cron.New()
	instance.Cron.AddFunc(after, func() {
		log.GetInstance().Debug("Calling recurrent function")
		for _, metric := range instance.Metrics {
			go instance.callUpdateOnMetricNoMsg(metric)
		}
	})
}
