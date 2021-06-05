package plugin

import (
	"fmt"
	"github.com/niftynei/glightning/glightning"
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
		fmt.Println("Close command received")
	default:
		fmt.Println("The node is up and runnning update the info")
	}
	return nil
}

func (plugin *MetricsPlugin) RegisterMethods() {
	rpcMethod := glightning.NewRpcMethod(&DiagnosticRpcMethod{}, "Example rpc method")
	rpcMethod.LongDesc = "Show the diagnostic data of the lightning network node"
	rpcMethod.Category = "metrics"
	plugin.Plugin.RegisterMethod(rpcMethod)
}
