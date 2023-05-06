package capsule

/* Documentation:
## 4 host functions: cacheSet, cacheGet, cacheDel, cacheKeys

*/

/* TODO
- implement LoadOrStore
- implement Filter
- implement ForEach
*/

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var memCache sync.Map

// DefineHostFuncCacheSet defines a host function
func DefineHostFuncCacheSet(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(cacheSet,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // value string position
				api.ValueTypeI32, // value string length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostCacheSet")
}

// cacheSet : host function called by the wasm function
// and then returning data to the wasm module
var cacheSet = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {
	
	// read the value of the arguments of the function
	keyPosition := uint32(params[0])
	keyLength := uint32(params[1])

	bufferKey, err := ReadBytesParameterFromMemory(module, keyPosition, keyLength)
	if err != nil {
		log.Panicf("Error (bufferKey): ReadBytesParameterFromMemory(%d, %d) out of range", keyPosition, keyLength)
	}

	stringValuePosition := uint32(params[2])
	stringValueLength := uint32(params[3])

	bufferStringValue, err := ReadBytesParameterFromMemory(module, stringValuePosition, stringValueLength)
	if err != nil {
		log.Panicf("Error (bufferStringValue): ReadBytesParameterFromMemory(%d, %d) out of range", stringValuePosition, stringValueLength)
	}

	// Execute the host function with the arguments and return a value
	var resultFromHost []byte
	
	// start the host work
	memCache.Store(string(bufferKey), bufferStringValue)
	resultFromHost = success(bufferKey)
	//! we cannot know if there is an error or not
	// end of the host work

	// return the result value (using the return buffer)
	positionReturnBuffer := uint32(params[4])
	lengthReturnBuffer := uint32(params[5])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})


// DefineHostFuncCacheGet defines a host function
func DefineHostFuncCacheGet(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(cacheGet,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostCacheGet")
}

// cacheGet : host function called by the wasm function
// and then returning data to the wasm module
var cacheGet = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {
	
	// read the value of the arguments of the function
	keyPosition := uint32(params[0])
	keyLength := uint32(params[1])

	bufferKey, err := ReadBytesParameterFromMemory(module, keyPosition, keyLength)
	if err != nil {
		log.Panicf("Error (bufferKey): ReadBytesParameterFromMemory(%d, %d) out of range", keyPosition, keyLength)
	}

	// Execute the host function with the arguments and return a value
	var resultFromHost []byte
	
	// start the host work
	result, ok := memCache.Load(string(bufferKey))

	if ok {
		resultFromHost = success(result.([]byte))
	} else {
		resultFromHost = failure([]byte("key not found"))
	}
	// end of the host work

	// return the result value (using the return buffer)
	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})


// DefineHostFuncCacheDel defines a host function
func DefineHostFuncCacheDel(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(cacheDel,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostCacheDel")
}

// cacheDel : host function called by the wasm function
// and then returning data to the wasm module
var cacheDel = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {
	
	// read the value of the arguments of the function
	keyPosition := uint32(params[0])
	keyLength := uint32(params[1])

	bufferKey, err := ReadBytesParameterFromMemory(module, keyPosition, keyLength)
	if err != nil {
		log.Panicf("Error (bufferKey): ReadBytesParameterFromMemory(%d, %d) out of range", keyPosition, keyLength)
	}

	// Execute the host function with the arguments and return a value
	var resultFromHost []byte
	
	// start the host work
	memCache.Delete(string(bufferKey))
	resultFromHost = success(bufferKey)
	//! we cannot know if there is an error or not
	// end of the host work

	// return the result value (using the return buffer)
	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})




// DefineHostFuncCacheKeys defines a host function
func DefineHostFuncCacheKeys(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(cacheKeys,
			[]api.ValueType{
				api.ValueTypeI32, // filter position
				api.ValueTypeI32, // filter length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostCacheKeys")
}

// cacheKeys : host function called by the wasm function
// and then returning data to the wasm module
var cacheKeys = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

	// read the value of the arguments of the function
	filterPosition := uint32(params[0])
	filterLength := uint32(params[1])

	bufferFilter, err := ReadBytesParameterFromMemory(module, filterPosition, filterLength)
	if err != nil {
		log.Panicf("Error (bufferFilter): ReadBytesParameterFromMemory(%d, %d) out of range", filterPosition, filterLength)
	}

	// Execute the host function with the arguments and return a value
	var resultFromHost []byte
	
	// start the host work

	var keys []string
	var keysMap = make(map[string][]string)

	if string(bufferFilter) == "*" {
		memCache.Range(func(key , value interface{}) bool {
			keys = append(keys, key.(string))
			return true
		})
	} else {
		//TODO: implement
		// starts with "something"
		// contains ""
		// ends with
	}
	keysMap["keys"] = keys 
	jsonStr, err := json.Marshal(keysMap)
	// {"keys":["key1","key2"]}
	if err != nil {
		resultFromHost = failure(jsonStr)
	} else {
		resultFromHost = success(jsonStr)
	}
	// end of the host work

	// return the result value (using the return buffer)
	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})
