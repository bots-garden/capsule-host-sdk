package capsule

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// GetBuilder returns a new instance of the HostModuleBuilder
// configured with the default host functions
func GetBuilder(runtime wazero.Runtime) wazero.HostModuleBuilder {
	builder := runtime.NewHostModuleBuilder("env")

	// Define default host functions
	DefineHostFuncLog(builder) // see hostfunc.log.go
	DefineHostFuncPrint(builder) // see hostfunc.print.go
	DefineHostFuncTalk(builder) // see hostfunc.talk.go //! this one is a kind of template
	DefineHostFuncGetEnv(builder) // see hostfunc.getenv.go
	DefineHostFuncWriteFile(builder) // see hostfunc.filewrite.go
	DefineHostFuncReadFile(builder) // see hostfunc.readfile.go
	DefineHostFuncHTTP(builder) // see hostfunc.http.go

	DefineHostFuncCacheGet(builder) // see hostfunc.memorycache.go
	DefineHostFuncCacheSet(builder) // see hostfunc.memorycache.go
	DefineHostFuncCacheDel(builder) // see hostfunc.memorycache.go
	DefineHostFuncCacheKeys(builder) // see hostfunc.memorycache.go

	DefineHostFuncRedisDel(builder)
	DefineHostFuncRedisGet(builder)
	DefineHostFuncRedisKeys(builder)
	DefineHostFuncRedisSet(builder)

	return builder
}

// GetRuntime returns the WebAssembly runtime.
// It takes a context and returns a wazero.Runtime object.
func GetRuntime(ctx context.Context) wazero.Runtime {
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return runtime
}


// CallHandleFunction calls the given handleFunction with the argument argFunction
// and returns the result. The function uses CopyDataToMemory to copy the argument
// to memory, and UnPackPosSize to unpack the result. Returns a byte slice and an
// error.
//
// ctx: The context.Context
// mod: The api.Module
// handleFunction: The api.Function to be called
// argFunction: The argument to the function
//
// Returns ([]byte, error).
func CallHandleFunction(ctx context.Context, mod api.Module, handleFunction api.Function, argFunction []byte) ([]byte, error) {

	// send argument to the function
	pos, size, err := CopyDataToMemory(ctx, mod, argFunction)
	if err != nil {
		return nil, err
	}

	// Now, we can call "callHandle" with the position and the size of "Bob Morane"
	// the result type is []uint64
	result, err := handleFunction.Call(ctx, pos, size)
	if err != nil {
		return nil, err
	}
	// read the result of the function
	resultMemoryPosition, resultSize := UnPackPosSize(result[0])

	bufferResult, err := ReadDataFromMemory(mod, resultMemoryPosition, resultSize)
	if err != nil {
		return nil, err
	}

	return Result(bufferResult)
}
