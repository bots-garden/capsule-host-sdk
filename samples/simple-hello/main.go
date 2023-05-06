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

	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	runtime := capsule.GetRuntime(ctx)

	// START: host functions
	builder := capsule.GetBuilder(runtime)

	// 🤚 Add your host functions here
	DefineHostFuncPrintHello(builder)

	// Instantiate builder and default host functions
	_, err := builder.Instantiate(ctx)
	if err != nil {
		log.Panicln("Error with env module and host function(s):", err)
	}
	// END: host functions

	// This closes everything this Runtime created.
	defer runtime.Close(ctx)

	// Load the WebAssembly module
	wasmPath := "../../../capsule-module-sdk/samples/simple-hello/simple-hello.wasm"
	helloWasm, err := os.ReadFile(wasmPath)
	if err != nil {
		log.Panicln("📝", err)
	}

	// 👀 see https://github.com/tetratelabs/wazero/blob/main/examples/concurrent-instantiation/main.go
	mod, err := runtime.Instantiate(ctx, helloWasm)
	if err != nil {
		log.Panicln("🥚", err)
	}

	// Get the reference to the WebAssembly function: "callHandle"
	// callHandle is exported by the Capsule plugin
	handleFunction := capsule.GetHandle(mod)

	pos, size, err := capsule.CopyDataToMemory(ctx, mod, []byte("Bob Morane"))
	if err != nil {
		log.Panicln(err)
	}

	// Now, we can call "callHandle" with the position and the size of "Bob Morane"
	// the result type is []uint64
	r, err := handleFunction.Call(ctx, pos, size)
	if err != nil {
		log.Panicln(err)
	}	

	rpos, rsize := capsule.UnPackPosSize(r[0])

	res, err := capsule.ReadBytesFromMemory(mod, rpos, rsize)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}
