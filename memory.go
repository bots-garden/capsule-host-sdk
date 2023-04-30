package capsule

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// CopyDataToMemory returns the position and size of the data in memory
func CopyDataToMemory(ctx context.Context, mod api.Module, data []byte) (uint64, uint64, error) {
	// These function are exported by TinyGo
	malloc := mod.ExportedFunction("malloc")
	free := mod.ExportedFunction("free")

	dataSize := uint64(len(data))

	// Allocate Memory for "Bob Morane"
	results, err := malloc.Call(ctx, dataSize)
	if err != nil {
		return 0, 0, err
	}
	dataPosition := results[0]

	// This pointer is managed by TinyGo,
	// but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, dataPosition)

	// Copy data to memory
	if !mod.Memory().Write(uint32(dataPosition), data) {
		return 0, 0, errors.New("out of range of memory size")
	} else {
		return dataPosition, dataSize, nil
	}
}

// UnPackPosSize extract the position and size from a unique value
func UnPackPosSize(pair uint64) (uint32, uint32) {
	// Extract the position and size of the returned value
	pos := uint32(pair >> 32)
	size := uint32(pair)
	return pos, size
}


// ReadDataFromMemory returns the data in memory
func ReadDataFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	// Read the value from the memory
	bytes, ok := mod.Memory().Read(pos, size)
	if !ok {
		return nil, errors.New("out of range of memory size")
	} 
	return bytes, nil
}

// ReadBytesFromMemory returns the data in memory
func ReadBytesFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	data, err := ReadDataFromMemory(mod, pos, size)
	if err != nil {
		return nil, err
	}
	result, err := Result(data)
	return result, err
}

//! When using host function

// ReadBytesParameterFromMemory → read the parameter(s) sent by the WASM guest
func ReadBytesParameterFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	buff, ok := mod.Memory().Read(pos, size)
	if !ok {
		return nil, errors.New("out of range of memory size")
	}
	return buff, nil
}

// ReturnBytesToMemory → return data to the WASM guest
func ReturnBytesToMemory(ctx context.Context, mod api.Module, positionReturnBuffer uint32, lengthReturnBuffer uint32, dataFromHost []byte) (bool, error) {
	dataFromHostLength := len(dataFromHost)
	// This is a wasm function defined in the capsule-module-sdk
	results, err := mod.ExportedFunction("allocateBuffer").Call(ctx, uint64(dataFromHostLength))
	if err != nil {
		return false, err
	}
	allocatedPosition := uint32(results[0])
	mod.Memory().WriteUint32Le(positionReturnBuffer, allocatedPosition)
	mod.Memory().WriteUint32Le(lengthReturnBuffer, uint32(dataFromHostLength))

	// add the message to the memory of the module
	return mod.Memory().Write(allocatedPosition, dataFromHost), nil

}

/* Documentation:

*/

