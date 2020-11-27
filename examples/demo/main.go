package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowEmulator()

	var ignoreFields = map[string][]string{
		"flow.AccountCodeUpdated": []string{"codeHash"},
		"flow.AccountKeyAdded":    []string{"publicKey"},
	}
	g.InitializeContracts()
	g.CreateAccount("first")
	g.CreateAccount("second")

	g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunPrintEvents(ignoreFields)

	g.TransactionFromFile("arguments").SignProposeAndPayAs("first").StringArgument("argument1").RunPrintEvents(ignoreFields)

	g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("first").AccountArgument("second").RunPrintEvents(ignoreFields)
	g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("first").PayloadSigner("second").RunPrintEvents(ignoreFields)

	g.ScriptFromFile("test").AccountArgument("second").Run()

	//Run script that returns
	result := g.ScriptFromFile("test").AccountArgument("second").RunReturns()
	log.Printf("Script returned %s", result)

}
