package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const RemainingPathKey = key("REMAINING_PATH")

type MuxWrapper struct {
	*http.ServeMux
	Middleware func(http.HandlerFunc) http.HandlerFunc
}

var _ http.Handler = (*MuxWrapper)(nil)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method is not allowed")
}

func (mux MuxWrapper) RegisterHandlers(path string, methodHandlers map[string]http.HandlerFunc) {
	mux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			val := strings.TrimPrefix(r.URL.EscapedPath(), path)
			if val == r.URL.EscapedPath() {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Path is malformed")
				log.Printf("Something funky going on with trimming path")
			}
			ctx := context.WithValue(r.Context(), RemainingPathKey, val)
			r = r.WithContext(ctx)

			// TODO add tracing to log messages
			handler, exist := methodHandlers[r.Method]
			if !exist {
				methodNotAllowedHandler(w, r)
				return
			}

			if mux.Middleware != nil {
				handler = mux.Middleware(handler)
			}

			handler(w, r)
		},
	)
}
