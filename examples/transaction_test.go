package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/stretchr/testify/assert"
)

/*
 Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestTransaction(t *testing.T) {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	events, err := g.TransactionFromFile("create_nft_collection").SignProposeAndPayAs("first").RunE()

	assert.NoError(t, err)
	assert.Empty(t, events)

}
