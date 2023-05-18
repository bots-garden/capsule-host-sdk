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
func DownloadWasmFile(wasmFileURL, wasmFilePath, authenticationHeader, authenticationHeaderValue string) ([]byte, error) {

	// authenticationHeader:
	// Example: "PRIVATE-TOKEN: ${GITLAB_WASM_TOKEN}"
	// SetHeader("Accept", "application/json").

	client := resty.New()

	if authenticationHeader != "" {
		client.SetHeader(authenticationHeader, authenticationHeaderValue)
	} 

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

// TODO:
// GitLab registry (with and without token) / it's http
// GitHub registry (with and without token) looks like I need to use an OCI lib
// From S3
// https://wapm.io/ (if possible)
// Other OCI registries
// From GitLab and GitHub release? (it's http?)

