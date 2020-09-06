package main

import "github.com/bjartek/go-with-the-flow/gwtf"

func main() {

	gwtf := gwtf.NewGoWithTheFlowEmulator()

	gwtf.ScriptFromFile("block").Run()

}
