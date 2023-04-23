package capsule

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// GetRuntime returns the runtime
func GetRuntime(ctx context.Context) wazero.Runtime {
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return runtime
}
