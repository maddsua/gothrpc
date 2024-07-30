package rest

import "goproc/proc"

type Method struct {
	GET    proc.Handler
	POST   proc.Handler
	PUT    proc.Handler
	DELETE proc.Handler
}

func (this *Method) Handle(ctx proc.Context) (any, error) {

	var useHandler proc.Handler

	switch ctx.Method {
	case "GET":
		useHandler = this.GET
	case "POST":
		useHandler = this.POST
	case "PUT":
		useHandler = this.PUT
	case "DELETE":
		useHandler = this.DELETE
	}

	if useHandler == nil {
		return nil, proc.ErrorMethodNotAllowed
	}

	return useHandler.Handle(ctx)
}
