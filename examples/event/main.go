package main

import (
	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	g := gwtf.NewGoWithTheFlowEmulator()

	eventBuilder := g.SendEventsTo("gwtf").
		TrackProgressIn("/tmp/gwtf_events").
		EventIgnoringFields("flow.AccountCodeUpdated", []string{"codeHash"}).
		EventIgnoringFields("flow.AccountKeyAdded", []string{"publicKey"})

	err, event := eventBuilder.Run()
	if err != nil {
		panic(err)
	}
	spew.Dump(event)

	g.CreateAccountPrintEvents("first")

	err, event = eventBuilder.Run()
	if err != nil {
		panic(err)
	}
	spew.Dump(event)

}
