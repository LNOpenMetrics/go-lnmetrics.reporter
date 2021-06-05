package plugin

import (
	"fmt"
	"github.com/niftynei/glightning/jrpc2"
)

type DiagnosticRpcMethod struct{}

func (rpc *DiagnosticRpcMethod) Name() string {
	return "diagnostic"
}

func (rpc *DiagnosticRpcMethod) New() interface{} {
	return &DiagnosticRpcMethod{}
}

func (rpc *DiagnosticRpcMethod) Call() (jrpc2.Result, error) {
	return fmt.Sprintf("Here will be the diagnostic"), nil
}
