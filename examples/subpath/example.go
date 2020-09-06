package main

import (
	"github.com/versus-flow/go-flow-tooling/tooling"
)

func main() {

	flow := tooling.NewFlowConfigLocalhostWithParentPath("..")
	flow.RunScriptReturns("block")

}
