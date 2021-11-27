package plugin

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/vincenzopalazzo/glightning/glightning"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

type MetricsPlugin struct {
	Plugin  *glightning.Plugin
	Metrics map[int]Metric
	Rpc     *glightning.Lightning
	Cron    *cron.Cron
	Server  *graphql.Client
	Storage db.PluginDatabase
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
		plugin.Cron.Stop()
		log.GetInstance().Info("Close command received")
	default:
		return nil
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMetrics(id int, metric Metric) error {
	_, ok := plugin.Metrics[id]
	if ok {
		log.GetInstance().Error(fmt.Sprintf("Metrics with is %d already registered.", id))
		return fmt.Errorf("Metrics with is %d already registered.", id)
	}
	plugin.Metrics[id] = metric
	return nil
}

func (plugin *MetricsPlugin) RegisterMethods() error {
	method := NewMetricPlugin(plugin)
	rpcMethod := glightning.NewRpcMethod(method, "Show diagnostic node")
	rpcMethod.LongDesc = "Show the diagnostic data of the lightning network node"
	rpcMethod.Category = "metrics"
	if err := plugin.Plugin.RegisterMethod(rpcMethod); err != nil {
		return err
	}

	infoMethod := NewPluginRpcMethod()
	infoRpcMethod := glightning.NewRpcMethod(infoMethod, "Show go-lnmetrics.reporter info")
	infoRpcMethod.Category = "metrics"
	infoRpcMethod.LongDesc = "Return a map where the key is the id of the method and the value is the payload of the metric. The metrics_id is a string that conatins the id divided by a comma. An example is \"diagnostic \"1,2,3\"\""
	if err := plugin.Plugin.RegisterMethod(infoRpcMethod); err != nil {
		return err
	}

	return nil
}

//nolint
func (instance *MetricsPlugin) callUpdateOnMetric(metric Metric, msg *Msg) {
	if err := metric.UpdateWithMsg(msg, instance.Rpc); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during update metrics event: %s", err))
	}
}

// Call on stop operation on the node when the caller are shoutdown it self.
func (instance *MetricsPlugin) callOnStopOnMetrics(metric Metric, msg *Msg) {
	err := metric.OnClose(msg, instance.Rpc)
	if err != nil {
		log.GetInstance().Error(err)
	}
}

// Update the metrics without any information received by the caller
func (instance *MetricsPlugin) callUpdateOnMetricNoMsg(metric Metric) {
	log.GetInstance().Debug("Calling Update on metrics")
	err := metric.Update(instance.Rpc)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

func (instance *MetricsPlugin) updateAndUploadMetric(metric Metric) {
	log.GetInstance().Info("Calling update and upload metric")
	instance.callUpdateOnMetricNoMsg(metric)
	if err := metric.UploadOnRepo(instance.Server, instance.Rpc); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

// Register internal recurrent methods
func (instance *MetricsPlugin) RegisterRecurrentEvt(after string) {
	instance.Cron = cron.New()
	// FIXME: Discover what is the first value
	_, err := instance.Cron.AddFunc(after, func() {
		log.GetInstance().Info("Update and Uploading metrics")
		for _, metric := range instance.Metrics {
			go instance.updateAndUploadMetric(metric)
		}
	})
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during registering recurrent event: %s", err))
	}
}

func (instance *MetricsPlugin) RegisterOneTimeEvt(after string) {
	duration, err := time.ParseDuration(after)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error in the on time evt: %s", err))
		return
	}
	time.AfterFunc(duration, func() {
		log.GetInstance().Debug("Calling on time function function")
		// TODO: Should C-Lightning send a on init event like notification?
		for _, metric := range instance.Metrics {
			go func(instance *MetricsPlugin, metric Metric) {
				err := metric.OnInit(instance.Rpc)
				if err != nil {
					log.GetInstance().Error(fmt.Sprintf("Error during on init call: %s", err))
				}

				// Init on server.
				if err := metric.InitOnRepo(instance.Server, instance.Rpc); err != nil {
					log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
				}

			}(instance, metric)
		}
	})
}
