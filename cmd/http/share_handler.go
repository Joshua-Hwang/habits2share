package main

import (
	"fmt"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"net/http"
	"strings"
)

func PostUserHabit(w http.ResponseWriter, r *http.Request) {
	var userId string
	var habitId string
	// first split is an empty string because we start with /
	splits := strings.SplitN(r.URL.EscapedPath(), "/", 5)
	if len(splits) != 5 || splits[1] != "user" || splits[3] != "habit" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User id or habit id invalid")
		return
	}
	userId = splits[2]
	habitId = splits[4]

	app, ok := injectApp(w, r)
	if !ok {
		// remember injection already returns the failure header
		return
	}

	err := app.ShareHabit(habitId, userId)
	if err != nil {
		if err == habit_share.PermissionDeniedError {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not allowed to share that habit")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to share the habit")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteUserHabit(w http.ResponseWriter, r *http.Request) {
	var userId string
	var habitId string
	// first split is an empty string because we start with /
	splits := strings.SplitN(r.URL.EscapedPath(), "/", 5)
	if len(splits) != 5 || splits[1] != "user" || splits[3] != "habit" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User id or habit id invalid")
		return
	}
	userId = splits[2]
	habitId = splits[4]

	app, ok := injectApp(w, r)
	if !ok {
		// remember injection already returns the failure header
		return
	}

	err := app.UnShareHabit(habitId, userId)
	if err != nil {
		if err == habit_share.PermissionDeniedError {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Not allowed to share that habit")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to share the habit")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
