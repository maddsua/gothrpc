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

type BeforeHandleHookFn func(req *http.Request, ctx *Context) error

type AfterHandleHookFn func(ctx *Context) (*AfterHandleHookResult, error)

type AfterHandleHookResult struct {
	StatusCode int
	Headers    http.Header
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
		//	todo: fix
		err := this.BeforeHandle(req, &ctx)
	}

	if this.ErrorHandler != nil {
		ctx.errorHandler = this.ErrorHandler
	} else {
		ctx.errorHandler = func(err error, _ Context) {
			log.Default().Print("gothrpc error: ", err.Error())
		}
	}

	result := execute(this.Router, ctx)

	if result.Headers() != nil {
		for header, entry := range result.Headers() {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	writer.Header().Set("content-type", "application/json")

	if this.AfterHandle != nil {
		//	todo: fix
		hookResult, err := this.AfterHandle(&ctx)
	}

	writer.WriteHeader(result.StatusCode())

	json.NewEncoder(writer).Encode(result)
}
