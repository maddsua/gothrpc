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

// todo: add procedure handlers with type definitions
var procRouter = gothrpc.Router{
	"test": &gothrpc.Method{
		GET: gothrpc.HandleFn(func(ctx gothrpc.Context) (any, error) {
			return map[string]any{
				"args":    "ctx.Args",
				"message": "well this didn't do much!",
			}, nil
		}),
		POST: gothrpc.HandleFn(func(ctx gothrpc.Context) (any, error) {
			return "whoops!", nil
		}),
		DELETE: gothrpc.HandleFn(func(ctx gothrpc.Context) (any, error) {
			return errorResult{"test": "test errors"}, nil
		}),
	},
	"next": gothrpc.Router{
		"test": &gothrpc.Procedure[any, string]{
			Query: gothrpc.QueryHandlerFn(func(ctx gothrpc.Context, args gothrpc.Args) (string, error) {
				return "whoa a next gen test fr", nil
			}),
			Mutation: gothrpc.MutationHandlerFn(func(ctx gothrpc.Context, args gothrpc.Args, p any) (string, error) {
				fmt.Printf("payload: %v\n", p)
				return "ok so this would imply that we did modify something, eh?", nil
			}),
		},
	},
	"props": &gothrpc.Procedure[any, any]{
		Query: gothrpc.QueryHandlerFn(func(ctx gothrpc.Context, args gothrpc.Args) (any, error) {
			return ctx.Props, nil
		}),
	},
	"panic": &gothrpc.Procedure[any, any]{
		Query: gothrpc.QueryHandlerFn(func(ctx gothrpc.Context, args gothrpc.Args) (any, error) {
			panic("test panic")
		}),
	},
}

func main() {

	const serverPort = "7774"
	const apiPrefix = "/api/rest/v1/"

	procHandler := &gothrpc.RestHandler{
		Router: procRouter,
		Prefix: apiPrefix,
		GetProps: func() any {
			return map[string]string{
				"test":      "ok",
				"some_data": "42",
			}
		},
		/*ErrorHandler: func(err error, ctx gothrpc.Context) {
			fmt.Printf("handler error: %s\n", err.Error())
		},*/
	}

	mux := http.NewServeMux()

	mux.Handle(apiPrefix, procHandler)

	fmt.Printf("Listening at: http://localhost:%s\n", serverPort)

	http.ListenAndServe(":"+serverPort, mux)
}
