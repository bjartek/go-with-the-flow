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
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

//EventHookBuilder builder to hold info about eventhook context
type EventHookBuilder struct {
	GoWithTheFlow         *GoWithTheFlow
	WebhookName           string
	EventsAndIgnoreFields map[string][]string
	FromIndex             int64
	EndAtCurrentHeight    bool
	EndIndex              uint64
	ProgressFile          string
}

// SendEventsTo starts a event hook builder
func (f *GoWithTheFlow) SendEventsTo(eventHookName string) EventHookBuilder {
	return EventHookBuilder{
		GoWithTheFlow:         f,
		WebhookName:           eventHookName,
		EventsAndIgnoreFields: map[string][]string{},
		EndAtCurrentHeight:    true,
		FromIndex:             -10,
		ProgressFile:          "",
	}
}

// Event fetches and Events and all its fields
func (e EventHookBuilder) Event(eventName string) EventHookBuilder {
	e.EventsAndIgnoreFields[eventName] = []string{}
	return e
}

//EventIgnoringFields fetch event and ignore the specified fields
func (e EventHookBuilder) EventIgnoringFields(eventName string, ignoreFields []string) EventHookBuilder {
	e.EventsAndIgnoreFields[eventName] = ignoreFields
	return e
}

//Start specify what blockHeight to fetch starting atm. This can be negative related to end/until
func (e EventHookBuilder) Start(blockHeight int64) EventHookBuilder {
	e.FromIndex = blockHeight
	return e
}

//From specify what blockHeight to fetch from. This can be negative related to end.
func (e EventHookBuilder) From(blockHeight int64) EventHookBuilder {
	e.FromIndex = blockHeight
	return e
}

//End specify what index to end at
func (e EventHookBuilder) End(blockHeight uint64) EventHookBuilder {
	e.EndIndex = blockHeight
	e.EndAtCurrentHeight = false
	return e
}

//Last fetch events from the number last blocks
func (e EventHookBuilder) Last(number uint64) EventHookBuilder {
	e.EndAtCurrentHeight = true
	e.FromIndex = -int64(number)
	return e
}

//Until specify what index to end at
func (e EventHookBuilder) Until(blockHeight uint64) EventHookBuilder {
	e.EndIndex = blockHeight
	e.EndAtCurrentHeight = false
	return e
}

//UntilCurrent Specify to fetch events until the current Block
func (e EventHookBuilder) UntilCurrent() EventHookBuilder {
	e.EndAtCurrentHeight = true
	e.EndIndex = 0
	return e
}

//TrackProgressIn Specify a file to store progress in
func (e EventHookBuilder) TrackProgressIn(fileName string) EventHookBuilder {
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

// Run the eventHook flow
func (e EventHookBuilder) Run() (*discordgo.Message, error) {

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

	eventHook, ok := e.GoWithTheFlow.WebHooks[e.WebhookName]
	if !ok {
		return nil, errors.New("Could not find webhook with name " + e.WebhookName)
	}

	ctx := context.Background()
	c, err := client.New(e.GoWithTheFlow.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	endIndex := e.EndIndex
	if e.EndAtCurrentHeight {
		header, err := c.GetLatestBlockHeader(ctx, true)
		if err != nil {
			return nil, err
		}
		endIndex = header.Height
	}

	fromIndex := e.FromIndex
	//if we have a negative fromIndex is is relative to endIndex
	if e.FromIndex <= 0 {
		fromIndex = int64(endIndex) + e.FromIndex
	}

	if fromIndex < 0 {
		return nil, errors.New("FromIndex is negative")
	}

	formatedEvents := []*FormatedEvent{}
	for contract, ignoreFields := range e.EventsAndIgnoreFields {
		events, err := fetchEvents(ctx, c,
			client.EventRangeQuery{
				Type:        contract,
				StartHeight: uint64(fromIndex),
				EndHeight:   endIndex,
			}, ignoreFields)
		if err != nil {
			return nil, err
		}
		formatedEvents = append(formatedEvents, events...)
	}

	sort.Slice(formatedEvents, func(i, j int) bool {
		return formatedEvents[i].BlockHeight < formatedEvents[j].BlockHeight
	})

	if len(formatedEvents) == 0 {
		return nil, nil
	}

	if e.ProgressFile != "" {
		err := writeProgressToFile(e.ProgressFile, endIndex)
		if err != nil {
			return nil, errors.Wrap(err, "Could not write progress to file")
		}

	}
	return eventHook.SendEventsToWebhook(formatedEvents)

}

// SendEventsToWebhook Sends events to the webhook with the given name from flow.json
func (dw DiscordWebhook) SendEventsToWebhook(events []*FormatedEvent) (*discordgo.Message, error) {

	discord, err := discordgo.New()
	if err != nil {
		return nil, err
	}

	status, err := discord.WebhookExecute(
		dw.ID,
		dw.Token,
		dw.Wait,
		EventsToWebhookParams(events))

	if err != nil {
		return nil, err
	}
	return status, nil
}

//FetchEvents fetches events for the given query and formats them
func fetchEvents(ctx context.Context, c *client.Client, query client.EventRangeQuery, ignoreFields []string) ([]*FormatedEvent, error) {

	formatedEvents := []*FormatedEvent{}
	blockEvents, err := c.GetEventsForHeightRange(ctx, query)
	if err != nil {
		return nil, err
	}
	for _, blockEvent := range blockEvents {
		for _, event := range blockEvent.Events {
			ev := ParseEvent(event, blockEvent.Height, blockEvent.BlockTimestamp, ignoreFields)
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

//ParseEvent parses a flow event into a more terse representation
func ParseEvent(event flow.Event, blockHeight uint64, time time.Time, ignoreFields []string) *FormatedEvent {
	var fieldNames []string
	for _, eventTypeFields := range event.Value.EventType.Fields {
		fieldNames = append(fieldNames, eventTypeFields.Identifier)
	}
	finalFields := map[string]string{}
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

		finalFields[name] = fmt.Sprintf("%v", field)
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
	Name        string            `json:"name"`
	BlockHeight uint64            `json:"blockHeight,omitempty"`
	Time        time.Time         `json:"time,omitempty"`
	Fields      map[string]string `json:"fields"`
}

//just copied this from here https://github.com/onflow/cadence/blob/master/encoding/json/encode.go
func encodeUFix64(v uint64) string {
	integer := v / sema.Fix64Factor
	fraction := v % sema.Fix64Factor

	return fmt.Sprintf(
		"%d.%08d",
		integer,
		fraction,
	)
}

func encodeFix64(v int64) string {
	integer := v / sema.Fix64Factor
	fraction := v % sema.Fix64Factor

	negative := fraction < 0

	var builder strings.Builder

	if negative {
		fraction = -fraction
		if integer == 0 {
			builder.WriteRune('-')
		}
	}

	builder.WriteString(fmt.Sprintf(
		"%d.%08d",
		integer,
		fraction,
	))

	return builder.String()
}
