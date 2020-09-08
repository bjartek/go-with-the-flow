package gwtf

import (
	"encoding/json"
	"fmt"
	"strings"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk"
)

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
			fmt.Println(field.Name)
		}
	}

	return &FormatedEvent{
		Name:   obj.Value.ID,
		Fields: fields,
	}, nil
}

//FormatedEvent event in a more condensed formated form
type FormatedEvent struct {
	Name   string            `json:"name"`
	Fields map[string]string `json:"fields"`
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
