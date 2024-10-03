package gothrpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Statuser interface {
	StatusCode() int
}

type Headerer interface {
	Headers() http.Header
}

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

	args := procedureGetArgs(ctx)

	return this.Query.Handle(ctx, args)
}

func (this *Procedure[P, Q, M]) handleMutation(ctx *Context) (any, error) {

	if this.Mutation == nil {
		return nil, errMethodNotAllowed
	}

	var payload P

	//	fail if P is a concrete type and body is empty
	if strings.Contains(ctx.Req.Header.Get("content-type"), "json") {
		if data, err := io.ReadAll(ctx.Req.Body); err == nil {
			if err := json.Unmarshal(data, &payload); err != nil {
				return nil, errors.New("failed to parse mutation props")
			}
		}
	}

	args := procedureGetArgs(ctx)

	return this.Mutation.Handle(ctx, args, payload)
}

func procedureGetArgs(ctx *Context) Args {

	args := map[string]string{}

	for key, entries := range ctx.Req.URL.Query() {

		if len(entries) == 0 {
			continue
		}

		args[key] = entries[len(entries)-1]
	}

	return args
}
