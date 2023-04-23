package capsule

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// logString : print a string to the console
var logString = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {

	// Extract the position and size of the returned value
	position, length := UnPackPosSize(stack[0])

	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}
	fmt.Println(time.Now(), ":", string(buffer))

	stack[0] = 0 // return 0
})


// DefineHostFuncLog defines a host function
func DefineHostFuncLog(builder wazero.HostModuleBuilder) {
		// hostLogString
		builder.NewFunctionBuilder().
		WithGoModuleFunction(logString, 
			[]api.ValueType{
				api.ValueTypeI64, // string position + length
			}, 
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostLogString")
}
