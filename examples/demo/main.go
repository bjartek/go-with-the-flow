package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/gateway"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"
)

func start(flowConfigName string) {

	network := "emulator"
	//this is only for emulator right now
	saAccountName := "emulator-account"

	loader := &afero.Afero{Fs: afero.NewOsFs()}
	p, err := flowkit.Load([]string{flowConfigName}, loader)
	if err != nil {
		log.Fatal(err)
	}

	logger := output.NewStdoutLogger(output.DebugLog)
	gateway, err := gateway.NewGrpcGateway(config.DefaultEmulatorNetwork().Host)
	if err != nil {
		log.Fatal(err)
	}

	service := services.NewServices(gateway, p, logger)

	log.Println("Deploying contracts")
	_, err = service.Project.Deploy(network, false)
	if err != nil {
		log.Fatal(err)
	}

	signerAccount := p.Accounts().ByName(saAccountName)

	accounts := p.AccountNamesForNetwork(network)

	for _, accountName := range accounts {
		log.Println(fmt.Sprintf("Ensuring account with name '%s' is present", accountName))
		account := p.Accounts().ByName(accountName)

		_, err := service.Accounts.Get(account.Address())
		if err == nil {
			logger.Debug("Account is present")
			continue
		}

		a, err := service.Accounts.Create(
			signerAccount,
			[]crypto.PublicKey{account.Key().ToConfig().PrivateKey.PublicKey()},
			[]int{1000},
			account.Key().SigAlgo(),
			account.Key().HashAlgo(),
			[]string{})
		if err != nil {
			log.Fatal(err)
		}
		if account.Address() != a.Address {
			log.Fatal(errors.New(fmt.Sprintf("Configured account address != created address, %s != %s", account.Address(), a.Address)))
		}
		log.Println("Account created " + a.Address.String())
	}

	fileName := "transactions/create_nft_collection.cdc"
	code, err := p.ReaderWriter().ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	tx, res, err := service.Transactions.Send(
		p.Accounts().ByName("first"),
		code,
		fileName,
		9999,
		[]cadence.Value{},
		"emulator")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(tx)
	spew.Dump(res)

	scriptName := "scripts/test.cdc"

	script, err := p.ReaderWriter().ReadFile(scriptName)
	if err != nil {
		log.Fatal(err)
	}
	value, err := service.Scripts.Execute(
		script,
		[]cadence.Value{cadence.Address(signerAccount.Address())},
		scriptName,
		"emulator")

	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(value)

}

func main() {

	start("flow.json")
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
