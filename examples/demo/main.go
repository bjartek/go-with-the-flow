package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowEmulator()

	gwtf.PrintEvents(g.DeployContract("nft"))

	g.DeployContract("ft")
	g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("ft").Run()

	g.TransactionFromFile("arguments").SignProposeAndPayAs("ft").StringArgument("argument1").Run()

	g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("ft").AccountArgument("nft").Run()
	g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("ft").PayloadSigner("nft").Run()

	g.ScriptFromFile("test").AccountArgument("nft").Run()

	//Run script that returns
	result := g.ScriptFromFile("test").AccountArgument("nft").RunReturns()
	log.Printf("Script returned %s", result)

}
