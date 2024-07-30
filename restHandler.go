package gothrpc

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RestHandler struct {
	Router Router
	Props  any
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	pathSegments := strings.Split(strings.TrimPrefix(strings.ReplaceAll(req.URL.Path, "//", "/"), "/"), "/")

	ctx := Context{
		Method:     req.Method,
		Headers:    req.Header,
		RemoteAddr: req.RemoteAddr,
		ProcPath:   NewProcedureStepper(pathSegments),
		Args:       parseArgs(req),
		Props:      this.Props,
	}

	result := this.Router.Exec(ctx)

	//	todo: error writer

	writer.WriteHeader(result.Status())

	if result.Header() != nil {
		for header, entry := range result.Header() {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	writer.Header().Set("content-type", "application/json")

	json.NewEncoder(writer).Encode(result)
}
