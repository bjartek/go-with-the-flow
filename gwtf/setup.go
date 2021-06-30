package gwtf

import (
	"log"

	"github.com/enescakir/emoji"
	"github.com/spf13/afero"

	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/gateway"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
)

// DiscordWebhook stores information about a webhook
type DiscordWebhook struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Wait  bool   `json:"wait"`
}

// GoWithTheFlow Entire configuration to work with Go With the Flow
type GoWithTheFlow struct {
	State    *flowkit.State
	Services *services.Services
	Network  string
}

//NewGoWithTheFlowEmulator create a new client
func NewGoWithTheFlowEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "emulator")
}

// NewGoWithTheFlowDevNet setup dev like in https://www.notion.so/Accessing-Flow-Devnet-ad35623797de48c08d8b88102ea38131
func NewGoWithTheFlowDevNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "devnet")
}

// NewGoWithTheFlowDevNet setup dev like in https://www.notion.so/Accessing-Flow-Devnet-ad35623797de48c08d8b88102ea38131
func NewGoWithTheFlowMainNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "mainnet")
}

// NewGoWithTheFlow with custom file panic on error
func NewGoWithTheFlow(filenames []string, network string) *GoWithTheFlow {
	gwtf, err := NewGoWithTheFlowError(filenames, network)
	if err != nil {
		log.Fatalf("%v error %+v", emoji.PileOfPoo, err)
	}
	return gwtf
}

// NewGoWithTheFlowError creates a new local go with the flow client
func NewGoWithTheFlowError(paths []string, network string) (*GoWithTheFlow, error) {

	loader := &afero.Afero{Fs: afero.NewOsFs()}
	p, err := flowkit.Load(paths, loader)
	if err != nil {
		return nil, err
	}

	logger := output.NewStdoutLogger(output.DebugLog)
	gateway, err := gateway.NewGrpcGateway(config.DefaultEmulatorNetwork().Host)
	if err != nil {
		log.Fatal(err)
	}

	service := services.NewServices(gateway, p, logger)

	return &GoWithTheFlow{
		State:    p,
		Services: service,
		Network:  network,
	}, nil

}
