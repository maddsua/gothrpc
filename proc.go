package gothrpc

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Procedure[P any, Q any, M any] struct {
	Query    QueryHandler[Q]
	Mutation MutationHandler[P, M]
}

func (this *Procedure[P, Q, M]) Handle(ctx *Context) (any, error) {

	if ctx.path.hasNext() {
		return nil, errProcNotFound
	}

	switch ctx.Req.Method {
	case "GET":
		return this.handleQuery(ctx)
	case "POST":
		return this.handleMutation(ctx)
	}

	return nil, errMethodNotAllowed
}

func (this *Procedure[P, Q, M]) handleQuery(ctx *Context) (any, error) {

	if this.Query == nil {
		return nil, errMethodNotAllowed
	}

	return this.Query.Handle(ctx, getQueryInput(ctx.Req))
}

func (this *Procedure[P, Q, M]) handleMutation(ctx *Context) (any, error) {

	if this.Mutation == nil {
		return nil, errMethodNotAllowed
	}

	var payload P

	if strings.Contains(ctx.Req.Header.Get("content-type"), "json") {

		//	construct some of the nil-by-default types
		switch any(payload).(type) {
		case QueryInput:
			payload = any(QueryInput{}).(P)

		case map[string]any:
			payload = any(map[string]any{}).(P)
		}

		//	parse json payload
		if data, err := io.ReadAll(ctx.Req.Body); err == nil {
			if err := json.Unmarshal(data, &payload); err != nil {
				return nil, Error{
					Message: "failed to unwrap mutation props",
					Extensions: map[string]any{
						"cause": err.Error(),
					},
				}
			}
		}

	} else {
		//	extract input from URL search aprams
		switch any(payload).(type) {
		case QueryInput:
			payload = any(getQueryInput(ctx.Req)).(P)
		}
	}

	return this.Mutation.Handle(ctx, payload)
}

func getQueryInput(req *http.Request) QueryInput {

	input := map[string]string{}

	for key, entries := range req.URL.Query() {

		if len(entries) == 0 {
			continue
		}

		input[key] = entries[len(entries)-1]
	}

	return input
}
