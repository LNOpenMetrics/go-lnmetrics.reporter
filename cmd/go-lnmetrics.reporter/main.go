package main

import (
	"fmt"

	metrics "github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/plugin"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	"github.com/vincenzopalazzo/glightning/glightning"
)

var metricsPlugin metrics.MetricsPlugin

func main() {

	// FIXME: I can remove the Plugin?
	metricsPlugin = metrics.MetricsPlugin{Plugin: nil,
		Metrics: make(map[int]metrics.Metric), Rpc: nil}

	plugin, err := metrics.ConfigureCLNPlugin[metrics.MetricsPlugin](&metricsPlugin)
	if err != nil {
		panic(err)
	}

	hook := &glightning.Hooks{RpcCommand: OnRpcCommand}
	if err := plugin.RegisterHooks(hook); err != nil {
		panic(err)
	}

	if err := metricsPlugin.RegisterMethods(); err != nil {
		panic(err)
	}

	// To set the time the following doc is followed
	// https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc
	metricsPlugin.RegisterRecurrentEvt("@every 30m")

	metricsPlugin.Cron.Start()

	plugin.Start()
}

func OnRpcCommand(event *glightning.RpcCommandEvent) (*glightning.RpcCommandResponse, error) {
	if err := metricsPlugin.HendlerRPCMessage(event); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during a hook handler: %s", err))
	}
	return event.Continue(), nil
}
