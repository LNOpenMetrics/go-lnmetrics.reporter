package plugin

import (
	"errors"
	"fmt"
	"github.com/niftynei/glightning/glightning"
	"sync"
	"time"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type MetricsPlugin struct {
	Plugin  *glightning.Plugin
	Metrics map[int]Metric
	Rpc     *glightning.Lightning
}

func (plugin *MetricsPlugin) HendlerRPCMessage(event *glightning.RpcCommandEvent) error {
	command := event.Cmd
	switch command.MethodName {
	case "stop":
		// Share to all the metrics, so we need a global method that iterate over the metrics map
		params := make(map[string]interface{})
		params["timestamp"] = time.Now()
		msg := Msg{"stop", params}
		var courutinesWait sync.WaitGroup
		courutinesWait.Add(len(plugin.Metrics))
		for _, metric := range plugin.Metrics {
			plugin.callOnStopOnMetrics(metric, &msg, &courutinesWait)
		}
		courutinesWait.Wait()
		log.GetInstance().Debug("Close command received")
	default:
		return nil
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMetrics(id int, metric Metric) error {
	_, ok := plugin.Metrics[id]
	if ok {
		//TODO add more information in the error message
		return errors.New("Metrics already registered")
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

func (instance *MetricsPlugin) callUpdateOnMetric(metric Metric, msg *Msg,
	corutine *sync.WaitGroup) {
	defer corutine.Done()
	metric.UpdateWithMsg(msg, instance.Rpc)
}

func (instance *MetricsPlugin) callOnStopOnMetrics(metric Metric, msg *Msg,
	corutine *sync.WaitGroup) {
	defer corutine.Done()
	err := metric.OnClose(msg, instance.Rpc)
	if err != nil {
		log.GetInstance().Error(err)
	}
}

func (instance *MetricsPlugin) callUpdateOnMetricNoMsg(metric Metric,
	corutine *sync.WaitGroup) {
	defer corutine.Done()
	err := metric.Update(instance.Rpc)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

func (instance *MetricsPlugin) RegisterRecurrentEvt(after time.Duration) {
	time.AfterFunc(after, func() {
		log.GetInstance().Debug("Recurrent event called")
		var courutinesWait sync.WaitGroup
		courutinesWait.Add(len(instance.Metrics))
		for _, metric := range instance.Metrics {
			instance.callUpdateOnMetricNoMsg(metric, &courutinesWait)
		}
		courutinesWait.Wait()
	})
}
