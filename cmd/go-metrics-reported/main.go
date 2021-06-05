package main

import (
	"fmt"
	"os"

	metrics "github.com/OpenLNMetrics/go-metrics-reported/internal/plugin"

	"github.com/niftynei/glightning/glightning"
)

var metricsPlugin metrics.MetricsPlugin

func main() {
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
	//TODO init proxy and read the ln config
}

func OnRpcCommand(event *glightning.RpcCommandEvent) (*glightning.RpcCommandResponse, error) {

	metricsPlugin.HendlerRPCMessage(event)
	return event.Continue(), nil
}
