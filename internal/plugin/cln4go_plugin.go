package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/cache"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/metrics"
	jsonv2 "github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/json"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/trace"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	"github.com/elastic/go-sysinfo"
	cln4go "github.com/vincenzopalazzo/cln4go/plugin"
)

func ConfigureCLNPlugin[T MetricsPluginState](state T) (*cln4go.Plugin[T], error) {
	plugin := cln4go.New(state, false, OnInit[T])
	plugin.SetEncoder(&jsonv2.FastJSON{})
	plugin.SetTracer(&trace.Tracer{})

	plugin.RegisterOption("lnmetrics-urls", "string", "", "URLs of remote servers", false)
	plugin.RegisterOption("lnmetrics-noproxy", "bool", "false",
		"Disable the usage of proxy in case only for the go-lnmmetrics.reporter", false)
	plugin.RegisterNotification("shutdown", &OnShoutdown[T]{})

	plugin.RegisterRPCMethod("raw-local-score", "", "return the local reputation raw data collected by the plugin", NewRawLocalScoreRPC[T]())
	// FIXME: register the force rpc command only in developer mode
	plugin.RegisterRPCMethod("lnmetrics-force-update", "", "trigget the update to the server (caution)", &ForceUpdateRPC[T]{})
	plugin.RegisterRPCMethod("lnmetrics-info", "", "return the information regarding the lnmetrics plugin", &LNMetricsInfoRPC[T]{})
	plugin.RegisterRPCMethod("lnmetrics-clean", "", "clean the lnmetrics cache", &LNMetricsCleanCacheRPC[T]{})
	return plugin, nil
}

type OnShoutdown[T MetricsPluginState] struct{}

func (self *OnShoutdown[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) {
	// Share to all the metrics, so we need a global method that iterate over the metrics map
	params := make(map[string]any)
	params["timestamp"] = time.Now()
	msg := metrics.NewMsg("stop", params)
	for _, metric := range plugin.State.GetMetrics() {
		plugin.GetState().CallOnStopOnMetrics(metric, &msg)
	}
	plugin.GetState().GetCron().Stop()
	log.GetInstance().Info("Close command received")
}

type LNMetricsInfoRPC[T MetricsPluginState] struct{}

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

func (self *LNMetricsInfoRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	metricsSupp := make([]string, 0)
	for key := range plugin.GetState().GetMetrics() {
		metricsSupp = append(metricsSupp, metrics.MetricsSupported[key])
	}
	goInfo := sysinfo.Go()
	resp := info{
		Name:         "go-lnmetrics.reporter",
		Version:      "v0.0.5-rc2",
		LangVersion:  goInfo.Version,
		Architecture: goInfo.Arch,
		MaxProcs:     goInfo.MaxProcs,
		StoragePath:  plugin.State.GetStorage().GetDBPath(),
		Metrics:      metricsSupp,
		ProxyEnabled: plugin.GetState().IsWithProxy(),
	}
	// FIXME: support encode function inside the method
	var res map[string]any
	bytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type LNMetricsCleanCacheRPC[T MetricsPluginState] struct{}

func (self *LNMetricsCleanCacheRPC[T]) Call(plugin *cln4go.Plugin[T], payload map[string]any) (map[string]any, error) {
	cmdStr, found := payload["cmd"]
	if !found {
		return nil, fmt.Errorf("please specify the command that you want to call")
	}
	cmd := cmdStr.(string)
	if cmd != "clean" {
		return nil, fmt.Errorf("command %s node found", cmd)
	}

	if err := cache.GetInstance().CleanCache(); err != nil {
		log.GetInstance().Errorf("CleanCacheRPC rpc call return the following error during the %s cmd: %s", cmd, err)
		return nil, fmt.Errorf("clean cache operation return the following error: %s", err)
	}
	response := map[string]any{
		"result": "cleaning operation succeeded",
	}
	return response, nil
}
