package rest

import (
	"encoding/json"
	"goproc/proc"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func parseArgs(req *http.Request) proc.Args {

	args := proc.Args{}

	args = parseUrlArgs(args, req.URL)

	switch req.Method {
	case "POST", "PUT", "PATCH":
	default:
		return args
	}

	if !strings.Contains(req.Header.Get("content-type"), "json") {
		return args
	}

	if data, err := io.ReadAll(req.Body); err == nil {
		args = parseBodyArgs(args, data)
	}

	return args
}

func parseUrlArgs(args proc.Args, url *url.URL) proc.Args {

	if args == nil {
		args = map[string]any{}
	}

	for key, entries := range url.Query() {

		if len(entries) == 0 {
			continue
		}

		args[key] = entries[len(entries)-1]
	}

	return args
}

func parseBodyArgs(args proc.Args, body []byte) proc.Args {

	if args == nil {
		args = map[string]any{}
	}

	payload := map[string]any{}

	if err := json.Unmarshal(body, &payload); err != nil {
		return args
	}

	for key, value := range payload {
		args[key] = value
	}

	return args
}
