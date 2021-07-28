package main

import (
	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupIntegration(t *testing.T) {

	t.Run("Should create inmemory emulator client", func(t *testing.T) {
		g := gwtf.NewGoWithTheFlowInMemoryEmulator()
		assert.Equal(t, g.Network, "emulator")
	})

	t.Run("Should create local emulator client", func(t *testing.T) {
		g := gwtf.NewGoWithTheFlowEmulator()
		assert.Equal(t, g.Network, "emulator")
	})

	t.Run("Should create testnet client", func(t *testing.T) {
		g := gwtf.NewGoWithTheFlowDevNet()
		assert.Equal(t, g.Network, "testnet")
	})

	t.Run("Should create testnet client with for network metdho", func(t *testing.T) {
		g := gwtf.NewGoWithTheFlowForNetwork("testnet")
		assert.Equal(t, g.Network, "testnet")
	})

	t.Run("Should create mainnet client", func(t *testing.T) {
		g := gwtf.NewGoWithTheFlowMainNet()
		assert.Equal(t, g.Network, "mainnet")
		assert.True(t, g.PrependNetworkToAccountNames)
		g = g.DoNotPrependNetworkToAccountNames()
		assert.False(t, g.PrependNetworkToAccountNames)

	})
}