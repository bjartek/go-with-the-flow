package gwtf

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

//FetchEvents fetches events for the given query and formats them
func (f *GoWithTheFlow) FetchEvents(query client.EventRangeQuery) ([]*FormatedEvent, error) {

	ctx := context.Background()

	c, err := client.New(f.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	formatedEvents := []*FormatedEvent{}
	blockEvents, err := c.GetEventsForHeightRange(ctx, query)
	for _, blockEvent := range blockEvents {
		var time time.Time
		if len(blockEvent.Events) > 0 {
			header, err2 := c.GetBlockHeaderByHeight(ctx, blockEvent.Height)
			if err2 != nil {
				return nil, err2
			}
			time = header.Timestamp
		}
		for _, event := range blockEvent.Events {
			ev := ParseEvent(event, blockEvent.Height, time)
			formatedEvents = append(formatedEvents, ev)
		}
	}
	return formatedEvents, nil
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

//PrintEvents prints th events
func PrintEvents(events []flow.Event) {
	if len(events) > 0 {
		fmt.Println("EVENTS")
		fmt.Println("======")
	}
	for _, event := range events {
		ev := ParseEvent(event, uint64(0), time.Now())
		prettyJSON, err := json.MarshalIndent(ev, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	}
	if len(events) > 0 {
		fmt.Println("======")
	}
}

//ParseEvent parses a flow event into a more terse representation
func ParseEvent(event flow.Event, blockHeight uint64, time time.Time) *FormatedEvent {
	var fieldNames []string
	for _, eventTypeFields := range event.Value.EventType.Fields {
		fieldNames = append(fieldNames, eventTypeFields.Identifier)
	}
	finalFields := map[string]string{}
	for id, field := range event.Value.Fields {

		name := fieldNames[id]
		value := fmt.Sprintf("%+v", field)
		var fieldValue string
		if strings.Contains(value, "Values") {
			fieldValue = between(value, "Values:[", "]}")
		} else {
			fieldValue = value
		}
		finalFields[name] = fieldValue
	}
	return &FormatedEvent{
		Name:        event.Value.EventType.Identifier,
		Fields:      finalFields,
		BlockHeight: blockHeight,
		Time:        time,
	}
}

//FormatedEvent event in a more condensed formated form
type FormatedEvent struct {
	Name        string            `json:"name"`
	BlockHeight uint64            `json:"blockHeight,omitempty"`
	Time        time.Time         `json:"time,omitempty"`
	Fields      map[string]string `json:"fields"`
}
