package gothrpc

import "net/http"

type ProcedureError interface {
	ProcError() *Error
}

type Error struct {
	Message    string          `json:"message"`
	Cause      string          `json:"cause,omitempty"`
	Extensions ErrorExtensions `json:"extensions,omitempty"`
	HttpStatus int             `json:"-"`
}

type ErrorExtensions map[string]any

func (this Error) ProcError() *Error {
	return &this
}

func (this Error) Error() string {
	return this.Message
}

func (this Error) StatusCode() int {

	if this.HttpStatus < http.StatusBadRequest {
		return http.StatusBadRequest
	}

	return this.HttpStatus
}

var errProcNotFound = Error{
	Message:    "procedure not found",
	HttpStatus: http.StatusNotFound,
}

var errMethodNotAllowed = Error{
	Message:    "method not allowed",
	HttpStatus: http.StatusMethodNotAllowed,
}
