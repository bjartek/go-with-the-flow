package gwtf

import (
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorsInAccountCreation(t *testing.T) {

	t.Run("Should give error on wrong saAccount name", func(t *testing.T) {
		g := NewGoWithTheFlow([]string{"../examples/flow.json"}, "emulator", true, output.NoneLog)
		_, err := g.CreateAccountsE("foobar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find account with name foobar")
	})

	t.Run("Should give erro on wrong account name", func(t *testing.T) {
		_, err := NewGoWithTheFlowError([]string{"fixtures/invalid_account_in_deployment.json"}, "emulator", true, output.NoneLog)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deployment contains nonexisting account emulator-firs")
	})

	t.Run("Should create inmemory emulator client", func(t *testing.T) {
		g := NewGoWithTheFlowInMemoryEmulator()
		assert.Equal(t, g.Network, "emulator")
	})

	t.Run("Should create local emulator client", func(t *testing.T) {
		g := NewGoWithTheFlowEmulator()
		assert.Equal(t, g.Network, "emulator")
	})

	t.Run("Should create testnet client", func(t *testing.T) {
		g := NewGoWithTheFlowDevNet()
		assert.Equal(t, g.Network, "testnet")
	})

	t.Run("Should create testnet client with for network metdho", func(t *testing.T) {
		g := NewGoWithTheFlowForNetwork("testnet")
		assert.Equal(t, g.Network, "testnet")
	})

	t.Run("Should create mainnet client", func(t *testing.T) {
		g := NewGoWithTheFlowMainNet()
		assert.Equal(t, g.Network, "mainnet")
		assert.True(t, g.PrependNetworkToAccountNames)
		g = g.DoNotPrependNetworkToAccountNames()
		assert.False(t, g.PrependNetworkToAccountNames)

	})

}
