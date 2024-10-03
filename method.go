package gothrpc

type Method struct {
	GET    Handler
	POST   Handler
	PUT    Handler
	DELETE Handler
}

func invokeMethodHandler(method Handler, ctx *Context) (any, error) {

	if method == nil {
		return nil, errMethodNotAllowed
	}

	return method.Handle(ctx)
}

func (this *Method) Handle(ctx *Context) (any, error) {

	if ctx.path.hasNext() {
		return nil, errProcNotFound
	}

	switch ctx.Req.Method {
	case "GET":
		return invokeMethodHandler(this.GET, ctx)
	case "POST":
		return invokeMethodHandler(this.POST, ctx)
	case "PUT":
		return invokeMethodHandler(this.PUT, ctx)
	case "DELETE":
		return invokeMethodHandler(this.DELETE, ctx)
	default:
		return nil, errMethodNotAllowed
	}
}
