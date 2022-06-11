package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	server "internal/baseserver"
)

type resw = http.ResponseWriter
type req = http.Request

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO middleware for timeouts and rate limiting
	mux := server.BuildBaseServer()

	mux.RegisterHandlers("/habits", map[string]http.HandlerFunc{
		"GET":  server.ParseSession(getHabitsHandler),
		"POST": server.ParseSession(postHabitsHandler),
	})

	mux.RegisterHandlers("/habit/", map[string]http.HandlerFunc{
		"GET":    server.ParseSession(getActivityHandler),
		"PUT":    server.ParseSession(postActivityHandler),
		"POST":   server.ParseSession(postActivityHandler),
		"DELETE": server.ParseSession(deleteHandler),
	})

	mux.RegisterHandlers("/shared/habits", map[string]http.HandlerFunc{
		"GET":  server.ParseSession(getSharedHabitsHandler),
		"POST": server.ParseSession(postShareHandler),
	})

	mux.RegisterHandlers("/shared/habits/", map[string]http.HandlerFunc{
		"DELETE": server.ParseSession(deleteShareHandler),
	})

	mux.RegisterHandlers("/friends", map[string]http.HandlerFunc{
		"GET": server.ParseSession(getFriendshandler),
	})

	// TODO THIS IS A TERRIBLE SOLUTION
	mux.RegisterHandlers("/accounts", map[string]http.HandlerFunc{
		"GET": server.ParseSession(getAccountshandler),
	})

	fs := http.FileServer(http.Dir("frontend"))
	mux.RegisterHandlers("/", map[string]http.HandlerFunc{
		"GET": fs.ServeHTTP,
	})
	mux.RegisterHandlers("/#/", map[string]http.HandlerFunc{
		"GET": http.StripPrefix("/#", fs).ServeHTTP,
	})

	log.Printf("Listening on port: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func deleteHandler(w resw, r *req) {
	remainingPath, ok := r.Context().Value(server.RemainingPathKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("RemainingPathKey did not contain a string")
	}
	urlParts := strings.SplitN(remainingPath, "/", 2)

	if len(urlParts) == 1 {
		deleteHabitHandler(w, r)
	} else {
		deleteActivityHandler(w, r)
	}
}
