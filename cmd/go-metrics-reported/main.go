package main

import (
	"fmt"
	"os"

	maker "github.com/OpenLNMetrics/go-metrics-reported/init/persistence"
	metrics "github.com/OpenLNMetrics/go-metrics-reported/internal/plugin"
	log "github.com/OpenLNMetrics/go-metrics-reported/pkg/log"

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
		fmt.Println(err)
	}
}

func onInit(plugin *glightning.Plugin,
	options map[string]glightning.Option, config *glightning.Config) {
	log.GetInstance().Debug("Options node have the following paramameters")
	log.GetInstance().Debug(options)
	log.GetInstance().Debug("Node with the following configuration")
	log.GetInstance().Debug(config)
	err := maker.PrepareHomeDirectory(config.LightningDir)
	if err != nil {
		panic(err)
	}
}

func OnRpcCommand(event *glightning.RpcCommandEvent) (*glightning.RpcCommandResponse, error) {

	metricsPlugin.HendlerRPCMessage(event)
	return event.Continue(), nil
}
