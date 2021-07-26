package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	maker "github.com/OpenLNMetrics/go-metrics-reported/init/persistence"
	metrics "github.com/OpenLNMetrics/go-metrics-reported/internal/plugin"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/niftynei/glightning/glightning"
)

var metricsPlugin metrics.MetricsPlugin

func main() {
	log.GetInstance().Info("Init plugin")
	plugin := glightning.NewPlugin(onInit)

	metricsPlugin = metrics.MetricsPlugin{Plugin: plugin,
		Metrics: make(map[int]metrics.Metric), Rpc: nil}

	plugin.RegisterHooks(&glightning.Hooks{
		RpcCommand: OnRpcCommand,
	})

	metricsPlugin.RegisterMethods()

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
	log.GetInstance().Debug("Options node have the following paramameters")
	log.GetInstance().Debug(options)
	log.GetInstance().Debug("Node with the following configuration")
	log.GetInstance().Debug(config)
	metricsPlugin.Rpc = glightning.NewLightning()

	metricsPlugin.Rpc.StartUp(config.RpcFile, config.LightningDir)
	metricsPath, err := maker.PrepareHomeDirectory(config.LightningDir)
	if err != nil {
		log.GetInstance().Error(err)
		panic(err)
	}
	db.GetInstance().InitDB(*metricsPath)

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
	metricsPlugin.HendlerRPCMessage(event)
	return event.Continue(), nil
}

//FIXME: Improve quality of Go style here
func loadMetricIfExist(id int) (metrics.Metric, error) {
	metricName, ok := metrics.MetricsSupported[id]
	if ok == false {
		log.GetInstance().Info(fmt.Sprintf("Metric with id %d not supported", id))
		return nil, errors.New(fmt.Sprintf("Metric with id %s not supported", id))
	}
	log.GetInstance().Info(fmt.Sprintf("Loading metrics with id %s end name", id, metricName))
	metricDb, err := db.GetInstance().GetValue(metricName)
	log.GetInstance().Info("value on db us " + metricDb)
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
			return nil, errors.New(fmt.Sprintf("Metric with id %d not supported", id))
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
		return nil, errors.New(fmt.Sprintf("Metric with id %d not supported", id))
	}
}
