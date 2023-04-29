package helpers

import (
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

// LoadWasmFile loads a wasm file locally
func LoadWasmFile(wasmFilePath string) ([]byte, error) {
	wasmFileToLoad, errLoadWasmFile := os.ReadFile(wasmFilePath)
	return wasmFileToLoad, errLoadWasmFile
}

// DownloadWasmFile downloads a wasm file from a remote location
func DownloadWasmFile(wasmFileURL, wasmFilePath string) ([]byte, error) {

	client := resty.New()
	resp, err := client.R().
		SetOutput(wasmFilePath).
		Get(wasmFileURL)

	if resp.IsError() {
		return nil, errors.New("error while downloading the wasm file")
	}

	if err != nil {
		return nil, err
	}
	
	return LoadWasmFile(wasmFilePath)

}
