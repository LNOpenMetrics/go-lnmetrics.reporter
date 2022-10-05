package plugin

import (
	pluginDB "github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	cln4go "github.com/vincenzopalazzo/cln4go/client"
)

type MetricsPluginState interface {
	GetRpc() *cln4go.UnixRPC

	NewClient(path string) error

	SetStorage(storage pluginDB.PluginDatabase)

	GetStorage() pluginDB.PluginDatabase

	RegisterOneTimeEvt(after string)

	RegisterMetrics(id int, metrics Metric) error

	SetProxy(withProxy bool)

	IsWithProxy() bool

	SetServer(server *graphql.Client)

	GetServer() *graphql.Client
}
