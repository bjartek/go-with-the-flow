package main

import (
	"log"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	//This special New method will setup and in memory emulator, deploy all contracts, create all acconts that does not have contracts in deploy block and prepare for unit testing or like this an demo script
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()
	g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunPrintEventsFull()
	g.TransactionFromFile("arguments").SignProposeAndPayAs("first").StringArgument("argument1").RunPrintEventsFull()
	g.TransactionFromFile("argumentsWithAccount").SignProposeAndPayAs("first").AccountArgument("second").RunPrintEventsFull()
	g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("first").PayloadSigner("second").StringArgument("asserts.go").RunPrintEventsFull()

	g.ScriptFromFile("asserts.go").AccountArgument("second").Run()
	g.TransactionFromFile("mint_tokens").SignProposeAndPayAsService().AccountArgument("first").UFix64Argument("10.0").RunPrintEventsFull()

	g.Script(`
pub fun main(account: Address): String {
    return getAccount(account).address.toString()
}`).AccountArgument("second").Run()

	g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(value:String) {
  prepare(acct: AuthAccount) {
	Debug.log(value)
 }
}`).SignProposeAndPayAs("first").StringArgument("foobar").RunPrintEventsFull()

	//Run script that returns
	result := g.ScriptFromFile("test").AccountArgument("second").RunReturns()
	log.Printf("Script returned %s", result)

}
