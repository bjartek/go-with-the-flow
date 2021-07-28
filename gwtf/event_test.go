package gwtf

import (
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvent(t *testing.T) {

	g := NewGoWithTheFlow([]string{"../examples/flow.json"}, "emulator", true, output.NoneLog)

	t.Run("Start argument", func(t *testing.T) {
		ef := g.EventFetcher().Start(100)
		assert.Equal(t, ef.FromIndex, int64(100))
	})

	t.Run("From argument", func(t *testing.T) {
		ef := g.EventFetcher().From(100)
		assert.Equal(t, ef.FromIndex, int64(100))
	})

	t.Run("End argument", func(t *testing.T) {
		ef := g.EventFetcher().End(100)
		assert.Equal(t, ef.EndIndex, uint64(100))
		assert.False(t, ef.EndAtCurrentHeight)
	})

	t.Run("Until argument", func(t *testing.T) {
		ef := g.EventFetcher().Until(100)
		assert.Equal(t, ef.EndIndex, uint64(100))
		assert.False(t, ef.EndAtCurrentHeight)
	})

	t.Run("Until current argument", func(t *testing.T) {
		ef := g.EventFetcher().UntilCurrent()
		assert.Equal(t, ef.EndIndex, uint64(0))
		assert.True(t, ef.EndAtCurrentHeight)
	})

	t.Run("workers argument", func(t *testing.T) {
		ef := g.EventFetcher().Workers(100)
		assert.Equal(t, ef.NumberOfWorkers, 100)
	})

	t.Run("batch size argument", func(t *testing.T) {
		ef := g.EventFetcher().BatchSize(100)
		assert.Equal(t, ef.EventBatchSize, uint64(100))
	})

	t.Run("event ignoring field argument", func(t *testing.T) {
		ef := g.EventFetcher().EventIgnoringFields("foo", []string{"bar", "baz"})
		assert.Equal(t, ef.EventsAndIgnoreFields["foo"], []string{"bar", "baz"})
	})
}
