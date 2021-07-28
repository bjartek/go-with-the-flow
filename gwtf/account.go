package gwtf

import (
	"errors"
	"fmt"
	"github.com/onflow/flow-go-sdk/crypto"
	"log"
	"sort"
)


//CreateAccount creates the accounts with the given names
func (f *GoWithTheFlow) CreateAccounts(saAccountName string) *GoWithTheFlow {

	p := f.State
	signerAccount := p.Accounts().ByName(saAccountName)

	accounts := p.AccountNamesForNetwork(f.Network)
	sort.Strings(accounts)

	f.Logger.Info(fmt.Sprintf("%v\n", accounts))

	for _, accountName := range accounts {
		f.Logger.Info(fmt.Sprintf("Ensuring account with name '%s' is present", accountName))

		account := p.Accounts().ByName(accountName)

		_, err2 := f.Services.Accounts.Get(account.Address())
		if err2 == nil {
			f.Logger.Info("Account is present")
			continue
		}
		a, err := f.Services.Accounts.Create(
			signerAccount,
			[]crypto.PublicKey{account.Key().ToConfig().PrivateKey.PublicKey()},
			[]int{1000},
			account.Key().SigAlgo(),
			account.Key().HashAlgo(),
			[]string{})
		if err != nil {
			log.Fatal(err)
		}
		if account.Address() != a.Address {
			log.Fatal(errors.New(fmt.Sprintf("Configured account address != created address, %s != %s", account.Address(), a.Address)))
		}
		f.Logger.Info("Account created " + a.Address.String())
	}
	return f
}

//InitializeContracts installs all contracts to the account with name accounts
func (f *GoWithTheFlow) InitializeContracts() *GoWithTheFlow {
	f.Logger.Info("Deploying contracts")
	_, err := f.Services.Project.Deploy(f.Network, false)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
