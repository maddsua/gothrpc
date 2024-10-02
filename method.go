package gothrpc

type Method struct {
	GET    Handler
	POST   Handler
	PUT    Handler
	DELETE Handler
}

func (this *Method) Handle(ctx *Context) (any, error) {

	if ctx.procPath.HasNext() {
		return nil, errProcNotFound
	}

	var useHandler Handler

	switch ctx.Req.Method {
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
		return nil, errMethodNotAllowed
	}

	return useHandler.Handle(ctx)
}
