package plugin

import (
	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

func ConfigureCLNPlugin[T MetricsPluginState](state T) (*cln4go.Plugin[T], error) {
	plugin := cln4go.New(state, false, OnInit[T])

	plugin.RegisterOption("lnmetrics-urls", "string", "", "URLs of remote servers", false)
	plugin.RegisterOption("lnmetrics-noproxy", "flag", "false",
		"Disable the usage of proxy in case only for the go-lnmmetrics.reporter", false)

	return plugin, nil
}
