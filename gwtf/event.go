package gwtf

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jsoncdc "github.com/onflow/cadence/encoding/json"
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
			ev, err := ParseEvent(event)
			if err != nil {
				return nil, err
			}
			ev.BlockHeight = blockEvent.Height
			ev.Time = time
			formatedEvents = append(formatedEvents, ev)
		}
	}
	return formatedEvents, nil
}

//Not really happy with this code, would love something simpler
//ParseEvent parses a flow.Event into a more condensed form without type info for viewing
func ParseEvent(ev flow.Event) (*FormatedEvent, error) {
	encodedValue, err := jsoncdc.Encode(ev.Value)
	if err != nil {
		return nil, err
	}

	var obj rawEvent
	if err := json.Unmarshal(encodedValue, &obj); err != nil {
		return nil, err
	}

	fields := map[string]string{}
	for _, field := range obj.Value.Fields {
		val := field.Value.Value
		switch val.(type) {
		case string:
			fields[field.Name] = val.(string)
		case []interface{}:
			f := []string{}
			for _, valField := range val.([]interface{}) {
				v := valField.(map[string]interface{})
				f = append(f, v["value"].(string))
			}
			fields[field.Name] = strings.Join(f, ",")
		default:
			fields[field.Name] = fmt.Sprintf("%s", val)
		}
	}

	return &FormatedEvent{
		Name:   obj.Value.ID,
		Fields: fields,
	}, nil
}

//FormatedEvent event in a more condensed formated form
type FormatedEvent struct {
	Name        string            `json:"name"`
	BlockHeight uint64            `json:"blockHeight,omitempty"`
	Time        time.Time         `json:"time,omitempty"`
	Fields      map[string]string `json:"fields"`
}

type rawEvent struct {
	Type  string `json:"type"`
	Value struct {
		ID     string `json:"id"`
		Fields []struct {
			Name  string    `json:"name"`
			Value typeValue `json:"value"`
		} `json:"fields"`
	} `json:"value"`
}

type typeValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
