package plugin

import (
	"time"

	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

func ConfigureCLNPlugin[T MetricsPluginState](state T) (*cln4go.Plugin[T], error) {
	plugin := cln4go.New(state, false, OnInit[T])

	plugin.RegisterOption("lnmetrics-urls", "string", "", "URLs of remote servers", false)
	plugin.RegisterOption("lnmetrics-noproxy", "flag", "false",
		"Disable the usage of proxy in case only for the go-lnmmetrics.reporter", false)
	plugin.RegisterNotification("shutdown", &OnShoutdown[T]{})
	return plugin, nil
}

type OnShoutdown[T MetricsPluginState] struct{}

func (self *OnShoutdown[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) {
	// Share to all the metrics, so we need a global method that iterate over the metrics map
	params := make(map[string]any)
	params["timestamp"] = time.Now()
	msg := Msg{"stop", params}
	for _, metric := range plugin.State.GetMetrics() {
		plugin.GetState().CallOnStopOnMetrics(metric, &msg)
	}
	plugin.GetState().GetCron().Stop()
	log.GetInstance().Info("Close command received")
}
