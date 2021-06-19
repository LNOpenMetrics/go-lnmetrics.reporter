package plugin

import (
	"fmt"
	"github.com/niftynei/glightning/jrpc2"
)

type DiagnosticRpcMethod struct {
	// all the error that will happen durng the metrics
	errors   map[int][]string
	memoryDB *MemoryDB
}

func (rpc *DiagnosticRpcMethod) Name() string {
	return "diagnostic"
}

func NewMetricPlugin() *DiagnosticRpcMethod {
	return &DiagnosticRpcMethod{errors: make(map[int][]string),
		memoryDB: &MemoryDB{nodeId: "", metrics: make(map[int]*Metric)}}
}

func (rpc *DiagnosticRpcMethod) New() interface{} {
	return NewMetricPlugin()
}

func (rpc *DiagnosticRpcMethod) Call() (jrpc2.Result, error) {
	return fmt.Sprintf("Here will be the diagnostic"), nil
}

func (rpc *DiagnosticRpcMethod) UpdateWithMessage(msg *Msg) (jrpc2.Result, error) {
	return fmt.Sprintf("Here will be the diagnostic"), nil
}
