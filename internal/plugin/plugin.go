package plugin

import (
	"github.com/niftynei/glightning/glightning"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

type MetricsPlugin struct {
	Plugin *glightning.Plugin
}

func (plugin *MetricsPlugin) HendlerRPCMessage(event *glightning.RpcCommandEvent) error {
	command := event.Cmd
	method, err := command.Get()
	if err != nil {
		return err
	}
	switch method.(type) {
	case glightning.CloseRequest:
		log.GetInstance().Debug("Close command received")
	default:
		log.GetInstance().Debug("The node is up and runnning update the info")
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMethods() {
	rpcMethod := glightning.NewRpcMethod(&DiagnosticRpcMethod{}, "Show diagnostic node")
	rpcMethod.LongDesc = "Show the diagnostic data of the lightning network node"
	rpcMethod.Category = "metrics"
	plugin.Plugin.RegisterMethod(rpcMethod)
}
