package gothrpc

import (
	"net/http"
)

type Handler interface {
	Handle(ctx Context) (any, error)
}

type QueryHandler[R any] interface {
	Handle(ctx Context, args Args) (R, error)
}

type MutationHandler[P, R any] interface {
	Handle(ctx Context, args Args, props P) (R, error)
}

type Context struct {
	//	Original http request pointer
	Req *http.Request
	//	A custom value passed between procedure and executor hooks
	Value    any
	procPath procStepper
}

type Args map[string]string

func HandleFn(handler func(ctx Context) (any, error)) Handler {
	return &handlerFuncWrapper{
		handler: handler,
	}
}

type handlerFuncWrapper struct {
	handler func(ctx Context) (any, error)
}

func (this *handlerFuncWrapper) Handle(ctx Context) (any, error) {

	if ctx.procPath.HasNext() {
		return nil, errProcNotFound
	}

	return this.handler(ctx)
}

func QueryHandlerFn[R any](handler func(ctx Context, args Args) (R, error)) QueryHandler[R] {
	return &queryHandlerFnWrapper[R]{
		handler: handler,
	}
}

type queryHandlerFnWrapper[R any] struct {
	handler func(ctx Context, args Args) (R, error)
}

func (this *queryHandlerFnWrapper[R]) Handle(ctx Context, args Args) (R, error) {
	return this.handler(ctx, args)
}

func MutationHandlerFn[P, R any](handler func(ctx Context, args Args, payload P) (R, error)) MutationHandler[P, R] {
	return &mutHandlerFnWrapper[P, R]{
		handler: handler,
	}
}

type mutHandlerFnWrapper[P, R any] struct {
	handler func(ctx Context, args Args, payload P) (R, error)
}

func (this *mutHandlerFnWrapper[P, R]) Handle(ctx Context, args Args, payload P) (R, error) {
	return this.handler(ctx, args, payload)
}
