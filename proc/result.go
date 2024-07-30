package proc

import "net/http"

type Result struct {
	Data   any        `json:"data"`
	Error  *ProcError `json:"error,omitempty"`
	status int
	header http.Header
}

func (this *Result) Status() int {
	return this.status
}

func (this *Result) Header() http.Header {
	return this.header
}

type Statuser interface {
	StatusCode() int
}

//	yeah make jokes about it's name, I don't care
type Headerer interface {
	Headers() http.Header
}
