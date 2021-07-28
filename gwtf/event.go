package gwtf

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/pkg/errors"
)

//EventFetcherBuilder builder to hold info about eventhook context
type EventFetcherBuilder struct {
	GoWithTheFlow         *GoWithTheFlow
	EventsAndIgnoreFields map[string][]string
	FromIndex             int64
	EndAtCurrentHeight    bool
	EndIndex              uint64
	ProgressFile          string
	context.Context
	NumberOfWorkers int
	EventBatchSize  uint64
}

//EventFetcher create an event fetcher builder
func (f *GoWithTheFlow) EventFetcher() EventFetcherBuilder {
	return EventFetcherBuilder{
		GoWithTheFlow:         f,
		EventsAndIgnoreFields: map[string][]string{},
		EndAtCurrentHeight:    true,
		FromIndex:             -10,
		ProgressFile:          "",
		EventBatchSize:        250,
		NumberOfWorkers:       20,
	}
}

//Workers sets the number of workers
func (e EventFetcherBuilder) Workers(workers int) EventFetcherBuilder {
	e.NumberOfWorkers = workers
	return e
}

//BatchSize sets the size of a batch
func (e EventFetcherBuilder) BatchSize(batchSize uint64) EventFetcherBuilder {
	e.EventBatchSize = batchSize
	return e
}

//Event fetches and Events and all its fields
func (e EventFetcherBuilder) Event(eventName string) EventFetcherBuilder {
	e.EventsAndIgnoreFields[eventName] = []string{}
	return e
}

//EventIgnoringFields fetch event and ignore the specified fields
func (e EventFetcherBuilder) EventIgnoringFields(eventName string, ignoreFields []string) EventFetcherBuilder {
	e.EventsAndIgnoreFields[eventName] = ignoreFields
	return e
}

//Start specify what blockHeight to fetch starting atm. This can be negative related to end/until
func (e EventFetcherBuilder) Start(blockHeight int64) EventFetcherBuilder {
	e.FromIndex = blockHeight
	return e
}

//From specify what blockHeight to fetch from. This can be negative related to end.
func (e EventFetcherBuilder) From(blockHeight int64) EventFetcherBuilder {
	e.FromIndex = blockHeight
	return e
}

//End specify what index to end at
func (e EventFetcherBuilder) End(blockHeight uint64) EventFetcherBuilder {
	e.EndIndex = blockHeight
	e.EndAtCurrentHeight = false
	return e
}

//Last fetch events from the number last blocks
func (e EventFetcherBuilder) Last(number uint64) EventFetcherBuilder {
	e.EndAtCurrentHeight = true
	e.FromIndex = -int64(number)
	return e
}

//Until specify what index to end at
func (e EventFetcherBuilder) Until(blockHeight uint64) EventFetcherBuilder {
	e.EndIndex = blockHeight
	e.EndAtCurrentHeight = false
	return e
}

//UntilCurrent Specify to fetch events until the current Block
func (e EventFetcherBuilder) UntilCurrent() EventFetcherBuilder {
	e.EndAtCurrentHeight = true
	e.EndIndex = 0
	return e
}

//TrackProgressIn Specify a file to store progress in
func (e EventFetcherBuilder) TrackProgressIn(fileName string) EventFetcherBuilder {
	e.ProgressFile = fileName
	e.EndIndex = 0
	e.FromIndex = 0
	e.EndAtCurrentHeight = true
	return e
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func writeProgressToFile(fileName string, blockHeight uint64) error {
	err := ioutil.WriteFile(fileName, []byte(fmt.Sprintf("%d", blockHeight)), 0644)
	if err != nil {
		return errors.Wrap(err, "Could not create initial progress file")
	}
	return nil
}

func readProgressFromFile(fileName string) (int64, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, errors.Wrap(err, "ProgressFile is not valid")
	}

	stringValue := strings.TrimSpace(string(dat))

	return strconv.ParseInt(stringValue, 10, 64)

}

//RunAndSendToWebhook runs the eventFetcher and sends the event to a webhook with the provided url
func (e EventFetcherBuilder) RunAndSendToWebhook(url string) (*discordgo.Message, error) {

	events, err := e.Run()
	if err != nil {
		return nil, err
	}

	return NewDiscordWebhook(url).SendEventsToWebhook(events)
}

//Run runs the eventfetcher returning events or an error
func (e EventFetcherBuilder) Run() ([]*FormatedEvent, error) {

	//if we have a progress file read the value from it and set it as oldHeight
	if e.ProgressFile != "" {
		//TODO if file does not exist that is OK

		present, err := exists(e.ProgressFile)
		if err != nil {
			return nil, err
		}
		if !present {
			err := writeProgressToFile(e.ProgressFile, 0)
			if err != nil {
				return nil, errors.Wrap(err, "Could not create initial progress file")
			}
			e.FromIndex = 0
		} else {
			oldHeight, err := readProgressFromFile(e.ProgressFile)
			if err != nil {
				return nil, errors.Wrap(err, "Could not parse progress file as block height")
			}
			e.FromIndex = oldHeight
		}
	}

	endIndex := e.EndIndex
	if e.EndAtCurrentHeight {
		blockHeight, err := e.GoWithTheFlow.Services.Blocks.GetLatestBlockHeight()
		if err != nil {
			return nil, err
		}
		endIndex = blockHeight
	}

	fromIndex := e.FromIndex
	//if we have a negative fromIndex is is relative to endIndex
	if e.FromIndex <= 0 {
		fromIndex = int64(endIndex) + e.FromIndex
	}

	if fromIndex < 0 {
		return nil, errors.New("FromIndex is negative")
	}

	e.GoWithTheFlow.Logger.Info(fmt.Sprintf("Fetching events from %d to %d", fromIndex, endIndex))

	var events []string
	for key := range e.EventsAndIgnoreFields {
		events = append(events, key)
	}

	blockEvents, err := e.GoWithTheFlow.Services.Events.Get(events, uint64(fromIndex), endIndex, e.EventBatchSize, e.NumberOfWorkers)
	if err != nil {
		return nil, err
	}

	formatedEvents := FormatEvents(blockEvents, e.EventsAndIgnoreFields)

	if e.ProgressFile != "" {
		err := writeProgressToFile(e.ProgressFile, endIndex+1)
		if err != nil {
			return nil, errors.Wrap(err, "Could not write progress to file")
		}
	}
	sort.Slice(formatedEvents, func(i, j int) bool {
		return formatedEvents[i].BlockHeight < formatedEvents[j].BlockHeight
	})

	return formatedEvents, nil

}

//PrintEvents prints th events, ignoring fields specified for the given event typeID
func PrintEvents(events []flow.Event, ignoreFields map[string][]string) {
	if len(events) > 0 {
		fmt.Println("EVENTS")
		fmt.Println("======")
	}
	for _, event := range events {

		eventType := string(event.Value.EventType.Location.ID())
		ignoreFieldsForType := ignoreFields[eventType]
		ev := ParseEvent(event, uint64(0), time.Now(), ignoreFieldsForType)
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

//FormatEvents
func FormatEvents(blockEvents []client.BlockEvents, ignoreFields map[string][]string) []*FormatedEvent {
	var events []*FormatedEvent
	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			ev := ParseEvent(event, blockEvent.Height, blockEvent.BlockTimestamp, ignoreFields[event.Type])
			events = append(events, ev)
		}
	}
	return events
}

//ParseEvent parses a flow event into a more terse representation
func ParseEvent(event flow.Event, blockHeight uint64, time time.Time, ignoreFields []string) *FormatedEvent {
	var fieldNames []string
	for _, eventTypeFields := range event.Value.EventType.Fields {
		fieldNames = append(fieldNames, eventTypeFields.Identifier)
	}
	finalFields := map[string]interface{}{}
	for id, field := range event.Value.Fields {
		skip := false
		name := fieldNames[id]
		for _, ignoreField := range ignoreFields {
			if ignoreField == name {
				skip = true
			}
		}
		if skip {
			continue
		}

		finalFields[name] = CadenceValueToInterface(field)
	}
	return &FormatedEvent{
		Name:        event.Type,
		Fields:      finalFields,
		BlockHeight: blockHeight,
		Time:        time,
	}
}

//FormatedEvent event in a more condensed formated form
type FormatedEvent struct {
	Name        string                 `json:"name"`
	BlockHeight uint64                 `json:"blockHeight,omitempty"`
	Time        time.Time              `json:"time,omitempty"`
	Fields      map[string]interface{} `json:"fields"`
}

func NewTestEvent(name string, fields map[string]interface{}) *FormatedEvent {
	return &FormatedEvent{
		Name:        name,
		BlockHeight: 0,
		Time:        time.Unix(0, 0),
		Fields:      fields,
	}
}

//String pretty print an event as a String
func (e FormatedEvent) String() string {
	j, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(j)
}
