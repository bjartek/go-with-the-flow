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

//Argument add an argument to the transaction
func (t FlowScriptBuilder) Argument(value cadence.Value) FlowScriptBuilder {
	t.Arguments = append(t.Arguments, value)
	return t
}

//StringArgument add a String Argument to the transaction
func (t FlowScriptBuilder) StringArgument(value string) FlowScriptBuilder {
	return t.Argument(cadence.String(value))
}

//BooleanArgument add a Boolean Argument to the transaction
func (t FlowScriptBuilder) BooleanArgument(value bool) FlowScriptBuilder {
	return t.Argument(cadence.NewBool(value))
}

//BytesArgument add a Bytes Argument to the transaction
func (t FlowScriptBuilder) BytesArgument(value []byte) FlowScriptBuilder {
	return t.Argument(cadence.NewBytes(value))
}

//IntArgument add an Int Argument to the transaction
func (t FlowScriptBuilder) IntArgument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt(value))
}

//Int8Argument add an Int8 Argument to the transaction
func (t FlowScriptBuilder) Int8Argument(value int8) FlowScriptBuilder {
	return t.Argument(cadence.NewInt8(value))
}

//Int16Argument add an Int16 Argument to the transaction
func (t FlowScriptBuilder) Int16Argument(value int16) FlowScriptBuilder {
	return t.Argument(cadence.NewInt16(value))
}

//Int32Argument add an Int32 Argument to the transaction
func (t FlowScriptBuilder) Int32Argument(value int32) FlowScriptBuilder {
	return t.Argument(cadence.NewInt32(value))
}

//Int64Argument add an Int64 Argument to the transaction
func (t FlowScriptBuilder) Int64Argument(value int64) FlowScriptBuilder {
	return t.Argument(cadence.NewInt64(value))
}

//Int128Argument add an Int128 Argument to the transaction
func (t FlowScriptBuilder) Int128Argument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt128(value))
}

//Int256Argument add an Int256 Argument to the transaction
func (t FlowScriptBuilder) Int256Argument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt256(value))
}

//UIntArgument add an UInt Argument to the transaction
func (t FlowScriptBuilder) UIntArgument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt(value))
}

//UInt8Argument add an UInt8 Argument to the transaction
func (t FlowScriptBuilder) UInt8Argument(value uint8) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt8(value))
}

//UInt16Argument add an UInt16 Argument to the transaction
func (t FlowScriptBuilder) UInt16Argument(value uint16) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt16(value))
}

//UInt32Argument add an UInt32 Argument to the transaction
func (t FlowScriptBuilder) UInt32Argument(value uint32) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt32(value))
}

//UInt64Argument add an UInt64 Argument to the transaction
func (t FlowScriptBuilder) UInt64Argument(value uint64) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt64(value))
}

//UInt128Argument add an UInt128 Argument to the transaction
func (t FlowScriptBuilder) UInt128Argument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt128(value))
}

//UInt256Argument add an UInt256 Argument to the transaction
func (t FlowScriptBuilder) UInt256Argument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt256(value))
}

//Word8Argument add a Word8 Argument to the transaction
func (t FlowScriptBuilder) Word8Argument(value uint8) FlowScriptBuilder {
	return t.Argument(cadence.NewWord8(value))
}

//Word16Argument add a Word16 Argument to the transaction
func (t FlowScriptBuilder) Word16Argument(value uint16) FlowScriptBuilder {
	return t.Argument(cadence.NewWord16(value))
}

//Word32Argument add a Word32 Argument to the transaction
func (t FlowScriptBuilder) Word32Argument(value uint32) FlowScriptBuilder {
	return t.Argument(cadence.NewWord32(value))
}

//Word64Argument add a Word64 Argument to the transaction
func (t FlowScriptBuilder) Word64Argument(value uint64) FlowScriptBuilder {
	return t.Argument(cadence.NewWord64(value))
}

//Fix64Argument add a Fix64 Argument to the transaction
func (t FlowScriptBuilder) Fix64Argument(value string) FlowScriptBuilder {
	amount, err := cadence.NewFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

//UFix64Argument add a UFix64 Argument to the transaction
func (t FlowScriptBuilder) UFix64Argument(value string) FlowScriptBuilder {
	amount, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
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

	log.Printf("%v Script run from path %s result: %v\n", emoji.Star, scriptFilePath, CadenceValueToJsonString(result))
	return result
}

func (t FlowScriptBuilder) RunReturnsJsonString() string{
	return CadenceValueToJsonString(t.RunReturns())
}

func (t FlowScriptBuilder) RunReturnsInterface() interface{}{
	return CadenceValueToInterface(t.RunReturns())
}
