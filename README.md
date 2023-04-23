# Capsule Host SDK

🚧 this is a work in progress

This SDK allows to create and manage a WebAssembly host application using [WASI Capsule plugins](https://github.com/bots-garden/capsule-module-sdk).

The Capsule Host SDK use the **[Wazero](https://github.com/tetratelabs/wazero)** runtime to run the host application.

## Getting started: the host application

```golang
// Package main: host runtime
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	capsule "github.com/bots-garden/capsule-host-sdk"
)

func main() {

	ctx := context.Background()
	runtime := capsule.GetRuntime(ctx)

	builder := capsule.GetBuilder(runtime)
	// Instantiate builder with default host functions
	builder.Instantiate(ctx)

	defer runtime.Close(ctx)

	// Load the WebAssembly module
	wasmPath := "./main.wasm"
	helloWasm, _ := os.ReadFile(wasmPath)
 
	mod, _ := runtime.Instantiate(ctx, helloWasm)

	// Get the reference to the WebAssembly handleFunction
	handleFunction := capsule.GetHandle(mod)

	pos, size, _ := capsule.CopyDataToMemory(ctx, mod, []byte("Bob Morane"))

	// Call handleFunction with the position and the size of "Bob Morane"
	res, _ := handleFunction.Call(ctx, pos, size)

	resPos, resSize := capsule.UnPackPosSize(res[0])

	bytesRes, err := capsule.ReadBytesFromMemory(mod, resPos, resSize)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(bytesRes))
	}
}
```

## Getting started: the capsule plugin

```golang
package main

import (
	capsule "github.com/bots-garden/capsule-module-sdk"
)

func main() {
	capsule.SetHandle(Handle)
}

// Handle function
func Handle(param []byte) ([]byte, error) {

	capsule.Log("🟣 from the plugin: " + string(param))
	capsule.Print("💜 from the plugin: " + string(param))

	return []byte("Hello " + string(param)), nil
}
```
