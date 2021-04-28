package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/onflow/flow-cli/pkg/flowcli/config"
	"github.com/onflow/flow-cli/pkg/flowcli/gateway"
	"github.com/onflow/flow-cli/pkg/flowcli/output"
	"github.com/onflow/flow-cli/pkg/flowcli/project"
	"github.com/onflow/flow-cli/pkg/flowcli/services"
	"github.com/onflow/flow-go-sdk/crypto"
)

func ca(flowConfigName string, accountName string) {

	saAccountName := "emulator-account"
	p, err := project.Load([]string{flowConfigName})
	if err != nil {
		log.Fatal(err)
	}

	logger := output.NewStdoutLogger(output.DebugLog)
	gateway, err := gateway.NewGrpcGateway(config.DefaultEmulatorNetwork().Host)
	if err != nil {
		log.Fatal(err)
	}

	service := services.NewServices(gateway, p, logger)

	//This feels realy cumbersome just to get a Public key from a private key in config
	key := p.AccountByName(accountName).DefaultKey().ToConfig()
	pk := key.Context[config.PrivateKeyField]
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, pk)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.PublicKey().String()

	a, err := service.Accounts.Create(
		saAccountName,
		[]string{publicKey},
		[]int{1000},
		key.SigAlgo.String(),
		key.HashAlgo.String(),
		[]string{})
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(a)

}
func main() {

	ca("flow.json", "first")
	/*
		g := gwtf.
			NewGoWithTheFlowEmulator().
			CreateAccountPrintEvents("first", "second")

		g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunPrintEventsFull()
		g.TransactionFromFile("arguments").SignProposeAndPayAs("first").StringArgument("argument1").RunPrintEventsFull()
		g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("first").AccountArgument("second").RunPrintEventsFull()
		g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("first").PayloadSigner("second").RunPrintEventsFull()

		g.ScriptFromFile("test").AccountArgument("second").Run()

		//Run script that returns
		result := g.ScriptFromFile("test").AccountArgument("second").RunReturns()
		log.Printf("Script returned %s", result)
	*/

}
