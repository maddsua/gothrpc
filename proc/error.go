package proc

import "net/http"

type ProcError struct {
	Message    string         `json:"message"`
	Extensions map[string]any `json:"extensions,omitempty"`
}

func (this ProcError) Error() string {
	return this.Message
}

type ProcErrorWithCode struct {
	ProcError
	code int
}

func (this ProcErrorWithCode) StatusCode() int {
	return this.code
}

var ErrorProcedureNotFound = ProcErrorWithCode{
	ProcError: ProcError{
		Message: "procedure not found",
	},
	code: http.StatusNotFound,
}

var ErrorMethodNotAllowed = ProcErrorWithCode{
	ProcError: ProcError{
		Message: "method not allowed",
	},
	code: http.StatusMethodNotAllowed,
}
