package plugin

import (
	"fmt"
	"time"

	cron "github.com/robfig/cron/v3"
	cln4go "github.com/vincenzopalazzo/cln4go/client"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/metrics"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/json"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/trace"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

// FIXME: move this to a generics to set the Client
// in this way we could support different implementation
type MetricsPlugin struct {
	Metrics   map[int]metrics.Metric
	Rpc       *cln4go.UnixRPC
	Cron      *cron.Cron
	Server    *graphql.Client
	Storage   db.PluginDatabase
	WithProxy bool
}

func (self *MetricsPlugin) GetRpc() *cln4go.UnixRPC {
	return self.Rpc
}

func (self *MetricsPlugin) NewClient(path string) error {
	rpc, err := cln4go.NewUnix(path)
	if err != nil {
		return err
	}
	rpc.SetEncoder(&json.FastJSON{})
	rpc.SetTracer(&trace.Tracer{})
	self.Rpc = rpc
	return nil
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

func (plugin *MetricsPlugin) RegisterMetrics(id int, metric metrics.Metric) error {
	_, ok := plugin.Metrics[id]
	if ok {
		log.GetInstance().Errorf("Metrics with is %d already registered.", id)
		return fmt.Errorf("Metrics with is %d already registered.", id)
	}
	plugin.Metrics[id] = metric
	return nil
}

func (self *MetricsPlugin) GetMetrics() map[int]metrics.Metric {
	return self.Metrics
}

func (self *MetricsPlugin) GetCron() *cron.Cron {
	return self.Cron
}

// nolint
func (plugin *MetricsPlugin) CallUpdateOnMetric(metric metrics.Metric, msg *metrics.Msg) {
	if err := metric.UpdateWithMsg(msg, plugin.Rpc); err != nil {
		log.GetInstance().Errorf("Error during update metrics event: %s", err)
	}
}

// callOnStopOnMetrics Call on stop operation on the node when the caller are shutdown itself.
func (plugin *MetricsPlugin) CallOnStopOnMetrics(metric metrics.Metric, msg *metrics.Msg) {
	err := metric.OnStop(msg, plugin.GetRpc())
	if err != nil {
		log.GetInstance().Error(err)
	}
}

// callUpdateOnMetricNoMsg Update the metrics without any information received by the caller
func (plugin *MetricsPlugin) callUpdateOnMetricNoMsg(metric metrics.Metric) {
	log.GetInstance().Debug("Calling Update on metrics")
	err := metric.Update(plugin.GetRpc())
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s", err))
	}
}

func (plugin *MetricsPlugin) updateAndUploadMetric(metric metrics.Metric) {
	log.GetInstance().Info("Calling update and upload metric")
	plugin.callUpdateOnMetricNoMsg(metric)
	if err := metric.UploadOnRepo(plugin.Server, plugin.GetRpc()); err != nil {
		log.GetInstance().Errorf("Error %s", err)
	}
	log.GetInstance().Info("Metrics shipped!")
}

// RegisterRecurrentEvt Register internal recurrent methods
func (plugin *MetricsPlugin) RegisterRecurrentEvt(after string) {
	log.GetInstance().Info(fmt.Sprintf("Register recurrent event each %s", after))
	plugin.Cron = cron.New()
	// FIXME: Discover what is the first value
	_, err := plugin.Cron.AddFunc(after, func() {
		log.GetInstance().Info("Update and Uploading metrics")
		for _, metric := range plugin.Metrics {
			// FIXME: we can do this in parallel, but not now :)
			plugin.updateAndUploadMetric(metric)
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
		// FIXME: Should C-Lightning send a on init event like notification?
		for _, metric := range plugin.Metrics {
			go func(instance *MetricsPlugin, metric metrics.Metric) {
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
