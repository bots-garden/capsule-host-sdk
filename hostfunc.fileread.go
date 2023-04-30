package capsule

// this hostfunction is a template for the other host functions
import (
	"context"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// DefineHostFuncReadFile defines a host function
func DefineHostFuncReadFile(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(readFile,
			[]api.ValueType{
				api.ValueTypeI32, // filePath position
				api.ValueTypeI32, // filePath length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostReadFile")
}

// readFile : host function called by the wasm function
// and then returning data to the wasm module
var readFile = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {
	
	filePathPosition := uint32(params[0])
	filePathLength := uint32(params[1])

	bufferFilePath, err := ReadBytesParameterFromMemory(module, filePathPosition, filePathLength)
	if err != nil {
		log.Panicf("Error (bufferFilePath): ReadBytesParameterFromMemory(%d, %d) out of range", filePathPosition, filePathLength)
	}

	var resultFromHost []byte
	data, errReadFile := os.ReadFile(string(bufferFilePath))


	if errReadFile != nil {
		resultFromHost = failure([]byte(errReadFile.Error()))
	} else {
		resultFromHost = success(data)
	}

	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})

/* Documentation:

*/
