package gothrpc

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewRestContext(req *http.Request, props any) Context {

	pathSegments := strings.Split(strings.TrimPrefix(strings.ReplaceAll(req.URL.Path, "//", "/"), "/"), "/")

	return Context{
		Req: req,
		procPath: procStepper{
			segments: pathSegments,
		},
		Props: props,
	}
}

type RestHandler struct {
	Router Router
	Props  any
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: also add methods to handle CORS and stuff

	ctx := NewRestContext(req, this.Props)
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
