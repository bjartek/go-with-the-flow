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

//FlowScriptBuilder is a struct to hold information for running a script
type FlowScriptBuilder struct {
	GoWithTheFlow *GoWithTheFlow
	FileName      string
	Arguments     []cadence.Value
}

//ScriptFromFile will start a flow script builder
func (f *GoWithTheFlow) ScriptFromFile(filename string) FlowScriptBuilder {
	return FlowScriptBuilder{
		GoWithTheFlow: f,
		FileName:      filename,
		Arguments:     []cadence.Value{},
	}
}

//AccountArgument add an account as an argument
func (t FlowScriptBuilder) AccountArgument(key string) FlowScriptBuilder {
	f := t.GoWithTheFlow
	address := cadence.BytesToAddress(f.Accounts[key].Address.Bytes())
	return t.Argument(address)
}

//UFix64Argument add an account as an argument
func (t FlowScriptBuilder) UFix64Argument(key string) FlowScriptBuilder {
	amount, err := cadence.NewUFix64(key)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

//StringArgument add an argument to the transaction
func (t FlowScriptBuilder) StringArgument(value string) FlowScriptBuilder {
	t.Arguments = append(t.Arguments, cadence.String(value))
	return t
}

//Argument add an argument to the transaction
func (t FlowScriptBuilder) Argument(value cadence.Value) FlowScriptBuilder {
	t.Arguments = append(t.Arguments, value)
	return t
}

// Run executes a read only script
func (t FlowScriptBuilder) Run() {
	_ = t.RunReturns()
}

// RunReturns executes a read only script
func (t FlowScriptBuilder) RunReturns() cadence.Value {

	f := t.GoWithTheFlow
	c, err := client.New(f.Address, grpc.WithInsecure(), grpc.WithMaxMsgSize(maxGRPCMessageSize))
	if err != nil {
		log.Fatalf("%v Error creating flow client", emoji.PileOfPoo)
	}

	scriptFilePath := fmt.Sprintf("./scripts/%s.cdc", t.FileName)
	code, err := ioutil.ReadFile(scriptFilePath)
	if err != nil {
		log.Fatalf("%v Could not read script file from path=%s", emoji.PileOfPoo, scriptFilePath)
	}

	log.Printf("Arguments %v\n", t.Arguments)
	ctx := context.Background()
	result, err := c.ExecuteScriptAtLatestBlock(ctx, code, t.Arguments)
	if err != nil {
		log.Fatalf("%v Error executing script: %s output %v", emoji.PileOfPoo, t.FileName, err)
	}

	log.Printf("%v Script run from path %s result: %v\n", emoji.Star, scriptFilePath, result)
	return result
}
