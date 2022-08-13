package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

type key string

const appKey = key("APP")
const todoAppKey = key("TODO_APP")
const dbKey = key("DB")
const todoDbKey = key("TODO_DB")
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

func injectDb(w http.ResponseWriter, r *http.Request) (habit_share.HabitsDatabase, bool) {
	db, ok := r.Context().Value(dbKey).(habit_share.HabitsDatabase)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for db")
		return nil, false
	}
	return db, true
}

func injectAuth(w http.ResponseWriter, r *http.Request) (habit_share.AuthInterface, bool) {
	auth, ok := r.Context().Value(authServiceKey).(habit_share.AuthInterface)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for auth")
		return nil, false
	}
	return auth, true
}

func injectTodoApp(w http.ResponseWriter, r *http.Request) (*todo.App, bool) {
	app, ok := r.Context().Value(todoAppKey).(*todo.App)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for todo app")
		return nil, false
	}
	return app, true
}
