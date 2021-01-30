package main

import "github.com/bjartek/go-with-the-flow/gwtf"

func main() {

	gwtf := gwtf.NewGoWithTheFlowEmulator()

	gwtf.CreateAccount("accounts", "first")
	gwtf.TransactionFromFile("mint_tokens").SignProposeAndPayAsService().AccountArgument("first").UFix64Argument("10.0").RunPrintEventsFull()

}
