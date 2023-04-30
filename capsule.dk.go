// Package capsule SDK for host applications
package capsule

import (
	"errors"

	"github.com/tetratelabs/wazero/api"
)

const isFailure = rune('F')
const isSuccess = rune('S')

/*
func main() {
	panic("not implemented")
}

func success(buffer []byte) uint64 {
	return copyBufferToMemory(append([]byte(string(isSuccess)), buffer...))
}

func failure(buffer []byte) uint64 {
	return copyBufferToMemory(append([]byte(string(isFailure)), buffer...))
}
*/

func success(buffer []byte) []byte {
	return append([]byte(string(isSuccess)), buffer...)
}

func failure(buffer []byte) []byte {
	return append([]byte(string(isFailure)), buffer...)
}




// Result function
func Result(data []byte,) ([]byte, error) {
	if data[0] == byte(isSuccess) {
		return data[1:], nil
	}
	return nil, errors.New(string(data[1:]))
}

// GetHandle returns the handle function
func GetHandle(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandle")
}

// GetHandleJSON returns the handle function
func GetHandleJSON(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandleJSON")
}

// GetHandleHTTP returns the handle function
func GetHandleHTTP(mod api.Module) api.Function {
	return mod.ExportedFunction("callHandleHTTP")
}

// TODO: handle the other handles
