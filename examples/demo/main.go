package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()
	g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunPrintEventsFull()
	g.TransactionFromFile("arguments").SignProposeAndPayAs("first").StringArgument("argument1").RunPrintEventsFull()
	g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("first").AccountArgument("second").RunPrintEventsFull()

	//multiple signers is not something I have handled yet.
	//g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("first").PayloadSigner("second").RunPrintEventsFull()

	g.ScriptFromFile("test").AccountArgument("second").Run()
	g.TransactionFromFile("mint_tokens").SignProposeAndPayAs("emulator-account").AccountArgument("first").UFix64Argument("10.0").RunPrintEventsFull()


	g.Script(`
import NonFungibleToken from "./contracts/NonFungibleToken.cdc"

pub fun main(account: Address): String {
    return getAccount(account).address.toString()
}`).AccountArgument("second").Run()

	g.Transaction(`
transaction(test:String) {
  prepare(acct: AuthAccount) {
    log(acct)
    log(test)
 }
}`).SignProposeAndPayAs("first").StringArgument("foobar").Run()

	//Run script that returns
	result := g.ScriptFromFile("test").AccountArgument("second").RunReturns()
	log.Printf("Script returned %s", result)

}
