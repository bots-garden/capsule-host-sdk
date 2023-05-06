package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bots-garden/capsule-host-sdk"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// printHello : print a string to the console
var printHello = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {

	// Extract the position and size of the returned value
	position, length := capsule.UnPackPosSize(stack[0])

	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}
	fmt.Println("👋 hello", string(buffer))

	stack[0] = 0 // return 0
})


// DefineHostFuncPrintHello defines a host function
func DefineHostFuncPrintHello(builder wazero.HostModuleBuilder) {
		// hostPrintHello
		builder.NewFunctionBuilder().
		WithGoModuleFunction(printHello, 
			[]api.ValueType{
				api.ValueTypeI64, // string position + length
			}, 
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostPrintHello")
}
