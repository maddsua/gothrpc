package gothrpc

import "net/http"

type ProcError struct {
	Message    string         `json:"message"`
	Extensions map[string]any `json:"extensions,omitempty"`
	HttpStatus int            `json:"-"`
}

func (this ProcError) Error() string {
	return this.Message
}

func (this ProcError) StatusCode() int {

	if this.HttpStatus < http.StatusBadRequest {
		return http.StatusBadRequest
	}

	return this.HttpStatus
}

var errProcNotFound = ProcError{
	Message:    "procedure not found",
	HttpStatus: http.StatusNotFound,
}

var errMethodNotAllowed = ProcError{
	Message:    "method not allowed",
	HttpStatus: http.StatusMethodNotAllowed,
}
