package gwtf

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
)

//CreateAccountPrintEvents create accounts and print events
func (f *GoWithTheFlow) CreateAccountPrintEvents(accountName ...string) *GoWithTheFlow {

	var ignoreFields = map[string][]string{
		"flow.AccountCodeUpdated": {"codeHash"},
		"flow.AccountKeyAdded":    {"publicKey"},
	}

	for _, account := range accountName {
		PrintEvents(f.CreateAccountWithContracts(account), ignoreFields)
	}
	return f
}

//CreateAccount creates the accounts with the given names
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

//InitializeContractsPrintEvents installs all contracts to the account with name accounts
func (f *GoWithTheFlow) InitializeContractsPrintEvents() *GoWithTheFlow {
	PrintEvents(f.CreateAccountWithAllContracts("accounts"), map[string][]string{})
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

// AddContractWithSignerAsInitArg will deploy a contract with the given name to an account with the same name from wallet.json
func (f *GoWithTheFlow) AddContractWithSignerAsInitArg(accountName string, contractName string) []flow.Event {
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
	tx := AddAccountContractWitwhAuthInit(user.Address, contract)
	events, err := f.performTransaction(tx, &user, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v adding contract %s %+v", emoji.PileOfPoo, contractName, err)
	}
	log.Printf("%v Contract: %s successfully deployed\n", emoji.Scroll, contractName)
	return events
}

// AddAccountContractWitwhAuthInit generates a transaction that deploys a contract to an account.
func AddAccountContractWitwhAuthInit(address flow.Address, contract templates.Contract) *flow.Transaction {
	cadenceName := cadence.NewString(contract.Name)
	cadenceCode := cadence.NewString(contract.SourceHex())

	return flow.NewTransaction().
		SetScript([]byte(addAccountContractTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

const addAccountContractTemplate = `
transaction(name: String, code: String) {
	prepare(signer: AuthAccount) {
		signer.contracts.add(name: name, code: code.decodeHex(), signer)
	}
}`

// UpdateContract will deploy a contract with the given name to an account with the same name from wallet.json
func (f *GoWithTheFlow) UpdateContract(accountName string, contractName string) []flow.Event {
	contractPath := fmt.Sprintf("./contracts/%s.cdc", contractName)
	log.Printf("Updating contract: %s at %s", contractName, contractPath)
	code, err := ioutil.ReadFile(contractPath)
	if err != nil {
		log.Fatalf("%v Could not read contract file from path=%s", emoji.PileOfPoo, contractPath)
	}
	user := f.Accounts[accountName]
	contract := templates.Contract{
		Name:   contractName,
		Source: string(code),
	}
	tx := templates.UpdateAccountContract(user.Address, contract)
	events, err := f.performTransaction(tx, &user, []*GoWithTheFlowAccount{}, []cadence.Value{})
	if err != nil {
		log.Fatalf("%v updateing contract %s %+v", emoji.PileOfPoo, contractName, err)
	}
	log.Printf("%v Contract: %s successfully deployed\n", emoji.Scroll, contractName)
	return events
}
