package gothrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RestHandler struct {
	Router       Router
	BeforeHandle HookHandlerFn
	Prefix       string
	ErrorHandler func(err error, ctx Context)
}

type HookHandlerFn func(ctx *Context) error

func defaultErrorHandler(err error, _ Context) {
	log.Default().Print("gothrpc error: ", err.Error())
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: also add methods to handle CORS and stuff

	path := req.URL.Path
	if this.Prefix != "" {
		path = strings.TrimPrefix(path, this.Prefix)
	}

	ctx := Context{
		Req:      req,
		procPath: newStepper(path),
	}

	//	todo: defer panic recover
	if this.BeforeHandle != nil {
		if err := this.BeforeHandle(&ctx); err != nil {
			writeErrorResponse(writer, err)
			return
		}
	}

	if this.ErrorHandler != nil {
		ctx.errorHandler = this.ErrorHandler
	} else {
		ctx.errorHandler = defaultErrorHandler
	}

	result := execute(this.Router, ctx)
	writeResponse(writer, result)
}

func writeResponse(writer http.ResponseWriter, result procedureResult) {

	if result.header != nil {
		for header, entry := range result.header {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(result.status)

	json.NewEncoder(writer).Encode(result)
}

func writeErrorResponse(writer http.ResponseWriter, err error) {

	result := procedureResult{
		Error: &ProcError{
			Message: err.Error(),
		},
	}

	if ext, ok := err.(Headerer); ok {
		result.header = ext.Headers()
	}

	if ext, ok := err.(Statuser); ok {
		result.status = ext.StatusCode()
	}

	writeResponse(writer, result)
}
