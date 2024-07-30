package main

import (
	"goproc/proc"
	"goproc/rest"
	"net/http"
)

type errorResult map[string]any

func (this errorResult) StatusCode() int {
	return 400
}
func (this errorResult) Headers() http.Header {
	return http.Header{}
}

// todo: add procedure handlers with type definitions
var procRouter = proc.Router{
	"test": &rest.Method{
		GET: proc.HandleFunc(func(ctx proc.Context) (any, error) {
			return map[string]any{
				"args":    ctx.Args,
				"message": "well this didn't do much!",
			}, nil
		}),
		POST: proc.HandleFunc(func(ctx proc.Context) (any, error) {
			return "whoops!", nil
		}),
		DELETE: proc.HandleFunc(func(ctx proc.Context) (any, error) {
			return errorResult{"test": "test errors"}, nil
		}),
	},
	"next": proc.Router{
		"test": &rest.Method{
			GET: proc.HandleFunc(func(ctx proc.Context) (any, error) {
				return "whoa a next gen test fr", nil
			}),
		},
	},
}

func main() {

	procHandler := &rest.Handler{
		Router: procRouter,
	}

	mux := http.NewServeMux()

	mux.Handle("/", procHandler)

	http.ListenAndServe(":7774", mux)
}
