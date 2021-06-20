package plugin

import (
	"errors"
	"sync"
	"time"

	"github.com/niftynei/glightning/glightning"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type MetricsPlugin struct {
	Plugin  *glightning.Plugin
	Metrics map[int]*Metric
	Rpc     *glightning.Lightning
}

func (plugin *MetricsPlugin) HendlerRPCMessage(event *glightning.RpcCommandEvent) error {
	command := event.Cmd
	method, err := command.Get()
	if err != nil {
		return err
	}
	switch method.(type) {
	case glightning.CloseRequest:
		// Share to all the metrics, so we need a global method that iterate over the metrics map
		params := make(map[string]interface{})
		params["timestamp"] = time.Now()
		msg := Msg{"close", params}
		var courutinesWait sync.WaitGroup
		courutinesWait.Add(len(plugin.Metrics))
		for _, metric := range plugin.Metrics {
			plugin.callUpdateOnMetric(metric, &msg, &courutinesWait)
		}
		courutinesWait.Wait()
		log.GetInstance().Debug("Close command received")
	default:
		var courutinesWait sync.WaitGroup
		courutinesWait.Add(len(plugin.Metrics))
		for _, metric := range plugin.Metrics {
			plugin.callUpdateOnMetricNoMsg(metric, &courutinesWait)
		}
		courutinesWait.Wait()
		log.GetInstance().Debug("The node is up and runnning update the info")
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMetrics(id int, metric *Metric) error {
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

func (instance *MetricsPlugin) callUpdateOnMetric(metric *Metric, msg *Msg,
	corutine *sync.WaitGroup) {
	defer corutine.Done()
	(*metric).UpdateWithMsg(msg, instance.Rpc)
}

func (instance *MetricsPlugin) callUpdateOnMetricNoMsg(metric *Metric,
	corutine *sync.WaitGroup) {
	defer corutine.Done()
	(*metric).Update(instance.Rpc)
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
