package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	events, err := g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunE()

	assert.NoError(t, err)
	assert.Empty(t, events)

}
