package capsule

// this hostfunction is a template for the other host functions
import (
	"context"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// DefineHostFuncWriteFile defines a host function
func DefineHostFuncWriteFile(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(writeFile,
			[]api.ValueType{
				api.ValueTypeI32, // filePath position
				api.ValueTypeI32, // filePath length
				api.ValueTypeI32, // content position
				api.ValueTypeI32, // content length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostWriteFile")
}

// writeFile : host function called by the wasm function
// and then returning data to the wasm module
var writeFile = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {
	
	filePathPosition := uint32(params[0])
	filePathLength := uint32(params[1])

	bufferFilePath, err := ReadBytesParameterFromMemory(module, filePathPosition, filePathLength)
	if err != nil {
		log.Panicf("Error (bufferFilePath): ReadBytesParameterFromMemory(%d, %d) out of range", filePathPosition, filePathLength)
	}

	contentPosition := uint32(params[2])
	contentLength := uint32(params[3])

	bufferContent, err := ReadBytesParameterFromMemory(module, contentPosition, contentLength)
	if err != nil {
		log.Panicf("Error (bufferContent): ReadBytesParameterFromMemory(%d, %d) out of range", contentPosition, contentLength)
	}

	var resultFromHost []byte
	errWriteFile := os.WriteFile(string(bufferFilePath), bufferContent, 0644)

	if errWriteFile != nil {
		resultFromHost = failure([]byte(errWriteFile.Error()))
	} else {
		resultFromHost = success(bufferFilePath)
	}

	positionReturnBuffer := uint32(params[4])
	lengthReturnBuffer := uint32(params[5])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})

/* Documentation:
! concurrent are not managed
? don't use it with the capsule-http runner
 */
