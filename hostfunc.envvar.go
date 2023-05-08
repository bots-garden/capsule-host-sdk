package capsule

import (
	"context"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// DefineHostFuncGetEnv defines a new host function to get the environment variable value.
//
// Parameters:
// - builder: the HostModuleBuilder to add the function to.
//
// Returns: nothing.
func DefineHostFuncGetEnv(builder wazero.HostModuleBuilder) {
		builder.NewFunctionBuilder().
		WithGoModuleFunction(getEnv, 
			[]api.ValueType{
				api.ValueTypeI32, // string position
				api.ValueTypeI32, // string length
				api.ValueTypeI32, // returned string position
				api.ValueTypeI32, // returned string length
			}, 
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostGetEnv")
}

// talk : host function called by the wasm function
// and then returning data to the wasm module
var getEnv = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

	// Position and size of the message coming from the WASM module
	position := uint32(params[0]) 
	length := uint32(params[1])
	// Read the buffer memory to retrieve the message
	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}
	variableName := string(buffer)
	variableValue := os.Getenv(variableName)

	variableValueLength := len(variableValue)

	// This is a wasm function defined in the capsule-module-sdk
	results, err := module.ExportedFunction("allocateBuffer").Call(ctx, uint64(variableValueLength))
	if err != nil {
		log.Panicln("Problem when calling allocateBuffer", err)
	}

	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	allocatedPosition := uint32(results[0])
	module.Memory().WriteUint32Le(positionReturnBuffer, allocatedPosition)
	module.Memory().WriteUint32Le(lengthReturnBuffer, uint32(variableValueLength))

	// add the message to the memory of the module
	module.Memory().Write(allocatedPosition, []byte(variableValue))

	params[0] = 0

})
