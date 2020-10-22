package gwtf

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
)

// CreateAccount will create an account for running transactions without a contract
func (f *GoWithTheFlow) CreateAccount(accountName string) []flow.Event {
	user := f.Accounts[accountName]
	tx := templates.CreateAccount([]*flow.AccountKey{user.NewAccountKey()}, nil, f.Service.Address)
	events, err := f.performTransaction(tx, &f.Service, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v error creating account %s %+v", emoji.PileOfPoo, accountName, err)
	}
	log.Printf("%v Account created: %s \n", emoji.Scroll, accountName)
	return events
}

// DeployContract will deploy a contract with the given name to an account with the same name from wallet.json
func (f *GoWithTheFlow) DeployContract(contractName string) []flow.Event {
	contractPath := fmt.Sprintf("./contracts/%s.cdc", contractName)
	//log.Printf("Deploying contract: %s at %s", contractName, contractPath)
	code, err := ioutil.ReadFile(contractPath)
	if err != nil {
		log.Fatalf("%v Could not read contract file from path=%s", emoji.PileOfPoo, contractPath)
	}
	user := f.Accounts[contractName]
	tx := templates.CreateAccount([]*flow.AccountKey{user.NewAccountKey()}, code, f.Service.Address)
	events, err := f.performTransaction(tx, &f.Service, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v error creating account %s %+v", emoji.PileOfPoo, contractName, err)
	}
	log.Printf("%v Contract: %s successfully deployed\n", emoji.Scroll, contractName)
	return events
}
