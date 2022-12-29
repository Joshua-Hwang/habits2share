package main

import (
	"fmt"
	"net/http"
)

const RemainingPathKey = key("REMAINING_PATH")

type MethodHandlers map[string]http.HandlerFunc

type MuxWrapper struct {
	*http.ServeMux
}

var _ http.Handler = (*MuxWrapper)(nil)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method is not allowed")
}

func (mux *MuxWrapper) RegisterHandlers(path string, methodHandlers MethodHandlers) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			// TODO add tracing to log messages
			handler, exist := methodHandlers[r.Method]
			if !exist {
				methodNotAllowedHandler(w, r)
				return
			}

			handler(w, r)
		},
	)
}
