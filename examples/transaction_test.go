package main

import (
	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

/*
 Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestTransaction(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()
	t.Parallel()

	t.Run("Create NFT collection", func(t *testing.T) {
		g.TransactionFromFile("create_nft_collection").
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess().
			AssertNoEvents()
	})

	t.Run("Mint tokens assert events", func(t *testing.T) {
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3).
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensDeposited").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.MinterCreated").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", "A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", "A.0ae53cb6e3f42a79.FlowToken.MinterCreated").
			AssertEmitEvent(gwtf.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", map[string]interface{}{"amount": "100.00000000"}))

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
			AssertDebugLog("foobar")

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

	t.Run("Date string argument", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(timestamp:UFix64) {
  prepare(acct: AuthAccount) {
	Debug.log(timestamp.toString())
 }
}`).
			SignProposeAndPayAsService().
			DateStringAsUnixTimestamp("July 29, 2021 08:00:00 AM", "America/New_York").
			Test(t).
			AssertSuccess().
			AssertDebugLog("1627560000.00000000")
	})

	t.Run("Argument test", func(t *testing.T) {

		fix, _ := cadence.NewFix64("-1.0")
		ufix, _ := cadence.NewUFix64("1.0")

		builder :=g.Transaction("").BooleanArgument(true) .
			BytesArgument([]byte{ 1}).
			Fix64Argument("-1.0").
			UFix64Argument("1.0")
		assert.Contains(t, builder.Arguments, cadence.NewBool(true))
		assert.Contains(t, builder.Arguments, cadence.NewBytes([]byte{1}))
		assert.Contains(t, builder.Arguments, fix)
		assert.Contains(t, builder.Arguments, ufix)
	})

	t.Run("Word argument test", func(t *testing.T) {
		builder :=g.Transaction("").
			Word8Argument(8).
			Word16Argument(16).
			Word32Argument(32).
			Word64Argument(64)

		assert.Contains(t, builder.Arguments, cadence.NewWord8(8))
		assert.Contains(t, builder.Arguments, cadence.NewWord16(16))
		assert.Contains(t, builder.Arguments, cadence.NewWord32(32))
		assert.Contains(t, builder.Arguments, cadence.NewWord64(64))
	})

	t.Run("UInt argument test", func(t *testing.T) {
		builder :=g.Transaction("").
			UIntArgument(1).
			UInt8Argument(8).
			UInt16Argument(16).
			UInt32Argument(32).
			UInt64Argument(64).
			UInt128Argument(128).
			UInt256Argument(256)

		assert.Contains(t, builder.Arguments, cadence.NewUInt(1))
		assert.Contains(t, builder.Arguments, cadence.NewUInt8(8))
		assert.Contains(t, builder.Arguments, cadence.NewUInt16(16))
		assert.Contains(t, builder.Arguments, cadence.NewUInt32(32))
		assert.Contains(t, builder.Arguments, cadence.NewUInt64(64))
		assert.Contains(t, builder.Arguments, cadence.NewUInt128(128))
		assert.Contains(t, builder.Arguments, cadence.NewUInt256(256))
	})

	t.Run("Int argument test", func(t *testing.T) {
		builder :=g.Transaction("").
			IntArgument(1).
			Int8Argument(-8).
			Int16Argument(16).
			Int32Argument(32).
			Int64Argument(64).
			Int128Argument(128).
			Int256Argument(256)

		assert.Contains(t, builder.Arguments, cadence.NewInt(1))
		assert.Contains(t, builder.Arguments, cadence.NewInt8(-8))
		assert.Contains(t, builder.Arguments, cadence.NewInt16(16))
		assert.Contains(t, builder.Arguments, cadence.NewInt32(32))
		assert.Contains(t, builder.Arguments, cadence.NewInt64(64))
		assert.Contains(t, builder.Arguments, cadence.NewInt128(128))
		assert.Contains(t, builder.Arguments, cadence.NewInt256(256))
	})
}

