package capsule

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// CopyDataToMemory copies data to memory.
//
// ctx: The context for this function.
// mod: The module to copy the data to.
// data: The data to be copied to memory.
//
// uint64, uint64, error: The position of the copied data, the size of the data,
// and an error if one occurs.
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
	}
	return dataPosition, dataSize, nil
	
}

// UnPackPosSize extracts the position and size of the returned value from a given pair.
//
// pair: 64-bit unsigned integer.
// Returns a pair of 32-bit unsigned integers.
func UnPackPosSize(pair uint64) (uint32, uint32) {
	// Extract the position and size of the returned value
	pos := uint32(pair >> 32)
	size := uint32(pair)
	return pos, size
}

// ReadDataFromMemory reads data from a given position in the memory of a module.
//
// Parameters:
// - mod: the module to read data from.
// - pos: the position in the memory to read from.
// - size: the size of the data to read.
//
// Returns:
// - a byte slice containing the read data.
// - an error if the position or size are out of range of the memory size.
func ReadDataFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	// Read the value from the memory
	bytes, ok := mod.Memory().Read(pos, size)
	if !ok {
		return nil, errors.New("out of range of memory size")
	} 
	return bytes, nil
}

// ReadBytesFromMemory reads a sequence of bytes from the given module's memory starting from pos and
// with a length of size. It returns the bytes read and any error encountered.
func ReadBytesFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	data, err := ReadDataFromMemory(mod, pos, size)
	if err != nil {
		return nil, err
	}
	result, err := Result(data)
	return result, err
}
// ReadBytesFromMemory returns the data in memory



//! When using host function

// ReadBytesParameterFromMemory reads a slice of bytes from the given position
// in memory of the provided module. Returns the slice of bytes and an error if
// the read operation failed due to the specified position being out of range.
// 
// mod: The module from which to read memory.
// pos: The starting position to read from.
// size: The number of bytes to read.
// 
// Returns: A slice of bytes read from memory and an error if the read operation
// failed.
func ReadBytesParameterFromMemory(mod api.Module, pos uint32, size uint32) ([]byte, error) {
	buff, ok := mod.Memory().Read(pos, size)
	if !ok {
		return nil, errors.New("out of range of memory size")
	}
	return buff, nil
}

// ReturnBytesToMemory writes data from the host to a buffer in the module's memory
// and updates the buffer information in the module. It returns a boolean value
// indicating whether the write was successful and an error if any.
//
// ctx: context required for the operation.
// mod: the module where the buffer is.
// positionReturnBuffer: the position in memory where the buffer's position will be written.
// lengthReturnBuffer: the position in memory where the buffer's length will be written.
// dataFromHost: the data to be written to the buffer.
//
// Returns:
// - a boolean indicating whether the write was successful.
// - an error if any.
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
// ReturnBytesToMemory → return data to the WASM guest


