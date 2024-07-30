package proc

import (
	"net/http"
)

type Context struct {
	Headers    http.Header
	Method     string
	RemoteAddr string
	Args       Args
	Props      any
	ProcPath   ProcedureStepper
}

type Args map[string]any

type Handler interface {
	Handle(ctx Context) (any, error)
}

func HandleFunc(handler func(ctx Context) (any, error)) Handler {
	return &handlerFuncWrapper{
		handler: handler,
	}
}

type handlerFuncWrapper struct {
	handler func(ctx Context) (any, error)
}

func (this *handlerFuncWrapper) Handle(ctx Context) (any, error) {

	if ctx.ProcPath.HasNext() {
		return nil, ErrorProcedureNotFound
	}

	return this.handler(ctx)
}
