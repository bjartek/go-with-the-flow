package gwtf

import (
	"fmt"
	"log"

	"github.com/enescakir/emoji"
	"github.com/spf13/afero"

	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/gateway"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
)

// GoWithTheFlow Entire configuration to work with Go With the Flow
type GoWithTheFlow struct {
	State                        *flowkit.State
	Services                     *services.Services
	Network                      string
	Logger                       output.Logger
	PrependNetworkToAccountNames bool
}

//NewGoWithTheFlowInMemoryEmulator this method is used to create an in memory emulator, deploy all contracts for the emulator and create all accounts
func NewGoWithTheFlowInMemoryEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "emulator", true).InitializeContracts().CreateAccounts("emulator-account")
}

func NewGoWithTheFlowForNetwork(network string) *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), network, false)

}

//NewGoWithTheFlowEmulator create a new client
func NewGoWithTheFlowEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "emulator", false)
}

func NewGoWithTheFlowDevNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "testnet", false)
}

func NewGoWithTheFlowMainNet() *GoWithTheFlow {
	return NewGoWithTheFlow(config.DefaultPaths(), "mainnet", false)
}

// NewGoWithTheFlow with custom file panic on error
func NewGoWithTheFlow(filenames []string, network string, inMemory bool) *GoWithTheFlow {
	gwtf, err := NewGoWithTheFlowError(filenames, network, inMemory)
	if err != nil {
		log.Fatalf("%v error %+v", emoji.PileOfPoo, err)
	}
	return gwtf
}

func (f *GoWithTheFlow) DoNotPrependNetworkToAccountNames() {
	f.PrependNetworkToAccountNames = false
}

func (f *GoWithTheFlow) Account(key string) *flowkit.Account {
	if f.PrependNetworkToAccountNames {
		key = fmt.Sprintf("%s-%s", f.Network, key)
	}

	account := f.State.Accounts().ByName(key)
	if account == nil {
		log.Fatalf("Could not find account with name %s", key)
	}

	return account

}

func (f *GoWithTheFlow) AccountName(key string) string {

	if f.PrependNetworkToAccountNames {
		return fmt.Sprintf("%s-%s", f.Network, key)
	}
	return key
}

// NewGoWithTheFlowError creates a new local go with the flow client
func NewGoWithTheFlowError(paths []string, network string, inMemory bool) (*GoWithTheFlow, error) {

	loader := &afero.Afero{Fs: afero.NewOsFs()}
	state, err := flowkit.Load(paths, loader)
	if err != nil {
		return nil, err
	}

	logger := output.NewStdoutLogger(output.InfoLog)

	var service *services.Services
	if inMemory {
		//YAY we can run it inline in memory!
		acc, _ := state.EmulatorServiceAccount()
		//TODO: How can i get the log output here? And enable verbose logging?
		gw := gateway.NewEmulatorGateway(acc)
		service = services.NewServices(gw, state, logger)
	} else {
		host := state.Networks().ByName(network).Host
		gw, err := gateway.NewGrpcGateway(host)
		if err != nil {
			log.Fatal(err)
		}
		service = services.NewServices(gw, state, logger)
	}
	return &GoWithTheFlow{
		State:                        state,
		Services:                     service,
		Network:                      network,
		Logger:                       logger,
		PrependNetworkToAccountNames: true,
	}, nil

}
