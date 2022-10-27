package main

import (
	metrics "github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/plugin"
)

var metricsPlugin metrics.MetricsPlugin

func main() {

	// FIXME: I can remove the Plugin?
	metricsPlugin = metrics.MetricsPlugin{Plugin: nil,
		Metrics: make(map[int]metrics.Metric), Rpc: nil}

	plugin, err := metrics.ConfigureCLNPlugin[*metrics.MetricsPlugin](&metricsPlugin)
	if err != nil {
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
