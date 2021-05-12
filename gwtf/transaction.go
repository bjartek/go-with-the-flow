package gwtf

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

//TransactionFromFile will start a flow transaction builder
func (f *GoWithTheFlow) TransactionFromFile(filename string) FlowTransactionBuilder {
	return FlowTransactionBuilder{
		GoWithTheFlow:  f,
		FileName:       filename,
		MainSigner:     nil,
		Arguments:      []cadence.Value{},
		PayloadSigners: []*GoWithTheFlowAccount{},
	}
}

//SignProposeAndPayAsService set the payer, proposer and envelope signer
func (t FlowTransactionBuilder) SignProposeAndPayAsService() FlowTransactionBuilder {
	mainSigner := t.GoWithTheFlow.Service
	t.MainSigner = &mainSigner
	return t
}

//SignProposeAndPayAs set the payer, proposer and envelope signer
func (t FlowTransactionBuilder) SignProposeAndPayAs(signer string) FlowTransactionBuilder {
	mainSigner := t.GoWithTheFlow.Accounts[signer]
	t.MainSigner = &mainSigner
	return t
}

//RawAccountArgument add an account from a string as an argument
func (t FlowTransactionBuilder) RawAccountArgument(key string) FlowTransactionBuilder {

	account := flow.HexToAddress(key)
	accountArg := cadence.BytesToAddress(account.Bytes())
	return t.Argument(accountArg)
}

//AccountArgument add an account as an argument
func (t FlowTransactionBuilder) AccountArgument(key string) FlowTransactionBuilder {
	f := t.GoWithTheFlow
	address := cadence.BytesToAddress(f.Accounts[key].Address.Bytes())
	return t.Argument(address)
}

//StringArgument add a String Argument to the transaction
func (t FlowTransactionBuilder) StringArgument(value string) FlowTransactionBuilder {
	return t.Argument(cadence.String(value))
}

//BooleanArgument add a Boolean Argument to the transaction
func (t FlowTransactionBuilder) BooleanArgument(value bool) FlowTransactionBuilder {
	return t.Argument(cadence.NewBool(value))
}

//BytesArgument add a Bytes Argument to the transaction
func (t FlowTransactionBuilder) BytesArgument(value []byte) FlowTransactionBuilder {
	return t.Argument(cadence.NewBytes(value))
}

//IntArgument add an Int Argument to the transaction
func (t FlowTransactionBuilder) IntArgument(value int) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt(value))
}

//Int8Argument add an Int8 Argument to the transaction
func (t FlowTransactionBuilder) Int8Argument(value int8) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt8(value))
}

//Int16Argument add an Int16 Argument to the transaction
func (t FlowTransactionBuilder) Int16Argument(value int16) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt16(value))
}

//Int32Argument add an Int32 Argument to the transaction
func (t FlowTransactionBuilder) Int32Argument(value int32) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt32(value))
}

//Int64Argument add an Int64 Argument to the transaction
func (t FlowTransactionBuilder) Int64Argument(value int64) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt64(value))
}

//Int128Argument add an Int128 Argument to the transaction
func (t FlowTransactionBuilder) Int128Argument(value int) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt128(value))
}

//Int256Argument add an Int256 Argument to the transaction
func (t FlowTransactionBuilder) Int256Argument(value int) FlowTransactionBuilder {
	return t.Argument(cadence.NewInt256(value))
}

//UIntArgument add an UInt Argument to the transaction
func (t FlowTransactionBuilder) UIntArgument(value uint) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt(value))
}

//UInt8Argument add an UInt8 Argument to the transaction
func (t FlowTransactionBuilder) UInt8Argument(value uint8) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt8(value))
}

//UInt16Argument add an UInt16 Argument to the transaction
func (t FlowTransactionBuilder) UInt16Argument(value uint16) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt16(value))
}

//UInt32Argument add an UInt32 Argument to the transaction
func (t FlowTransactionBuilder) UInt32Argument(value uint32) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt32(value))
}

//UInt64Argument add an UInt64 Argument to the transaction
func (t FlowTransactionBuilder) UInt64Argument(value uint64) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt64(value))
}

//UInt128Argument add an UInt128 Argument to the transaction
func (t FlowTransactionBuilder) UInt128Argument(value uint) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt128(value))
}

//UInt256Argument add an UInt256 Argument to the transaction
func (t FlowTransactionBuilder) UInt256Argument(value uint) FlowTransactionBuilder {
	return t.Argument(cadence.NewUInt256(value))
}

//Word8Argument add a Word8 Argument to the transaction
func (t FlowTransactionBuilder) Word8Argument(value uint8) FlowTransactionBuilder {
	return t.Argument(cadence.NewWord8(value))
}

//Word16Argument add a Word16 Argument to the transaction
func (t FlowTransactionBuilder) Word16Argument(value uint16) FlowTransactionBuilder {
	return t.Argument(cadence.NewWord16(value))
}

//Word32Argument add a Word32 Argument to the transaction
func (t FlowTransactionBuilder) Word32Argument(value uint32) FlowTransactionBuilder {
	return t.Argument(cadence.NewWord32(value))
}

//Word64Argument add a Word64 Argument to the transaction
func (t FlowTransactionBuilder) Word64Argument(value uint64) FlowTransactionBuilder {
	return t.Argument(cadence.NewWord64(value))
}

//Fix64Argument add a Fix64 Argument to the transaction
func (t FlowTransactionBuilder) Fix64Argument(value string) FlowTransactionBuilder {
	amount, err := cadence.NewFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

//UFix64Argument add a UFix64 Argument to the transaction
func (t FlowTransactionBuilder) UFix64Argument(value string) FlowTransactionBuilder {
	amount, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

//Argument add an argument to the transaction
func (t FlowTransactionBuilder) Argument(value cadence.Value) FlowTransactionBuilder {
	t.Arguments = append(t.Arguments, value)
	return t
}

//PayloadSigner set a signer for the payload
func (t FlowTransactionBuilder) PayloadSigner(value string) FlowTransactionBuilder {
	signer := t.GoWithTheFlow.Accounts[value]
	t.PayloadSigners = append(t.PayloadSigners, &signer)
	return t
}

//RunPrintEventsFull will run a transaction and print all events
func (t FlowTransactionBuilder) RunPrintEventsFull() {
	PrintEvents(t.Run(), map[string][]string{})
}

//RunPrintEvents will run a transaction and print all events
func (t FlowTransactionBuilder) RunPrintEvents(ignoreFields map[string][]string) {
	PrintEvents(t.Run(), ignoreFields)
}

//Run run the transaction
func (t FlowTransactionBuilder) Run() []flow.Event {
	if t.MainSigner == nil {
		log.Fatalf("%v You need to set the main signer", emoji.PileOfPoo)
	}
	txFilePath := fmt.Sprintf("./transactions/%s.cdc", t.FileName)
	code, err := ioutil.ReadFile(txFilePath)
	if err != nil {
		log.Fatalf("%v Could not read transaction file from path=%s", emoji.PileOfPoo, txFilePath)
	}
	tx := flow.NewTransaction().SetScript(code)

	events, err := t.GoWithTheFlow.performTransaction(tx, t.MainSigner, t.PayloadSigners, t.Arguments)
	if err != nil {
		log.Fatalf("%v error sending transaction %s %+v", emoji.PileOfPoo, t.FileName, err)
	}
	log.Printf("%v Transaction %s successfully applied\n", emoji.OkHand, t.FileName)
	return events
}

//FlowTransactionBuilder used to create a builder pattern for a transaction
type FlowTransactionBuilder struct {
	GoWithTheFlow  *GoWithTheFlow
	FileName       string
	Arguments      []cadence.Value
	MainSigner     *GoWithTheFlowAccount
	PayloadSigners []*GoWithTheFlowAccount
}
