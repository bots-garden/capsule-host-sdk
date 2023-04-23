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



// Success Failure
