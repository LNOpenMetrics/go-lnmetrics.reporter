package plugin

import (
	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/vincenzopalazzo/glightning/jrpc2"
)

type PluginRpcMethod struct {
	metricsPlugin *MetricsPlugin `json:"-"`
}

// Would be cool if it is possible auto-generate a structure from a
// file, so this will be inside the binary and we can avoid the hard coded
// file.
type info struct {
	Name         string
	Version      string
	LangVersion  string
	Architecture string
	MaxProcs     int
	StoragePath  string
	Metrics      []string
	ProxyEnabled bool
}

func (instance PluginRpcMethod) Name() string {
	return "lnmetrics-info"
}

func NewPluginRpcMethod(pluginMetrics *MetricsPlugin) *PluginRpcMethod {
	return &PluginRpcMethod{
		metricsPlugin: pluginMetrics,
	}
}

func (instance *PluginRpcMethod) New() interface{} {
	return instance
}

func (instance *PluginRpcMethod) Call() (jrpc2.Result, error) {
	metricsSupp := make([]string, 0)
	for key := range instance.metricsPlugin.Metrics {
		metricsSupp = append(metricsSupp, MetricsSupported[key])
	}
	goInfo := sysinfo.Go()
	return info{
		Name:         "go-lnmetrics.reporter",
		Version:      "v0.0.4-rc4",
		LangVersion:  goInfo.Version,
		Architecture: goInfo.Arch,
		MaxProcs:     goInfo.MaxProcs,
		StoragePath:  instance.metricsPlugin.Storage.GetDBPath(),
		Metrics:      metricsSupp,
		ProxyEnabled: instance.metricsPlugin.WithProxy,
	}, nil
}
