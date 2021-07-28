package main

import (
	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScript(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()
	t.Parallel()

	t.Run("Raw account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").RawAccountArgument("0x1cf0e2f2f715450").RunReturnsInterface()
		assert.Equal(t, value, "0x1cf0e2f2f715450")
	})

	t.Run("Raw account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").AccountArgument("first").RunReturnsInterface()
		assert.Equal(t, value, "0x1cf0e2f2f715450")
	})

}