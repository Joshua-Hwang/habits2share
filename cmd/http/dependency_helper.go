package http

import (
	"fmt"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"log"
	"net/http"
)

type key string

const appKey = key("APP")
const authDbKey = key("AUTH_DB")
const authServiceKey = key("AUTH_SERVICE")
const tokenParserKey = key("TOKEN_PARSER")

// helper function
func injectApp(w http.ResponseWriter, r *http.Request) (*habit_share.App, bool) {
	app, ok := r.Context().Value(appKey).(*habit_share.App)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for app")
		return nil, false
	}
	return app, true
}
