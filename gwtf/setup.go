package gwtf

import (
	"context"
	"fmt"
	"log"

	"github.com/enescakir/emoji"
	"github.com/mitchellh/go-homedir"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/pkg/errors"
)

// DiscordWebhook stores information about a webhook
type DiscordWebhook struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Wait  bool   `json:"wait"`
}

// GoWithTheFlow Entire configuration to work with Go With the Flow
type GoWithTheFlow struct {
	Service  GoWithTheFlowAccount
	Accounts map[string]GoWithTheFlowAccount
	WebHooks map[string]DiscordWebhook
	Address  string
	Gas      uint64
}

// GoWithTheFlowAccount represents an account for flow with resolves types
type GoWithTheFlowAccount struct {
	Address    flow.Address
	SigAlgo    crypto.SignatureAlgorithm
	HashAlgo   crypto.HashAlgorithm
	PrivateKey crypto.PrivateKey

	//These three are set the first time they are accessed
	Signer  *crypto.InMemorySigner
	Account *flow.Account
	Key     *flow.AccountKey
}

//NewGoWithTheFlowEmulator create a new client
func NewGoWithTheFlowEmulator() *GoWithTheFlow {
	return NewGoWithTheFlow("./flow.json")
}

// NewGoWithTheFlowDevNet setup dev like in https://www.notion.so/Accessing-Flow-Devnet-ad35623797de48c08d8b88102ea38131
func NewGoWithTheFlowDevNet() *GoWithTheFlow {
	flowConfigFile, err := homedir.Expand("~/.flow-dev.json")
	if err != nil {
		log.Fatalf("%v error %v", emoji.PileOfPoo, err)
	}

	return NewGoWithTheFlow(flowConfigFile)
}

// NewGoWithTheFlow with custom file panic on error
func NewGoWithTheFlow(filename string) *GoWithTheFlow {
	gwtf, err := NewGoWithTheFlowError(filename)
	if err != nil {
		log.Fatalf("%v error %+v", emoji.PileOfPoo, err)
	}
	return gwtf
}

// NewGoWithTheFlowError creates a new local go with the flow client
func NewGoWithTheFlowError(fileName string) (*GoWithTheFlow, error) {

	config, err := NewRawFlowConfig(fileName)
	if err != nil {
		return nil, err
	}

	if _, ok := config.Accounts["emulator-account"]; !ok {
		return nil, errors.New("Could not find emulator-account account in flow.json")
	}

	//loop over all the

	address := "127.0.0.1:3569"
	if config.Address != "" {
		address = config.Address
	}

	gas := uint64(1000)
	if config.GasLimit != 0 {
		gas = config.GasLimit
	}

	rawAccounts := config.Accounts
	for account, key := range config.EmulatorAccounts {
		rawAccounts[account] = RawAccount{
			Address: key,
			Keys:    "d5457a187e9642a8e49d4032b3b4f85c92da7202c79681d9302c6e444e7033a8",
		}
	}

	var accounts = map[string]GoWithTheFlowAccount{}
	var serviceAccount GoWithTheFlowAccount
	for name, rawAccount := range rawAccounts {
		//TODO: Only support the default for now
		sigAlgo := crypto.StringToSignatureAlgorithm("ECDSA_P256")
		hashAlgo := crypto.StringToHashAlgorithm("SHA3_256")
		privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, rawAccount.Keys)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Could not decode private key for %s", name))
		}

		address := flow.HexToAddress(rawAccount.Address)

		gwtfAccount := GoWithTheFlowAccount{
			Address:    address,
			SigAlgo:    sigAlgo,
			HashAlgo:   hashAlgo,
			PrivateKey: privateKey,
		}
		if name == "emulator-account" {
			serviceAccount = gwtfAccount
		} else {
			accounts[name] = gwtfAccount
		}
	}

	return &GoWithTheFlow{
		Address:  address,
		Gas:      gas,
		Service:  serviceAccount,
		Accounts: accounts,
		WebHooks: config.Webhooks,
	}, nil

}

//EnrichWithAccountSignerAndKey enriches and Account
func (a *GoWithTheFlowAccount) EnrichWithAccountSignerAndKey(c *client.Client) (*GoWithTheFlowAccount, error) {
	ctx := context.Background()
	serviceAccount, err := c.GetAccount(ctx, a.Address)
	if err != nil {
		return nil, err
	}
	serviceAccountKey := serviceAccount.Keys[1]
	a.Account = serviceAccount
	signer := crypto.NewInMemorySigner(a.PrivateKey, serviceAccountKey.HashAlgo)
	a.Signer = &signer
	a.Key = serviceAccountKey

	return a, nil
}

//NewAccountKey creates a NewFlowAccountKey
func (a *GoWithTheFlowAccount) NewAccountKey() *flow.AccountKey {
	return flow.NewAccountKey().
		SetPublicKey(a.PrivateKey.PublicKey()).
		SetSigAlgo(a.SigAlgo).
		SetHashAlgo(a.HashAlgo).
		SetWeight(flow.AccountKeyWeightThreshold)
}

//FindAddress finds an candence.Address value from a given key in your wallet
func (f *GoWithTheFlow) FindAddress(key string) cadence.Address {
	return cadence.BytesToAddress(f.Accounts[key].Address.Bytes())
}
