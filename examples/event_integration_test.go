package main

import (
	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
)

func TestEvents(t *testing.T) {

	t.Run("Fetch last events", func(t *testing.T) {
		g := gwtf.NewTestingEmulator()
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3)

		ev, err := g.EventFetcher().Last(2).Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").Run()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

	t.Run("Fetch last write progress file", func(t *testing.T) {
		g := gwtf.NewTestingEmulator()
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3)

		ev, err := g.EventFetcher().Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").TrackProgressIn("progress").Run()
		defer os.Remove("progress")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

	t.Run("Fetch last write progress file that exists", func(t *testing.T) {

		err := os.WriteFile("progress", []byte("1"), fs.ModePerm)
		assert.NoError(t, err)

		g := gwtf.NewTestingEmulator()
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3)

		ev, err := g.EventFetcher().Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").TrackProgressIn("progress").Run()
		defer os.Remove("progress")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

}
