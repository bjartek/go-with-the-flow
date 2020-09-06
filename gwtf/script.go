package gwtf

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

// RunScript executes a read only script with a given filename on the blockchain
func (f *GoWithTheFlow) RunScript(filename string, arguments ...cadence.Value) {
	_ = f.RunScriptReturns(filename, arguments...)
}

// RunScriptReturns executes a read only script with a given filename on the blockchain
func (f *GoWithTheFlow) RunScriptReturns(filename string, arguments ...cadence.Value) cadence.Value {

	c, err := client.New(f.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v Error creating flow client", emoji.PileOfPoo)
	}

	scriptFilePath := fmt.Sprintf("./scripts/%s.cdc", filename)
	code, err := ioutil.ReadFile(scriptFilePath)
	if err != nil {
		log.Fatalf("%v Could not read script file from path=%s", emoji.PileOfPoo, scriptFilePath)
	}

	log.Printf("Arguments %v\n", arguments)
	ctx := context.Background()
	result, err := c.ExecuteScriptAtLatestBlock(ctx, code, arguments)
	if err != nil {
		log.Fatalf("%v Error executing script: %s output %v", emoji.PileOfPoo, filename, err)
	}

	log.Printf("%v Script run from path %s result: %v\n", emoji.Star, scriptFilePath, result)
	return result
}
