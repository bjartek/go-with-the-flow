package gwtf

import (
	"fmt"
	"log"
	"sort"

	"github.com/onflow/flow-go-sdk/crypto"
)

// CreateAccounts ensures that all accounts present in the deployment block for the given network is present
func (f *GoWithTheFlow) CreateAccounts(saAccountName string) *GoWithTheFlow {
	p := f.State
	signerAccount, err := p.Accounts().ByName(saAccountName)
	if err != nil {
		log.Fatal(err)
	}

	accounts := p.AccountNamesForNetwork(f.Network)
	sort.Strings(accounts)

	log.Printf("%v\n", accounts)

	for _, accountName := range accounts {
		log.Println(fmt.Sprintf("Ensuring account with name '%s' is present", accountName))

		account, err := p.Accounts().ByName(accountName)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := f.Services.Accounts.Get(account.Address()); err == nil {
			log.Println("Account is present")
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
			log.Fatal(fmt.Errorf("configured account address != created address, %s != %s", account.Address(), a.Address))
		}

		log.Println("Account created " + a.Address.String())
	}
	return f
}

// InitializeContracts installs all contracts in the deployment block for the configured network
func (f *GoWithTheFlow) InitializeContracts() *GoWithTheFlow {
	log.Println("Deploying contracts")
	if _, err := f.Services.Project.Deploy(f.Network, false); err != nil {
		log.Fatal(err)
	}

	return f
}
