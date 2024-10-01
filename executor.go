package gothrpc

import (
	"errors"
	"fmt"
)

func execute(handler Handler, ctx Context) procedureResult {

	var data any
	var err error
	var isPanic bool

	func() {

		defer func() {

			if re := recover(); re != nil {

				isPanic = true

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

		data, err = handler.Handle(ctx)
	}()

	if err == nil {

		result := procedureResult{
			status: 200,
			Data:   data,
		}

		if data == nil {
			result.status = 204
		}

		if ext, ok := data.(Statuser); ok {
			result.status = ext.StatusCode()
		}
		if ext, ok := data.(Headerer); ok {
			result.header = ext.Headers()
		}

		return result
	}

	result := procedureResult{
		Error: &ProcError{},
	}

	if exterr, valid := err.(ProcError); valid {
		*result.Error = exterr
	} else {
		result.Error.Message = err.Error()
	}

	if ext, ok := err.(Statuser); ok {
		result.status = ext.StatusCode()
	} else if isPanic {
		result.status = 500
	} else {
		result.status = 400
	}

	if ext, ok := err.(Headerer); ok {
		result.header = ext.Headers()
	}

	return result
}
