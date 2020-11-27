package gwtf

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
)

//CreateAccount with no contracts
func (f *GoWithTheFlow) CreateAccount(accountName ...string) *GoWithTheFlow {
	for _, account := range accountName {
		f.CreateAccountWithContracts(account)
	}
	return f
}

//InitializeContracts installs all contracts to the account with name accounts
func (f *GoWithTheFlow) InitializeContracts() *GoWithTheFlow {
	f.CreateAccountWithAllContracts("accounts")
	return f
}

//CreateAccountWithAllContracts with all contracts in folder
func (f *GoWithTheFlow) CreateAccountWithAllContracts(accountName string) []flow.Event {
	files, err := ioutil.ReadDir("./contracts/")
	if err != nil {
		log.Fatal(err)
	}
	var contractNames []string
	for _, f := range files {
		fileName := f.Name()
		name := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		contractNames = append(contractNames, name)
	}
	return f.CreateAccountWithContracts(accountName, contractNames...)
}

// CreateAccountWithContracts will create an account for running transactions without a contract
func (f *GoWithTheFlow) CreateAccountWithContracts(accountName string, contractsStrings ...string) []flow.Event {

	var contracts []templates.Contract
	for _, contractName := range contractsStrings {
		contractPath := fmt.Sprintf("./contracts/%s.cdc", contractName)
		code, err := ioutil.ReadFile(contractPath)
		if err != nil {
			log.Fatalf("%v Could not read contract file from path=%s", emoji.PileOfPoo, contractPath)
		}
		log.Printf("%v create contract %s in %s\n", emoji.Scroll, contractName, accountName)
		contracts = append(contracts, templates.Contract{
			Name:   contractName,
			Source: string(code),
		})
	}
	user := f.Accounts[accountName]
	tx := templates.CreateAccount([]*flow.AccountKey{user.NewAccountKey()}, contracts, f.Service.Address)
	events, err := f.performTransaction(tx, &f.Service, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v error creating account %s %+v", emoji.PileOfPoo, accountName, err)
	}
	log.Printf("%v Account created: %s \n", emoji.Person, accountName)
	return events
}

// AddContract will deploy a contract with the given name to an account with the same name from wallet.json
func (f *GoWithTheFlow) AddContract(accountName string, contractName string) []flow.Event {
	contractPath := fmt.Sprintf("./contracts/%s.cdc", contractName)
	log.Printf("Deploying contract: %s at %s", contractName, contractPath)
	code, err := ioutil.ReadFile(contractPath)
	if err != nil {
		log.Fatalf("%v Could not read contract file from path=%s", emoji.PileOfPoo, contractPath)
	}
	user := f.Accounts[accountName]
	contract := templates.Contract{
		Name:   contractName,
		Source: string(code),
	}
	tx := templates.AddAccountContract(user.Address, contract)
	events, err := f.performTransaction(tx, &user, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v adding contract %s %+v", emoji.PileOfPoo, contractName, err)
	}
	log.Printf("%v Contract: %s successfully deployed\n", emoji.Scroll, contractName)
	return events
}
