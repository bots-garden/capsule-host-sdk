package capsule
// this hostfunction is a template for the other host functions
import (
	"context"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// DefineHostFuncTalk defines a host function
func DefineHostFuncTalk(builder wazero.HostModuleBuilder) {
		// hostLogString
		builder.NewFunctionBuilder().
		WithGoModuleFunction(talk, 
			[]api.ValueType{
				api.ValueTypeI32, // string position
				api.ValueTypeI32, // string length
				api.ValueTypeI32, // returned string position
				api.ValueTypeI32, // returned string length
			}, 
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostTalk")
}

// talk : host function called by the wasm function
// and then returning data to the wasm module
var talk = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

	// Position and size of the message coming from the WASM module
	position := uint32(params[0]) 
	length := uint32(params[1])
	// Read the buffer memory to retrieve the message
	buffer, ok := module.Memory().Read(position, length)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", position, length)
	}

	messageFromModule := string(buffer)
	fmt.Println("🟣 message from the WASM module:", messageFromModule)

	// Create a message from the host to reply the guest WASM module
	messageFromHost := "Hello 😀 I'm the host" 

	messageFromHostLength := len(messageFromHost)

	// This is a wasm function defined in the capsule-module-sdk
	results, err := module.ExportedFunction("allocateBuffer").Call(ctx, uint64(messageFromHostLength))
	if err != nil {
		log.Panicln("Problem when callibg allocateBuffer", err)
	}

	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	allocatedPosition := uint32(results[0])
	module.Memory().WriteUint32Le(positionReturnBuffer, allocatedPosition)
	module.Memory().WriteUint32Le(lengthReturnBuffer, uint32(messageFromHostLength))

	// add the message to the memory of the module
	module.Memory().Write(allocatedPosition, []byte(messageFromHost))

	params[0] = 0

})
