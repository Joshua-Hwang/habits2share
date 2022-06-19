package http

import (
	"fmt"
	"net/http"
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

func (mux *MuxWrapper) Handle(path string, handler http.Handler) {
	mux.ServeMux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			if mux.Middleware != nil {
				handler = mux.Middleware(handler.ServeHTTP)
			}

			handler.ServeHTTP(w, r)
		},
	)
}

func (mux *MuxWrapper) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	mux.ServeMux.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			if mux.Middleware != nil {
				handler = mux.Middleware(handler)
			}

			handler(w, r)
		},
	)
}

func (mux *MuxWrapper) RegisterHandlers(path string, methodHandlers map[string]http.HandlerFunc) {
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
