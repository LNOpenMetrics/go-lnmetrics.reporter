package main

import (
	"os"

	maker "github.com/OpenLNMetrics/go-metrics-reported/init/persistence"
	metrics "github.com/OpenLNMetrics/go-metrics-reported/internal/plugin"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"

	"github.com/niftynei/glightning/glightning"
)

var metricsPlugin metrics.MetricsPlugin

func main() {
	log.GetInstance().Info("Init plugin")
	plugin := glightning.NewPlugin(onInit)

	metricsPlugin = metrics.MetricsPlugin{plugin}

	plugin.RegisterHooks(&glightning.Hooks{
		RpcCommand: OnRpcCommand,
	})

	metricsPlugin.RegisterMethods()

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
	metricsPath, err := maker.PrepareHomeDirectory(config.LightningDir)
	if err != nil {
		log.GetInstance().Error(err)
		panic(err)
	}
	db.GetInstance().InitDB(*metricsPath)
}

func OnRpcCommand(event *glightning.RpcCommandEvent) (*glightning.RpcCommandResponse, error) {
	method := event.Cmd.MethodName
	log.GetInstance().Debug("hook throws by the following rpc command" + method)
	metricsPlugin.HendlerRPCMessage(event)
	return event.Continue(), nil
}
