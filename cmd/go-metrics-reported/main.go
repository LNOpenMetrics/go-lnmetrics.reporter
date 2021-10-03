package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	maker "github.com/OpenLNMetrics/go-metrics-reported/init/persistence"
	metrics "github.com/OpenLNMetrics/go-metrics-reported/internal/plugin"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/graphql"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/vincenzopalazzo/glightning/glightning"
)

var metricsPlugin metrics.MetricsPlugin

func main() {
	plugin := glightning.NewPlugin(onInit)

	metricsPlugin = metrics.MetricsPlugin{Plugin: plugin,
		Metrics: make(map[int]metrics.Metric), Rpc: nil}

	hook := &glightning.Hooks{RpcCommand: OnRpcCommand}
	if err := plugin.RegisterHooks(hook); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		panic(err)
	}
	if err := metricsPlugin.RegisterMethods(); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		panic(err)
	}

	// To set the time the following doc is followed
	// https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc
	metricsPlugin.RegisterRecurrentEvt("@every 30m")

	metricsPlugin.Cron.Start()

	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.GetInstance().Error(err)
		panic(err)
	}
}

func onInit(plugin *glightning.Plugin,
	options map[string]glightning.Option, config *glightning.Config) {
	metricsPlugin.Rpc = glightning.NewLightning()

	metricsPlugin.Rpc.StartUp(config.RpcFile, config.LightningDir)
	metricsPath, err := maker.PrepareHomeDirectory(config.LightningDir)
	if err != nil {
		log.GetInstance().Error(err)
		panic(err)
	}
	if err := db.GetInstance().InitDB(*metricsPath); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		panic(err)
	}
	err = parseOptionsPlugin(options)
	if err != nil {
		log.GetInstance().Error(err)
		panic(err)
	}
	//TODO: Load all the metrics in the datatabase that are registered from
	// the user
	metric, err := loadMetricIfExist(1)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error received %s", err))
		panic(err)
	}

	if err := metricsPlugin.RegisterMetrics(1, metric); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error received %s", err))
		panic(err)
	}

	metricsPlugin.RegisterOneTimeEvt("10s")
}

func OnRpcCommand(event *glightning.RpcCommandEvent) (*glightning.RpcCommandResponse, error) {
	if err := metricsPlugin.HendlerRPCMessage(event); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during a hook handler: %s", err))
	}
	return event.Continue(), nil
}

// This method include the code to parse the configuration options of the plugin.
func parseOptionsPlugin(options map[string]glightning.Option) error {
	urlsAsString, found := options["lnmetrics-urls"]
	urls := make([]string, 0)
	if found {
		urls = strings.FieldsFunc(urlsAsString.GetValue().(string), func(r rune) bool {
			return r == ','
		})
	}
	metricsPlugin.Server = graphql.New(urls)
	// FIXME: Store the urls on db.
	return nil
}

//FIXME: Improve quality of Go style here
func loadMetricIfExist(id int) (metrics.Metric, error) {
	metricName, found := metrics.MetricsSupported[id]
	if !found {
		log.GetInstance().Info(fmt.Sprintf("Metric with id %d not supported", id))
		return nil, fmt.Errorf("Metric with id %d not supported", id)
	}
	log.GetInstance().Info(fmt.Sprintf("Loading metrics with id %d end name %s", id, metricName))
	metricDb, err := db.GetInstance().GetValue(metricName)
	if err != nil {
		log.GetInstance().Info("No metrics available yet")
		log.GetInstance().Debug(fmt.Sprintf("Error received %s", err))
		sys, err := sysinfo.Host()
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error during get the system information, error description %s", err))
			return nil, err
		}
		switch id {
		case 1:
			one := metrics.NewMetricOne("", sys.Info())
			return one, nil

		default:
			return nil, fmt.Errorf("Metric with id %d not supported", id)
		}
	}
	log.GetInstance().Info("Metrics available on DB, loading them.")
	switch id {
	case 1:
		var metric metrics.MetricOne
		err = json.Unmarshal([]byte(metricDb), &metric)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error received %s", err))
			return nil, err
		}
		return &metric, nil
	default:
		return nil, fmt.Errorf("Metric with id %d not supported", id)
	}
}
