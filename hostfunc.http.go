package capsule

import (
	"context"
	"encoding/json"
	"errors"
	//"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//! THIS IS A WORK IN PROGRESS → it does not work very well with big requests

// request embeds the data of the http request
type request struct {
	JSONBody map[string]interface{} `json:"JSONBody"`
	TextBody string                 `json:"TextBody"`
	//Body     string                 `json:"Body"`
	URI      string                 `json:"URI"`
	Method   string                 `json:"Method"`
	Headers  map[string]string      `json:"Headers"`
}

// DefineHostFuncHTTP defines a host function
func DefineHostFuncHTTP(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(http,
			[]api.ValueType{
				api.ValueTypeI32, // request position
				api.ValueTypeI32, // request length
				api.ValueTypeI32, // returned value position
				api.ValueTypeI32, // returned value length
			},
			[]api.ValueType{api.ValueTypeI32}).
		Export("hostHTTP")
}

// http : host function called by the wasm function
// and then returning data to the wasm module
var http = api.GoModuleFunc(func(ctx context.Context, module api.Module, params []uint64) {

	requestPosition := uint32(params[0])
	requestLength := uint32(params[1])

	bufferRequest, err := ReadBytesParameterFromMemory(module, requestPosition, requestLength)
	if err != nil {
		log.Panicf("❌ Error (bufferRequest): ReadBytesParameterFromMemory(%d, %d) out of range", requestPosition, requestLength)
	}
	// unmarshal the request
	var req request
	errMarshal := json.Unmarshal(bufferRequest, &req)
	if errMarshal != nil {
		log.Println("❌ Error when unmarshal the request", errMarshal)
	}

	var resultFromHost []byte

	httpClient := resty.New()

	for key, value := range req.Headers {
		httpClient.SetHeader(key, value)
	}

	switch what := req.Method; what {
	case "GET":

		resp, err := httpClient.R().EnableTrace().Get(req.URI)

		if err != nil {
			resultFromHost = failure([]byte(err.Error()))
		} else {

			jsonHTTPResponse, err := buildResponseJSONString(resp)
			if err != nil {
				resultFromHost = failure([]byte(err.Error()))
			}
			resultFromHost = success([]byte(jsonHTTPResponse))
		}

	case "POST":
		var body string
		/*
		if req.Body != "" { // TODO: remove Body
			body = req.Body
		} else if req.JSONBody != nil {
			buff, _ := json.Marshal(req.JSONBody)
			// TODO: handle error
			body = string(buff)
		} else if req.TextBody != "" {
			body = req.TextBody
		}
		*/
		if req.JSONBody != nil {
			buff, _ := json.Marshal(req.JSONBody)
			// TODO: handle error
			body = string(buff)
		} else if req.TextBody != "" {
			body = req.TextBody
		}


		resp, err := httpClient.R().EnableTrace().SetBody(body).Post(req.URI)

		if err != nil {
			resultFromHost = failure([]byte(err.Error()))
		} else {
			jsonHTTPResponse, err := buildResponseJSONString(resp)
			if err != nil {
				resultFromHost = failure([]byte(err.Error()))
			}
			resultFromHost = success([]byte(jsonHTTPResponse))
		}

	case "PUT":
		// TODO: test it
		var body string
		if req.JSONBody != nil {
			buff, _ := json.Marshal(req.JSONBody)
			// TODO: handle error
			body = string(buff)
		} else if req.TextBody != "" {
			body = req.TextBody
		}

		resp, err := httpClient.R().EnableTrace().SetBody(body).Put(req.URI)

		if err != nil {
			resultFromHost = failure([]byte(err.Error()))
		} else {
			jsonHTTPResponse, err := buildResponseJSONString(resp)
			if err != nil {
				resultFromHost = failure([]byte(err.Error()))
			}
			resultFromHost = success([]byte(jsonHTTPResponse))
		}


	case "DELETE":
		// TODO: test it
		resp, err := httpClient.R().EnableTrace().Delete(req.URI)

		if err != nil {
			resultFromHost = failure([]byte(err.Error()))
		} else {

			jsonHTTPResponse, err := buildResponseJSONString(resp)
			if err != nil {
				resultFromHost = failure([]byte(err.Error()))
			}
			resultFromHost = success([]byte(jsonHTTPResponse))
		}

	default:
		resultFromHost = failure([]byte(errors.New("❌ Error: " + req.Method + " is not yet implemented").Error()))
	}

	positionReturnBuffer := uint32(params[2])
	lengthReturnBuffer := uint32(params[3])

	_, errReturn := ReturnBytesToMemory(ctx, module, positionReturnBuffer, lengthReturnBuffer, resultFromHost)
	if errReturn != nil {
		log.Panicf("❌ Error: ReturnBytesToMemory(%d, %d) out of range", positionReturnBuffer, lengthReturnBuffer)
	}

	params[0] = 0

})

func buildResponseJSONString(resp *resty.Response) (string, error) {
	// build headers JSON string
	// ! ATTENTION resp.Header() return a map[string]string[] (instead of map[string]string)
	// TODO: on the guest side, add method to the structure to read the headers
	// TODO: rebuild the headers and copy it to a map[string]string
	/*
		   for key, value := range resp.Header() {
		       fmt.Println(key, value[0])
			}
	*/
	// TODO: or try with another library

	jsonHeaders, err := json.Marshal(resp.Header())
	responseBody := resp.String() //? marshall or not?
	statusCode := resp.StatusCode()

	isJSON := false
	contentType, ok := resp.Header()["Content-Type"]
	if ok {
		isJSON = resty.IsJSONType(contentType[0])
	}

	var jsonHTTPResponse string
	
	if  isJSON {
		jsonHTTPResponse = `{"JSONBody":` + responseBody + `,"Headers":` + string(jsonHeaders) + `,"StatusCode":` + strconv.Itoa(statusCode) + `}`

	} else {
		// add double quotes for body
		jsonHTTPResponse = `{"TextBody":"` + responseBody + `","Headers":` + string(jsonHeaders) + `,"StatusCode":` + strconv.Itoa(statusCode) + `}`
	}


	return jsonHTTPResponse, err
}

/* Documentation:

 */
