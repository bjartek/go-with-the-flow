package main

import (
	"fmt"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

func main() {

	g := gwtf.NewGoWithTheFlowMainNet()

	eventsFetcher := g.EventFetcher().
		Last(1000).
		Event("A.0b2a3299cc857e29.TopShot.Withdraw")

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
