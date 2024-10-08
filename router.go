package gothrpc

type Router map[string]Handler

func (this Router) Handle(ctx *Context) (any, error) {

	procname, hasProcname := ctx.path.next()
	if !hasProcname {
		return nil, errProcNotFound
	}

	proc, has := this[procname]
	if !has {
		return nil, errProcNotFound
	}

	return proc.Handle(ctx)
}
