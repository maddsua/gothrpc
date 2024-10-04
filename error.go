package gothrpc

import "net/http"

type Error struct {
	Message    string         `json:"message"`
	Extensions map[string]any `json:"extensions,omitempty"`
	HttpStatus int            `json:"-"`
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
