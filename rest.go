package gothrpc

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewRestContext(req *http.Request, ctx *RestHandler) Context {

	path := req.URL.Path
	if ctx.Prefix != "" {
		path = strings.TrimPrefix(path, ctx.Prefix)
	}

	return Context{
		Req: req,
		procPath: procStepper{
			segments: strings.Split(strings.TrimPrefix(path, "/"), "/"),
		},
		Props: ctx.Props,
	}
}

type RestHandler struct {
	Router Router
	Props  any
	Prefix string
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: also add methods to handle CORS and stuff

	ctx := NewRestContext(req, this)
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
