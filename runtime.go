package capsule

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// GetBuilder returns the builder
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

	return builder
}

// GetRuntime returns the runtime
func GetRuntime(ctx context.Context) wazero.Runtime {
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return runtime
}
