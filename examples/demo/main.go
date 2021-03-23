package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

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

}
