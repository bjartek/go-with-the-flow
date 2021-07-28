package gwtf

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TransactionResult struct {
	Err     error
	Events  []*FormatedEvent
	Testing *testing.T
}

func (f FlowTransactionBuilder) Test(t *testing.T) TransactionResult {
	events, err := f.RunE()
	var formattedEvents []*FormatedEvent
	for _, event := range events {
		ev := ParseEvent(event, uint64(0), time.Unix(0, 0), []string{})
		formattedEvents = append(formattedEvents, ev)
	}
	return TransactionResult{
		Err:     err,
		Events:  formattedEvents,
		Testing: t,
	}
}

func(t TransactionResult) AssertFailure(msg string) TransactionResult {
	assert.Error(t.Testing, t.Err)
	assert.Contains(t.Testing, t.Err.Error(), msg)
	return t
}
func (t TransactionResult) AssertSuccess() TransactionResult {
	assert.NoError(t.Testing, t.Err)
	return t
}

func (t TransactionResult) AssertEventCount(number int) TransactionResult {
	assert.Equal(t.Testing, len(t.Events), number)
	return t

}
func (t TransactionResult) AssertNoEvents() TransactionResult {
	assert.Empty(t.Testing, t.Err)
	return t
}

func (t TransactionResult) AssertEmitEventName(event ...string) TransactionResult {
	var eventNames []string
	for _, fe := range t.Events {
		eventNames = append(eventNames, fe.Name)
	}

	for _, ev := range event {
		assert.Contains(t.Testing, eventNames, ev)
	}
	return t
}

func (t TransactionResult) AssertEmitEvent(event ...*FormatedEvent) TransactionResult {
	for _, ev := range event {
		assert.Contains(t.Testing, t.Events, ev)
	}
	return t
}

func (t TransactionResult) AssertDebugLog(message ...string) TransactionResult {
	var logMessages []interface{}
	for _, fe := range t.Events {
		if strings.HasSuffix(fe.Name, "Debug.Log") {
			logMessages = append(logMessages, fe.Fields["msg"])
		}
	}

	for _, ev := range message {
		assert.Contains(t.Testing, logMessages, ev)
	}
	return t
}
