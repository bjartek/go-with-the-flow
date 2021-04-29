package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/onflow/flow-cli/pkg/flowcli/config"
	"github.com/onflow/flow-cli/pkg/flowcli/gateway"
	"github.com/onflow/flow-cli/pkg/flowcli/output"
	"github.com/onflow/flow-cli/pkg/flowcli/project"
	"github.com/onflow/flow-cli/pkg/flowcli/services"
	"github.com/onflow/flow-go-sdk/crypto"
)

func start(flowConfigName string) {

	network := "emulator"
	//this is only for emulator right now
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

	log.Println("Deploying contracts")
	_, err = service.Project.Deploy(network, false)
	if err != nil {
		log.Fatal(err)
	}

	accounts := p.AccountNamesForNetwork(network)

	for _, accountName := range accounts {
		log.Println(fmt.Sprintf("Ensuring account with name '%s' is present", accountName))
		account := p.AccountByName(accountName)

		_, err := service.Accounts.Get(account.Address().String())
		if err == nil {
			logger.Debug("Account is present")
			continue
		}

		configuredAccount := p.AccountByName(accountName)
		key := configuredAccount.DefaultKey().ToConfig()
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
		if configuredAccount.Address() != a.Address {
			log.Fatal(errors.New(fmt.Sprintf("Configured account address != created address, %s != %s", configuredAccount.Address(), a.Address)))
		}
		log.Println("Account created " + a.Address.String())
	}

	tx, res, err := service.Transactions.Send(
		"transactions/create_nft_collection.cdc",
		"first",
		[]string{},
		"",
		"emulator")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(tx)
	spew.Dump(res)

	account := p.AccountByName("second").Address().Hex()
	value, err := service.Scripts.Execute(
		"scripts/test.cdc",
		[]string{fmt.Sprintf("Address:%s", account)},
		"",
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
