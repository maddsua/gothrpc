package gothrpc

import (
	"errors"
	"fmt"
)

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

// todo: consider moving this into "executor"
func (this Router) Exec(ctx Context) Result {

	var data any
	var err error

	func() {

		defer func() {

			if re := recover(); re != nil {

				switch err.(type) {
				case error:
					err = re.(error)
					if ctx.errorHandler != nil {
						ctx.errorHandler(err, ctx)
					}

				default:
					err = errors.New("runtime error")
					if ctx.errorHandler != nil {
						ctx.errorHandler(errors.New(fmt.Sprintf("%v", re)), ctx)
					}
				}
			}
		}()

		data, err = this.Handle(ctx)

	}()

	if err == nil {

		//	todo: write different codes based on returned content
		//	eg: return 204 for successful responses with no data
		//	idk if that's a good idea tho. we'll see

		result := Result{
			status: 200,
			Data:   data,
		}

		if ext, ok := data.(Statuser); ok {
			result.status = ext.StatusCode()
		}
		if ext, ok := data.(Headerer); ok {
			result.header = ext.Headers()
		}

		return result
	}

	result := Result{
		//	todo: detect runtime errors and set code to 500
		status: 400,
	}

	result.Error = &ProcError{}
	if exterr, valid := err.(ProcError); valid {
		*result.Error = exterr
	} else {
		result.Error.Message = err.Error()
	}

	if ext, ok := err.(Statuser); ok {
		result.status = ext.StatusCode()
	}
	if ext, ok := err.(Headerer); ok {
		result.header = ext.Headers()
	}

	return result
}
