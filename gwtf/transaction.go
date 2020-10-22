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

//SignProposeAndPayAs set the payer, proposer and envelope signer
func (t FlowTransactionBuilder) SignProposeAndPayAs(signer string) FlowTransactionBuilder {
	mainSigner := t.GoWithTheFlow.Accounts[signer]
	t.MainSigner = &mainSigner
	return t
}

//AccountArgument add an account as an argument
func (t FlowTransactionBuilder) AccountArgument(key string) FlowTransactionBuilder {
	f := t.GoWithTheFlow
	address := cadence.BytesToAddress(f.Accounts[key].Address.Bytes())
	return t.Argument(address)
}

//StringArgument add an argument to the transaction
func (t FlowTransactionBuilder) StringArgument(value string) FlowTransactionBuilder {
	t.Arguments = append(t.Arguments, cadence.String(value))
	return t
}

//UFix64Argument add an account as an argument
func (t FlowTransactionBuilder) UFix64Argument(key string) FlowTransactionBuilder {
	amount, err := cadence.NewUFix64(key)
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
