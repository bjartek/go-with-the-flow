package main

import "github.com/bjartek/go-with-the-flow/gwtf"

func main() {

	gwtf := gwtf.NewGoWithTheFlowEmulator()

	gwtf.DeployContract("nft")
	gwtf.DeployContract("ft")
	gwtf.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("ft").Run()

	gwtf.TransactionFromFile("arguments").SignProposeAndPayAs("ft").StringArgument("argument1").Run()

	gwtf.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("ft").AccountArgument("nft").Run()
	gwtf.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("ft").PayloadSigner("nft").Run()
	/*
		// Run Script
		gwtf.RunScriptReturns("test", gwtf.FindAddress("nft"))

		//Run script that returns
		result := gwtf.RunScriptReturns("test", gwtf.FindAddress("nft"))
		log.Printf("Script returned %s", result)
	*/

}
