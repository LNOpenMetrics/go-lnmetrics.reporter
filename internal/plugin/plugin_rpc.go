package plugin

import (
	"fmt"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/cache"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/vincenzopalazzo/glightning/jrpc2"
)

type PluginRpcMethod struct {
	metricsPlugin *MetricsPlugin `json:"-"`
}

// Would be cool if it is possible auto-generate a structure from a
// file, so this will be inside the binary, and we can avoid the hard coded
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

func NewPluginRpcMethod(pluginMetrics *MetricsPlugin) *PluginRpcMethod {
	return &PluginRpcMethod{
		metricsPlugin: pluginMetrics,
	}
}

func (instance PluginRpcMethod) Name() string {
	return "lnmetrics-info"
}

func (instance *PluginRpcMethod) New() any {
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
		Version:      "v0.0.5-rc1",
		LangVersion:  goInfo.Version,
		Architecture: goInfo.Arch,
		MaxProcs:     goInfo.MaxProcs,
		StoragePath:  instance.metricsPlugin.Storage.GetDBPath(),
		Metrics:      metricsSupp,
		ProxyEnabled: instance.metricsPlugin.WithProxy,
	}, nil
}

// CleanCacheRPC RPC call to clean up the plugin cache.
type CleanCacheRPC struct {
	Cmd string
}

func NewCleanCacheRPC(plugin *MetricsPlugin) *CleanCacheRPC {
	return &CleanCacheRPC{}
}

func (instance *CleanCacheRPC) New() any {
	return instance
}

func (instance CleanCacheRPC) Name() string {
	return "lnmetrics-cache"
}

func (instance *CleanCacheRPC) Call() (jrpc2.Result, error) {

	if instance.Cmd != "clean" {
		return nil, fmt.Errorf("command %s node found", instance.Cmd)
	}

	if err := cache.GetInstance().CleanCache(); err != nil {
		log.GetInstance().Errorf("CleanCacheRPC rpc call return the following error during the %s cmd: %s", instance.Cmd, err)
		return nil, fmt.Errorf("clean cache operation return the following error: %s", err)
	}
	response := struct {
		Result string `json:"result"`
	}{Result: "cleaning operation succeeded"}
	return response, nil
}
