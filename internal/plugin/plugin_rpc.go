package plugin

import (
	"github.com/vincenzopalazzo/glightning/jrpc2"
)

type PluginRpcMethod struct{}

// Would be cool if it is possible auto-generate a structure from a
// file, so this will be inside the binary and we can avoid the hard coded
// file.
type info struct {
	Name         string
	Version      string
	LangVersion  string
	Architecture string
}

func (instance PluginRpcMethod) Name() string {
	return "lnmetrics-reporter"
}

func NewPluginRpcMethod() *PluginRpcMethod {
	return &PluginRpcMethod{}
}

func (instance PluginRpcMethod) New() interface{} {
	return NewPluginRpcMethod()
}

func (instance *PluginRpcMethod) Call() (jrpc2.Result, error) {
	return info{Name: "go-lnmetrics-reporter", Version: "0.1", LangVersion: "Go lang 1.15.8", Architecture: "amd64"}, nil
}
