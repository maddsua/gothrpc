package rest

import (
	"encoding/json"
	"goproc/proc"
	"net/http"
)

type Handler struct {
	Router proc.Router
	Props  any
}

func (this *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	ctx := NewContext(req, this.Props)

	result := this.Router.Exec(ctx)

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
