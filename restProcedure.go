package gothrpc

type Procedure struct {
	Query    Handler
	Mutation Handler
}

func (this *Procedure) Handle(ctx Context) (any, error) {

	var useHandler Handler

	switch ctx.Method {
	case "GET":
		useHandler = this.Query
	case "POST":
		useHandler = this.Mutation
	}

	if useHandler == nil {
		return nil, ErrorMethodNotAllowed
	}

	return useHandler.Handle(ctx)
}
