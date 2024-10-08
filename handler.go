package gothrpc

import (
	"net/http"
)

type Handler interface {
	Handle(ctx *Context) (any, error)
}

type QueryHandler[R any] interface {
	Handle(ctx *Context, input QueryInput) (R, error)
}

type QueryInput map[string]string

type MutationHandler[P, R any] interface {
	Handle(ctx *Context, input P) (R, error)
}

type Context struct {
	//	Original http request pointer
	Req *http.Request
	//	Origina http request writer
	Writer http.ResponseWriter
	//	A custom value passed between procedure and esxecutor hooks
	Value any
	//	Rest procedure path steps
	path procPath
}

func (this *Context) ProcName() string {
	return this.path.at()
}

func HandleFn(handler func(ctx *Context) (any, error)) Handler {
	return &handlerFuncWrapper{
		handler: handler,
	}
}

type handlerFuncWrapper struct {
	handler func(ctx *Context) (any, error)
}

func (this *handlerFuncWrapper) Handle(ctx *Context) (any, error) {

	if ctx.path.hasNext() {
		return nil, errProcNotFound
	}

	return this.handler(ctx)
}

func QueryHandlerFn[R any](handler func(ctx *Context, input QueryInput) (R, error)) QueryHandler[R] {
	return &queryHandlerFnWrapper[R]{
		handler: handler,
	}
}

type queryHandlerFnWrapper[R any] struct {
	handler func(ctx *Context, input QueryInput) (R, error)
}

func (this *queryHandlerFnWrapper[R]) Handle(ctx *Context, input QueryInput) (R, error) {
	return this.handler(ctx, input)
}

func MutationHandlerFn[P, R any](handler func(ctx *Context, input P) (R, error)) MutationHandler[P, R] {
	return &mutHandlerFnWrapper[P, R]{
		handler: handler,
	}
}

type mutHandlerFnWrapper[P, R any] struct {
	handler func(ctx *Context, input P) (R, error)
}

func (this *mutHandlerFnWrapper[P, R]) Handle(ctx *Context, input P) (R, error) {
	return this.handler(ctx, input)
}
