package gothrpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RestHandler struct {
	//	Source procedures
	Router Router
	//	HTTP path prefix
	Prefix string
	//	Hook that's called before request is passed to the router
	OnBeforeHandle HookPreHandlerFn
	//	Hook that's called after request is processed by the router
	OnAfterHandle HookPostHandlerFn
	//	Error handler callback (logger)
	OnError ErrorHandlerFn
}

type HookPreHandlerFn func(ctx *Context) error
type HookPostHandlerFn func(ctx *Context, result *RestResponse) error

type ErrorHandlerFn func(err error, ctx *Context)

type RestResponse struct {
	Data    any         `json:"data"`
	Error   *Error      `json:"error,omitempty"`
	Status  int         `json:"-"`
	Headers http.Header `json:"-"`
}

type Statuser interface {
	StatusCode() int
}

type Headerer interface {
	Headers() http.Header
}

func defaultErrorHandler(err error) {
	log.Default().Print("gothrpc error: ", err.Error())
}

func (this *RestHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//	todo: add CORS handler

	//	this avoids nil dereferencing later on
	if this.Router == nil {
		this.Router = Router{}
	}

	path := req.URL.Path
	if this.Prefix != "" {
		path = strings.TrimPrefix(path, this.Prefix)
	}

	ctx := &Context{
		Req:    req,
		Writer: writer,
		path:   newProcPath(path),
	}

	defer func() {

		if re := recover(); re != nil {

			err := Error{
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
			} else {
				defaultErrorHandler(err)
			}

			writeResponse(writer, wrapErrorResult(err))
		}
	}()

	if this.OnBeforeHandle != nil {
		if err := this.OnBeforeHandle(ctx); err != nil {
			writeResponse(writer, wrapErrorResult(err))
			return
		}
	}

	result := wrapHandlerResult(this.Router.Handle(ctx))

	if this.OnAfterHandle != nil {
		if err := this.OnAfterHandle(ctx, &result); err != nil {
			writeResponse(writer, wrapErrorResult(err))
			return
		}
	}

	writeResponse(writer, result)
}

func wrapDataResult(data any) RestResponse {

	result := RestResponse{
		Status: http.StatusOK,
		Data:   data,
	}

	if data == nil {
		result.Status = 204
	}

	if ext, ok := data.(Statuser); ok {
		result.Status = ext.StatusCode()
	}

	if ext, ok := data.(Headerer); ok {
		result.Headers = ext.Headers()
	}

	return result
}

func wrapErrorResult(err error) RestResponse {

	response := RestResponse{
		Status: http.StatusBadRequest,
	}

	if ext, valid := err.(Error); valid {
		response.Error = &ext
	} else {
		response.Error = &Error{
			Message: err.Error(),
		}
	}

	if ext, valid := err.(Headerer); valid {
		response.Headers = ext.Headers()
	}

	if ext, valid := err.(Statuser); valid {
		response.Status = ext.StatusCode()
	}

	return response
}

func wrapHandlerResult(data any, err error) RestResponse {

	if err != nil {
		return wrapErrorResult(err)
	}

	return wrapDataResult(data)
}

func writeResponse(writer http.ResponseWriter, response RestResponse) {

	if response.Headers != nil {
		for header, entry := range response.Headers {
			for _, value := range entry {
				writer.Header().Set(header, value)
			}
		}
	}

	if response.Status < http.StatusContinue {
		response.Status = http.StatusOK
	}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(response.Status)

	json.NewEncoder(writer).Encode(response)
}
