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

func (this *Result) StatusCode() int {
	return this.status
}

func (this *Result) Headers() http.Header {
	return this.header
}
