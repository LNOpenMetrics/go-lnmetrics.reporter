package main

import (
	"runtime/debug"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/metrics"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/plugin"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

var metricsPlugin plugin.MetricsPlugin

func main() {
	defer func() {
		if x := recover(); x != nil {
			// recovering from a panic; x contains whatever was passed to panic()
			log.GetInstance().Errorf("run time panic: %v", x)
			log.GetInstance().Errorf("stacktrace %s", string(debug.Stack()))
		}
	}()

	metricsPlugin = plugin.MetricsPlugin{
		Metrics: make(map[int]metrics.Metric), Rpc: nil}

	plugin, err := plugin.ConfigureCLNPlugin[*plugin.MetricsPlugin](&metricsPlugin)
	if err != nil {
		panic(err)
	}

	// To set the time the following doc is followed
	// https://pkg.go.dev/github.com/robfig/cron?utm_source=godoc
	metricsPlugin.RegisterRecurrentEvt("@every 30m")

	metricsPlugin.Cron.Start()

	plugin.Start()
}
