package capsule

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/bots-garden/capsule-host-sdk/helpers"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	"github.com/redis/go-redis/v9"
)

var redisDb *redis.Client

// InitRedisClient initializes a Redis client instance if it is not already initialized.
func InitRedisClient() {

	if redisDb == nil {

		redisURI := helpers.GetEnv("REDIS_URI", "")

		if redisURI == "" {
			addr := helpers.GetEnv("REDIS_ADDR", "localhost:6379")
			password := helpers.GetEnv("REDIS_PWD", "") // no password set

			defaultDb, _ := strconv.Atoi(helpers.GetEnv("REDIS_DEFAULTDB", "0"))

			redisDb = redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
				DB:       defaultDb, // use default DB
			})
		} else {
			addr, err := redis.ParseURL(redisURI)
			if err != nil {
				// TODO: handle this error
				panic(err)
			}
			//? how to handle the "non connection"?
			redisDb = redis.NewClient(addr)
		}

	}
}

// getRedisClient returns a pointer to a Redis client.
//
// This function takes no parameters.
// It returns a pointer to a Redis client.
func getRedisClient() *redis.Client {
	return redisDb
}

// DefineHostFuncRedisSet defines a Go function that sets a value in Redis.
//
// It takes in the key and value string positions and lengths as well as the
// positions and lengths of the returned value. It returns an integer value.
func DefineHostFuncRedisSet(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(redisSet,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // value string position
				api.ValueTypeI32, // value string length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostRedisSet")
}

// redisSet : host function called by the wasm function
// and then returning data to the wasm module
var redisSet = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

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

	// start the host work (using Redis client)
	InitRedisClient()                                                                      // initialize the redis client only if it does not exist
	err = getRedisClient().Set(ctx, string(bufferKey), string(bufferStringValue), 0).Err() // TODO: check if []byte is ok for the value
	if err != nil {
		resultFromHost = failure([]byte(err.Error()))
	} else {
		resultFromHost = success(bufferKey)
	}
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

// DefineHostFuncRedisGet defines a function that gets a value from Redis cache.
func DefineHostFuncRedisGet(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(redisGet,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostRedisGet")
}

// redisGet : host function called by the wasm function
// and then returning data to the wasm module
var redisGet = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

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
	result, err := getRedisClient().Get(ctx, string(bufferKey)).Result()
	if err != nil {
		resultFromHost = failure([]byte(err.Error()))
	} else {
		resultFromHost = success([]byte(result))
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

// DefineHostFuncRedisDel defines a Redis Del operation for the host module builder.
//
// This function takes in a `builder` of type `wazero.HostModuleBuilder` and creates a new
// function builder for Redis Del operation. The function builder is then configured with
// parameters and exports the function with name "hostCacheDel".
func DefineHostFuncRedisDel(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(redisDel,
			[]api.ValueType{
				api.ValueTypeI32, // key position
				api.ValueTypeI32, // key length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostRedisDel")
}

// redisDel : host function called by the wasm function
// and then returning data to the wasm module
var redisDel = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

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
	_, err = getRedisClient().Del(ctx, string(bufferKey)).Result()
	if err != nil {
		resultFromHost = failure([]byte(err.Error()))
	} else {
		resultFromHost = success(bufferKey)
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

// DefineHostFuncRedisKeys defines a function that exports a host module function
// that retrieves Redis cache keys. It takes in four parameters: filter position,
// filter length, returned position and returned length. It returns an integer.
func DefineHostFuncRedisKeys(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(redisKeys,
			[]api.ValueType{
				api.ValueTypeI32, // filter position
				api.ValueTypeI32, // filter length
				api.ValueTypeI32, // returned position
				api.ValueTypeI32, // returned length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostRedisKeys")
}

// redisKeys : host function called by the wasm function
// and then returning data to the wasm module
var redisKeys = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

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
	var keysMap = make(map[string][]string)

	// call the redis KEYS command
	keys, err := getRedisClient().Keys(ctx, string(bufferFilter)).Result()
	if err != nil {
		resultFromHost = failure([]byte(err.Error()))
	} else {
		keysMap["keys"] = keys
		jsonStr, err := json.Marshal(keysMap)
		// {"keys":["key1","key2"]}
		if err != nil {
			resultFromHost = failure(jsonStr)
		} else {
			resultFromHost = success(jsonStr)
		}
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
