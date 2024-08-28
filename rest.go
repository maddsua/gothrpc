package gothrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RestHandler struct {
	Router       Router
	GetProps     func() any
	Prefix       string
	ErrorHandler func(err error, ctx Context)
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: also add methods to handle CORS and stuff

	path := req.URL.Path
	if this.Prefix != "" {
		path = strings.TrimPrefix(path, this.Prefix)
	}

	ctx := Context{
		Req: req,
		procPath: procStepper{
			segments: strings.Split(strings.TrimPrefix(path, "/"), "/"),
		},
	}

	if this.GetProps != nil {
		ctx.Props = this.GetProps()
	}

	if this.ErrorHandler != nil {
		ctx.errorHandler = this.ErrorHandler
	} else {
		ctx.errorHandler = func(err error, _ Context) {
			log.Default().Print("gothrpc error: ", err.Error())
		}
	}

	result := this.Router.Exec(ctx)

	//	todo: error writer

	if result.Header() != nil {
		for header, entry := range result.Header() {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	writer.Header().Set("content-type", "application/json")

	writer.WriteHeader(result.Status())

	json.NewEncoder(writer).Encode(result)
}
