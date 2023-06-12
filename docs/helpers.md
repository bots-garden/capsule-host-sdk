# 🧰 Helpers

> 🚧 this is a work in progress

## Call OnStart exported method (from the wasm module)
> introduced in v0.0.4

```golang
// Package main
package main

import (
	"strconv"
	"github.com/bots-garden/capsule-module-sdk"
)

func main() {
	capsule.SetHandleHTTP(func (param capsule.HTTPRequest) (capsule.HTTPResponse, error) {
		return capsule.HTTPResponse{
			TextBody: "👋 Hey",
			Headers: `{"Content-Type": "text/plain; charset=utf-8"}`,
			StatusCode: 200,
		}, nil
		
	})
}

// OnStart function
//export OnStart
func OnStart() {
	capsule.Print("🚗 OnStart")
}
```
> 👋 don't forget to export the `OnStart` function

## Call OnStop exported method (from the wasm module)
> introduced in v0.0.4

```golang
// Package main
package main

import (
	"strconv"
	"github.com/bots-garden/capsule-module-sdk"
)

func main() {
	capsule.SetHandleHTTP(func (param capsule.HTTPRequest) (capsule.HTTPResponse, error) {
		return capsule.HTTPResponse{
			TextBody: "👋 Hey",
			Headers: `{"Content-Type": "text/plain; charset=utf-8"}`,
			StatusCode: 200,
		}, nil
		
	})
}


// OnStop function
//export OnStop
func OnStop() {
	capsule.Print("🚙 OnStop")
}
```
> 👋 don't forget to export the `OnStop` function
