package gothrpc

import "net/http"

type Statuser interface {
	StatusCode() int
}

type Headerer interface {
	Headers() http.Header
}

type Result struct {
	Data   any        `json:"data"`
	Error  *ProcError `json:"error,omitempty"`
	status int
	header http.Header
}
