package main

import (
	"fmt"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowMainNet()

	eventsFetcher := g.EventFetcher().
		Start(14415071).
		UntilCurrent().
		Event("A.d796ff17107bbff6.Versus.Bid")

	events, err := eventsFetcher.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", events)

	/*
		//if you add an eventHook in discord to the flow.json in the example folder you can use this to send the events to a discord channel
		msg, err := eventsFetcher.SendEventsToWebhook("gwtf")
		if err != nil {
			panic(err)
		}
		fmt.Println("send message with id", msg.ID)
	*/

}
