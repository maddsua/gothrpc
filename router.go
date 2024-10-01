package gothrpc

type Router map[string]Handler

func (this Router) Handle(ctx Context) (any, error) {

	procname, hasProcname := ctx.procPath.Next()
	if !hasProcname {
		return nil, ErrorProcedureNotFound
	}

	proc, has := this[procname]
	if !has {
		return nil, ErrorProcedureNotFound
	}

	return proc.Handle(ctx)
}
