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
	gwtf.PrintEvents(g.DeployContract("nft"), ignoreFields)

	g.DeployContract("ft")
	g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("ft").RunPrintEvents(ignoreFields)

	g.TransactionFromFile("arguments").SignProposeAndPayAs("ft").StringArgument("argument1").RunPrintEvents(ignoreFields)

	g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("ft").AccountArgument("nft").RunPrintEvents(ignoreFields)
	g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("ft").PayloadSigner("nft").RunPrintEvents(ignoreFields)

	g.ScriptFromFile("test").AccountArgument("nft").Run()

	//Run script that returns
	result := g.ScriptFromFile("test").AccountArgument("nft").RunReturns()
	log.Printf("Script returned %s", result)

}
