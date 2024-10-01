package gothrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RestHandler struct {
	Router       Router
	BeforeHandle BeforeHandleHookFn
	AfterHandle  AfterHandleHookFn
	Prefix       string
	ErrorHandler func(err error, ctx Context)
}

type BeforeHandleHookFn func(ctx *Context) error

type AfterHandleHookFn func(ctx *Context) (*AfterHandleHookResult, error)

type AfterHandleHookResult struct {
	StatusCode int
	Headers    http.Header
}

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
	//	todo: get status and headers from hook return values

	if this.BeforeHandle != nil {
		if err := this.BeforeHandle(&ctx); err != nil {
			writeResponse(writer, Result{
				Error: &ProcError{
					Message: err.Error(),
				},
			})
			return
		}
	}

	if this.ErrorHandler != nil {
		ctx.errorHandler = this.ErrorHandler
	} else {
		ctx.errorHandler = defaultErrorHandler
	}

	result := execute(this.Router, ctx)

	if result.header != nil {
		writeHeaders(writer, result.header)
	}

	if this.AfterHandle != nil {

		hookResult, err := this.AfterHandle(&ctx)
		if err != nil {
			writeResponse(writer, Result{
				Error: &ProcError{
					Message: err.Error(),
				},
			})
			return
		}

		if hookResult != nil {

			if hookResult.Headers != nil {
				writeHeaders(writer, hookResult.Headers)
			}

			if hookResult.StatusCode > http.StatusOK {
				result.status = hookResult.StatusCode
			}

		}
	}

	writeResponse(writer, result)
}

func writeResponse(writer http.ResponseWriter, result Result) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(result.status)
	json.NewEncoder(writer).Encode(result)
}

func writeHeaders(writer http.ResponseWriter, headers http.Header) {
	for header, entry := range headers {
		for _, value := range entry {
			writer.Header().Set(header, value)
		}
	}
}
