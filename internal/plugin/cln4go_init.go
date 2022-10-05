package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	maker "github.com/LNOpenMetrics/go-lnmetrics.reporter/init/persistence"
	pluginDB "github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	sysinfo "github.com/elastic/go-sysinfo"
	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

func OnInit[T MetricsPluginState](plugin *cln4go.Plugin[T], request map[string]any) map[string]any {
	metricsPlugin := plugin.GetState()
	// FIXME: init rpc with the new interface

	lightningDir, _ := plugin.GetConf("lightning-dir")
	rpcFile, _ := plugin.GetConf("rpc-file")

	// FIXME(vincenzopalazzo): make possible that the user will choose the log level.
	if err := log.InitLogger(lightningDir.(string), "debug", false); err != nil {
		log.GetInstance().Error(err)
	}

	rpcPath := strings.Join([]string{lightningDir.(string), rpcFile.(string)}, "/")
	if err := metricsPlugin.NewClient(rpcPath); err != nil {
		panic(err)
	}
	metricsPath, err := maker.PrepareHomeDirectory(lightningDir.(string))
	if err != nil {
		panic(err)
	}

	dbPlugin, err := pluginDB.NewLevelDB(*metricsPath)
	if err != nil {
		panic(err)
	}
	metricsPlugin.SetStorage(dbPlugin)

	err = parseOptionsPlugin(plugin)
	if err != nil {
		panic(err)
	}
	// FIXME: Load all the metrics in the datatabase that are registered from the user
	metric, err := loadMetricIfExist(plugin, 1)

	if err != nil {
		panic(err)
	}

	if err := metricsPlugin.GetStorage().Migrate([]*string{metric.MetricName()}); err != nil {
		panic(err)
	}
	if err := metricsPlugin.RegisterMetrics(1, metric); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error received %s", err))
		panic(err)
	}
	// FIXME: After on init event c-lightning should be ready to accept request
	// from any plugin.
	metricsPlugin.RegisterOneTimeEvt("10s")
	return map[string]any{}
}

// This method include the code to parse the configuration options of the plugin.
func parseOptionsPlugin[T MetricsPluginState](plugin *cln4go.Plugin[T]) error {
	metricsPlugin := plugin.GetState()
	urlsAsString, found := plugin.GetOpt("lnmetrics-urls")
	urls := make([]string, 0)
	if found {
		urls = strings.FieldsFunc(urlsAsString.(string), func(r rune) bool {
			return r == ','
		})
	}

	noProxy, _ := plugin.GetOpt("lnmetrics-noproxy")
	proxy, _ := plugin.GetConf("proxy")
	if proxy != nil && !noProxy.(bool) {
		addr, _ := plugin.GetConf("address")
		port, _ := plugin.GetConf("port")
		server, err := graphql.NewWithProxy(urls, addr.(string), port.(uint64))
		if err != nil {
			return err
		}
		metricsPlugin.SetServer(server)
		metricsPlugin.SetProxy(true)
	} else {
		metricsPlugin.SetServer(graphql.New(urls))
		metricsPlugin.SetProxy(false)
	}
	// FIXME: Store the urls on db.
	return nil
}

func loadMetricIfExist[T MetricsPluginState](plugin *cln4go.Plugin[T], id int) (Metric, error) {
	metricName, found := MetricsSupported[id]
	if !found {
		log.GetInstance().Infof("Metric with id %d not supported", id)
		return nil, fmt.Errorf("Metric with id %d not supported", id)
	}
	log.GetInstance().Info(fmt.Sprintf("Loading metrics with id %d end name %s", id, metricName))

	switch id {
	case 1:
		return loadLastMetricOne(plugin)
	default:
		return nil, fmt.Errorf("Metric with is %d and name %s not supported", id, metricName)
	}
}

func loadLastMetricOne[T MetricsPluginState](plugin *cln4go.Plugin[T]) (*MetricOne, error) {
	metricsPlugin := plugin.GetState()
	metricDb, err := metricsPlugin.GetStorage().LoadLastMetricOne()
	if err != nil {
		log.GetInstance().Info("No metrics available yet")
		log.GetInstance().Debugf("Error received %s", err)
		sys, err := sysinfo.Host()
		if err != nil {
			log.GetInstance().Errorf("Error during get the system information, error description %s", err)
			return nil, err
		}
		one := NewMetricOne("", sys.Info(), metricsPlugin.GetStorage())
		return one, nil
	}
	log.GetInstance().Info("Metrics One available on DB, loading it.")
	var metric MetricOne
	/// FIXME: try to use the plugin encoder here?
	err = json.Unmarshal([]byte(*metricDb), &metric)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error received %s", err))
		return nil, err
	}
	metric.Storage = metricsPlugin.GetStorage()
	return &metric, nil
}
