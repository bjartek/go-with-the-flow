package main

import (
	"fmt"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowDevNet()

	eventsFetcher := g.EventFetcher().
		Last(1000).
		Event("A.0b2a3299cc857e29.TopShot.Withdraw")
		//EventIgnoringFields("A.0b2a3299cc857e29.TopShot.Withdraw", []string{"field1", "field"})

	events, err := eventsFetcher.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", events)

	//to send events to a discord eventhook use
	//	message, err := gwtf.NewDiscordWebhook("http://your-webhook-url").SendEventsToWebhook(events)

}
