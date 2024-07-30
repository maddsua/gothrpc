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
		GET: gothrpc.HandleFunc(func(ctx gothrpc.Context) (any, error) {
			return map[string]any{
				"args":    ctx.Args,
				"message": "well this didn't do much!",
			}, nil
		}),
		POST: gothrpc.HandleFunc(func(ctx gothrpc.Context) (any, error) {
			return "whoops!", nil
		}),
		DELETE: gothrpc.HandleFunc(func(ctx gothrpc.Context) (any, error) {
			return errorResult{"test": "test errors"}, nil
		}),
	},
	"next": gothrpc.Router{
		"test": &gothrpc.Method{
			GET: gothrpc.HandleFunc(func(ctx gothrpc.Context) (any, error) {
				return "whoa a next gen test fr", nil
			}),
		},
	},
}

func main() {

	const serverPort = "7774"

	procHandler := &gothrpc.RestHandler{
		Router: procRouter,
	}

	mux := http.NewServeMux()

	mux.Handle("/", procHandler)

	fmt.Printf("Listening at: http://localhost:%s\n", serverPort)

	http.ListenAndServe(":"+serverPort, mux)
}
