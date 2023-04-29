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
	DefineHostFuncLog(builder)
	DefineHostFuncPrint(builder)
	DefineHostFuncTalk(builder)

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
