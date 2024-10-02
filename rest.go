package gothrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RestHandler struct {
	Router         Router
	Prefix         string
	OnBeforeHandle HookHandlerFn
	OnAfterHandle  HookHandlerFn
	OnError        ErrorHandlerFn
}

// todo: pass response.writer hooks
type HookHandlerFn func(ctx *Context) error
type ErrorHandlerFn func(err error, ctx Context)

func defaultErrorHandler(err error, _ Context) {
	log.Default().Print("gothrpc error: ", err.Error())
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: also add methods to handle CORS and stuff

	path := req.URL.Path
	if this.Prefix != "" {
		path = strings.TrimPrefix(path, this.Prefix)
	}

	ctx := Context{
		Req:      req,
		procPath: newStepper(path),
	}

	defer func() {

		if re := recover(); re != nil {

			err := ProcError{
				HttpStatus: http.StatusInternalServerError,
			}

			switch re.(type) {
			case error:
				err.Message = re.(error).Error()

			case string:
				err.Message = re.(string)

			default:
				err.Message = "runtime error"
			}

			if this.OnError != nil {
				this.OnError(err, ctx)
			}

			writeErrorResponse(writer, err)
		}
	}()

	if this.OnBeforeHandle != nil {
		if err := this.OnBeforeHandle(&ctx); err != nil {
			writeErrorResponse(writer, err)
			return
		}
	}

	result, err := this.Router.Handle(ctx)
	if err != nil {
		writeErrorResponse(writer, err)
		return
	}

	//	todo: add onAfterHandle

	writeDataResponse(writer, result)
}

func writeResponse(writer http.ResponseWriter, response procResult) {

	if response.header != nil {
		for header, entry := range response.header {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(response.status)

	json.NewEncoder(writer).Encode(response)
}

func writeErrorResponse(writer http.ResponseWriter, err error) {

	response := procResult{
		Error: &ProcError{
			Message: err.Error(),
		},
	}

	if ext, ok := err.(Headerer); ok {
		response.header = ext.Headers()
	}

	if ext, ok := err.(Statuser); ok {
		response.status = ext.StatusCode()
	}

	writeResponse(writer, response)
}

func writeDataResponse(writer http.ResponseWriter, result any) {

	response := procResult{
		status: 200,
		Data:   result,
	}

	if result == nil {
		response.status = 204
	}

	if ext, ok := result.(Statuser); ok {
		response.status = ext.StatusCode()
	}
	if ext, ok := result.(Headerer); ok {
		response.header = ext.Headers()
	}

	writeResponse(writer, response)
}
