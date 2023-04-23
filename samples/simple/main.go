package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bots-garden/capsule-host-sdk"
)

func main() {

	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	runtime := capsule.GetRuntime(ctx)
	
	// This closes everything this Runtime created.
	defer runtime.Close(ctx)

	// Load the WebAssembly module
	wasmPath := "../../../capsule-module-sdk/samples/simple/main.wasm"
	helloWasm, err := os.ReadFile(wasmPath)
	if err != nil {
		log.Panicln("📝", err)
	}

	mod, err := runtime.Instantiate(ctx, helloWasm)
	if err != nil {
		log.Panicln("🥚", err)
	}

	// Get the reference to the WebAssembly function: "callHandle"
	// callHandle is exported by the Capsule plugin
	handleFunction := mod.ExportedFunction("callHandle")


	pos, size, err := capsule.CopyDataToMemory(ctx, mod, []byte("Bob Morane"))
	if err != nil {
		log.Panicln(err)
	}

	// Now, we can call "callHandle" with the position and the size of "Bob Morane"
	// the result type is []uint64
	result, err := handleFunction.Call(ctx, pos, size)
	if err != nil {
		log.Panicln(err)
	}	

	rpos, rsize := capsule.UnPackPosSize(result[0])

	bRes, err := capsule.ReadDataFromMemory(mod, rpos, rsize)
	if err != nil {
		log.Panicln(err)
	}
	
	res, err := capsule.Result(bRes)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(res))
	}
}
