package rest

import (
	"goproc/proc"
	"net/http"
	"strings"
)

func NewContext(req *http.Request, props any) proc.Context {

	pathSegments := strings.Split(strings.TrimPrefix(strings.ReplaceAll(req.URL.Path, "//", "/"), "/"), "/")

	return proc.Context{
		Method:     req.Method,
		Headers:    req.Header,
		RemoteAddr: req.RemoteAddr,
		ProcPath:   proc.NewProcedureStepper(pathSegments),
		Args:       parseArgs(req),
		Props:      props,
	}
}
