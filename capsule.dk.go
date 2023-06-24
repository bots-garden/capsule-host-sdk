// Package capsule SDK for host applications
package capsule

import (
	"context"
	"errors"
	"log"

	//"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const isFailure = rune('F')
const isSuccess = rune('S')

// success appends the isSuccess byte to the beginning of the input buffer and returns the result.
//
// buffer: byte slice to append isSuccess byte to.
// []byte: byte slice with the appended isSuccess byte.
func success(buffer []byte) []byte {
	return append([]byte(string(isSuccess)), buffer...)
}

// failure appends a string "isFailure" to the given byte slice buffer and returns the new slice.
//
// buffer: the byte slice to which "isFailure" is appended.
// Returns the new byte slice with the string "isFailure" appended to it.
func failure(buffer []byte) []byte {
	return append([]byte(string(isFailure)), buffer...)
}

// Result returns the data without the first byte if the first byte is isSuccess.
// Otherwise, it returns nil and an error with the data starting from the second byte.
//
// data: A byte slice containing the data to check.
// []byte: The data without the first byte if the first byte is isSuccess.
// error: If the first byte is not isSuccess, it returns an error with the data starting from the second byte.
func Result(data []byte,) ([]byte, error) {
	if data[0] == byte(isSuccess) {
		return data[1:], nil
	}
	return nil, errors.New(string(data[1:]))
}

// GetHandle returns an exported function named "callHandle" from the given module.
//
// mod: The module to retrieve the function from.
//
// Returns: An exported function with the name "callHandle".
func GetHandle(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandle")
}

// GetHandleJSON returns the exported "callHandleJSON" function from the given module.
//
// mod: the module to retrieve the function from.
//
// returns: the exported "callHandleJSON" function.
func GetHandleJSON(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandleJSON")
}

// GetHandleHTTP returns the exported 'callHandleHTTP' function from a given module.
//
// mod: The module containing the exported function.
//
// returns:
//     - api.Function: the exported 'callHandleHTTP' function.
func GetHandleHTTP(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandleHTTP")
}

// CallOnStart calls the OnStart function (if it exists) from the given module.
func CallOnStart(ctx context.Context, mod api.Module , wasmFile []byte) {

	onStart := mod.ExportedFunction("OnStart")
	if onStart != nil {
		_, err := onStart.Call(ctx)
		if err != nil {
			log.Println("❌ Error calling OnStart", err)
			panic(err)
		}
	}
}

// CallOnStop calls the OnStop function (if it exists) from the given module.
func CallOnStop(ctx context.Context, mod api.Module, wasmFile []byte) {

	onStop := mod.ExportedFunction("OnStop")
	if onStop != nil {
		_, err := onStop.Call(ctx)
		if err != nil {
			log.Println("❌ Error calling OnStop", err)
			panic(err)
		}
	}
}
