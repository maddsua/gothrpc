package main

import (
	"fmt"
	"net/http"

	"github.com/maddsua/gothrpc"
)

type errorResult map[string]any

func (this errorResult) StatusCode() int {
	return 400
}

var procRouter = gothrpc.Router{
	"test": &gothrpc.Method{
		GET: gothrpc.HandleFn(func(ctx *gothrpc.Context) (any, error) {
			return map[string]any{
				"args":    "ctx.Args",
				"message": "well this didn't do much!",
			}, nil
		}),
		POST: gothrpc.HandleFn(func(ctx *gothrpc.Context) (any, error) {
			val := ctx.Value.(int)
			ctx.Value = val + 1
			return "whoops!", nil
		}),
		DELETE: gothrpc.HandleFn(func(ctx *gothrpc.Context) (any, error) {
			return errorResult{"test": "test errors"}, nil
		}),
	},
	"next": gothrpc.Router{
		"test": &gothrpc.Procedure[any, string, string]{
			Query: gothrpc.QueryHandlerFn(func(ctx *gothrpc.Context, args gothrpc.Args) (string, error) {
				return "whoa a next gen test fr", nil
			}),
			Mutation: gothrpc.MutationHandlerFn(func(ctx *gothrpc.Context, args gothrpc.Args, p any) (string, error) {
				fmt.Printf("payload: %v\n", p)
				return "ok so this would imply that we did modify something, eh?", nil
			}),
		},
	},
	"props": &gothrpc.Procedure[any, any, any]{
		Query: gothrpc.QueryHandlerFn(func(ctx *gothrpc.Context, args gothrpc.Args) (any, error) {
			return ctx.Value, nil
		}),
	},
	"panic": &gothrpc.Procedure[any, any, any]{
		Query: gothrpc.QueryHandlerFn(func(ctx *gothrpc.Context, args gothrpc.Args) (any, error) {
			panic("test panic")
		}),
	},
}

func main() {

	const serverPort = "7774"
	const apiPrefix = "/"

	procHandler := &gothrpc.RestHandler{
		Router: procRouter,
		Prefix: apiPrefix,
		OnBeforeHandle: func(ctx *gothrpc.Context) error {
			ctx.Value = int(42)
			return nil
		},
		OnAfterHandle: func(ctx *gothrpc.Context, result *gothrpc.RestResponse) error {
			result.Headers = http.Header{}
			result.Headers.Set("X-Value", fmt.Sprintf("%v", ctx.Value))
			return nil
		},
		OnError: func(err error, ctx *gothrpc.Context) {
			fmt.Printf("handler error: %s\n", err.Error())
		},
	}

	mux := http.NewServeMux()

	mux.Handle(apiPrefix, procHandler)

	fmt.Printf("Listening at: http://localhost:%s\n", serverPort)

	http.ListenAndServe(":"+serverPort, mux)
}
