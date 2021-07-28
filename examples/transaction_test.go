package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

/*
 Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestTransaction(t *testing.T) {
	g := gwtf.NewTestingEmulator()
	t.Parallel()

	t.Run("Create NFT collection", func(t *testing.T) {
		g.TransactionFromFile("create_nft_collection").
			SignProposeAndPayAs("first").
			Test(t). //This method will return a TransactionResult that we can assert upon
			AssertSuccess().  //Assert that there are no errors and that the transactions succeeds
			AssertNoEvents()  //Assert that we did not emit any events.
	})

	t.Run("Mint tokens assert events", func(t *testing.T) {
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3). //assert the number of events returned
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted"). //assert the name of a single event
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", "A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", "A.0ae53cb6e3f42a79.FlowToken.MinterCreated"). //or assert more then one eventname in a go
			AssertEmitEvent(gwtf.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", map[string]interface{}{"amount": "100.00000000"})) //assert a given event, can also take multiple events if you like

	})

	t.Run("Inline transaction with debug log", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(message:String) {
  prepare(acct: AuthAccount, account2: AuthAccount) {
	Debug.log(message)
 }
}`).
			SignProposeAndPayAs("first").
			PayloadSigner("second").
			StringArgument("foobar").
			Test(t).
			AssertSuccess().
			AssertDebugLog("foobar") //assert that we have debug logged something. The assertion is contains so you do not need to write the entire debug log output if you do not like

	})

	t.Run("Raw account argument", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(user:Address) {
  prepare(acct: AuthAccount) {
	Debug.log(user.toString())
 }
}`).
			SignProposeAndPayAsService().
			RawAccountArgument("0x1cf0e2f2f715450").
			Test(t).
			AssertSuccess().
			AssertDebugLog("0x1cf0e2f2f715450")
	})

	t.Run("transaction that should fail", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(user:Address) {
  prepare(acct: AuthAccount) {
	Debug.log(user.toStrig())
 }
}`).
			SignProposeAndPayAsService().
			RawAccountArgument("0x1cf0e2f2f715450").
			Test(t).
			AssertFailure("has no member `toStrig`") //assert failure with an error message. uses contains so you do not need to write entire message
	})

}
