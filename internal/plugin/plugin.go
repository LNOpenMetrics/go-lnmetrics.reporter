package plugin

import (
	"fmt"
	"time"

	cron "github.com/robfig/cron/v3"
	cln4go "github.com/vincenzopalazzo/cln4go/client"
	"github.com/vincenzopalazzo/glightning/glightning"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

// FIXME: move this to a generics to set the Client
// in this way we could support different implementation
type MetricsPlugin struct {
	Plugin    *glightning.Plugin
	Metrics   map[int]Metric
	Rpc       *cln4go.UnixRPC
	Cron      *cron.Cron
	Server    *graphql.Client
	Storage   db.PluginDatabase
	WithProxy bool
}

func (self *MetricsPlugin) GetRpc() *cln4go.UnixRPC {
	return self.Rpc
}

func (self *MetricsPlugin) NewClient(path string) {
	self.Rpc, _ = cln4go.NewUnix(path)
}

func (self *MetricsPlugin) SetStorage(storage db.PluginDatabase) {
	self.Storage = storage
}

func (self *MetricsPlugin) GetStorage() db.PluginDatabase {
	return self.Storage
}

func (self *MetricsPlugin) SetProxy(withProxy bool) {
	self.WithProxy = withProxy
}

func (self *MetricsPlugin) IsWithProxy() bool {
	return self.WithProxy
}

func (self *MetricsPlugin) SetServer(server *graphql.Client) {
	self.Server = server
}

func (self *MetricsPlugin) GetServer() *graphql.Client {
	return self.Server
}

// FIXME: switch to the shutdown notification
func (plugin *MetricsPlugin) HandlerRPMMessage(event *glightning.RpcCommandEvent) error {
	command := event.Cmd
	switch command.MethodName {
	case "stop":
		// Share to all the metrics, so we need a global method that iterate over the metrics map
		params := make(map[string]any)
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
		log.GetInstance().Errorf("Metrics with is %d already registered.", id)
		return fmt.Errorf("Metrics with is %d already registered.", id)
	}
	plugin.Metrics[id] = metric
	return nil
}

func (plugin *MetricsPlugin) RegisterMethods() error {
	method := NewMetricPlugin(plugin)
	rpcMethod := glightning.NewRpcMethod(method, "Show diagnostic node")
	rpcMethod.LongDesc = "Show the metric one data of the lightning network node. An example metric_one start=last"
	rpcMethod.Category = "lnmetrics"
	if err := plugin.Plugin.RegisterMethod(rpcMethod); err != nil {
		return err
	}

	infoMethod := NewPluginRpcMethod(plugin)
	infoRpcMethod := glightning.NewRpcMethod(infoMethod, "Show go-lnmetrics.reporter info")
	infoRpcMethod.Category = "lnmetrics"
	infoRpcMethod.LongDesc = "Return the info od the env where the plugin is running. An example is \"lnmetrics-info"
	if err := plugin.Plugin.RegisterMethod(infoRpcMethod); err != nil {
		return err
	}

	cacheMethod := NewCleanCacheRPC(plugin)
	cacheRPCMethod := glightning.NewRpcMethod(cacheMethod, "Clean all the lnmetrics cache")
	cacheRPCMethod.Category = "lnmetrics"
	cacheRPCMethod.LongDesc = "Clean the cache made by the plugin during the time. An example is \"lnmetrics-cache clean"
	if err := plugin.Plugin.RegisterMethod(cacheRPCMethod); err != nil {
		return err
	}

	forceUpdate := NewForceUpdateRPC(plugin)
	forceUpdateRPC := glightning.NewRpcMethod(forceUpdate, "call the update on all the plugin")
	forceUpdateRPC.Category = "lnmetrics"
	if err := plugin.Plugin.RegisterMethod(forceUpdateRPC); err != nil {
		return err
	}
	return nil
}

// nolint
func (plugin *MetricsPlugin) callUpdateOnMetric(metric Metric, msg *Msg) {
	if err := metric.UpdateWithMsg(msg, plugin.Rpc); err != nil {
		log.GetInstance().Errorf("Error during update metrics event: %s", err)
	}
}

// callOnStopOnMetrics Call on stop operation on the node when the caller are shutdown itself.
func (plugin *MetricsPlugin) callOnStopOnMetrics(metric Metric, msg *Msg) {
	err := metric.OnStop(msg, plugin.GetRpc())
	if err != nil {
		log.GetInstance().Error(err)
	}
}

// callUpdateOnMetricNoMsg Update the metrics without any information received by the caller
func (plugin *MetricsPlugin) callUpdateOnMetricNoMsg(metric Metric) {
	log.GetInstance().Debug("Calling Update on metrics")
	err := metric.Update(plugin.GetRpc())
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

func (plugin *MetricsPlugin) updateAndUploadMetric(metric Metric) {
	log.GetInstance().Info("Calling update and upload metric")
	plugin.callUpdateOnMetricNoMsg(metric)
	if err := metric.UploadOnRepo(plugin.Server, plugin.GetRpc()); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

// RegisterRecurrentEvt Register internal recurrent methods
func (plugin *MetricsPlugin) RegisterRecurrentEvt(after string) {
	log.GetInstance().Info(fmt.Sprintf("Register recurrent event each %s", after))
	plugin.Cron = cron.New()
	// FIXME: Discover what is the first value
	_, err := plugin.Cron.AddFunc(after, func() {
		log.GetInstance().Info("Update and Uploading metrics")
		for _, metric := range plugin.Metrics {
			go plugin.updateAndUploadMetric(metric)
		}
	})
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during registering recurrent event: %s", err))
	}
}

func (plugin *MetricsPlugin) RegisterOneTimeEvt(after string) {
	log.GetInstance().Info(fmt.Sprintf("Register one time event after %s", after))
	duration, err := time.ParseDuration(after)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error in the on time evt: %s", err))
		return
	}
	time.AfterFunc(duration, func() {
		log.GetInstance().Debug("Calling on time function function")
		// TODO: Should C-Lightning send a on init event like notification?
		for _, metric := range plugin.Metrics {
			go func(instance *MetricsPlugin, metric Metric) {
				err := metric.OnInit(instance.GetRpc())
				if err != nil {
					log.GetInstance().Error(fmt.Sprintf("Error during on init call: %s", err))
				}

				// Init on server.
				if err := metric.InitOnRepo(instance.Server, instance.GetRpc()); err != nil {
					log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
				}

			}(plugin, metric)
		}
	})
}
