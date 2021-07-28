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



}
