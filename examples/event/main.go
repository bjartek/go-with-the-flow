package main

import (
	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	g := gwtf.NewGoWithTheFlowEmulator()

	event, err := g.SendEventsTo("gwtf").
		From(0).
		Until(6).
		EventIgnoringFields("flow.AccountCodeUpdated", []string{"codeHash"}).
		EventIgnoringFields("flow.AccountKeyAdded", []string{"publicKey"}).
		Run()
	if err != nil {
		panic(err)
	}
	spew.Dump(event)

}
