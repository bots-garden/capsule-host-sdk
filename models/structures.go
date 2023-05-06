// Package models ...
package models

// Request embeds the data of the http request
type Request struct {
	Body     string
	JSONBody string //! to use in the future
	TextBody string //! to use in the future
	URI      string
	Method   string
	Headers  string
}

// Response embeds the data of the http response
type Response struct {
	JSONBody   map[string]interface{} `json:"JSONBody"`
	TextBody   string                 `json:"TextBody"`
	Headers    map[string]string      `json:"Headers"`
	StatusCode int                    `json:"StatusCode"`
}

// 	Body    map[string]interface{} `json:"Body"`
