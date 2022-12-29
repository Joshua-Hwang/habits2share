package main

import (
	"fmt"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"net/http"
	"strings"
)

func (s Server) PostUserHabit(w http.ResponseWriter, r *http.Request) {
	var userId string
	var habitId string
	var err error
	// first split is an empty string because we start with /
	splits := strings.SplitN(r.URL.EscapedPath(), "/", 5)
	if len(splits) != 5 || splits[1] != "user" || splits[3] != "habit" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User id or habit id invalid")
		return
	}
	userId = splits[2]
	habitId = splits[4]

	// check if user exists
	if found, err := s.AuthDatabase.UserExists(r.Context(), userId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to share the habit")
		return
	} else if !found {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User doesn't exist")
		return
	}

	reqDeps, err := s.BuildRequestDependenciesOrReject(w, r)
	if err != nil {
		return
	}
	app := reqDeps.HabitApp

	err = app.ShareHabit(habitId, userId)
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

func (s Server) DeleteUserHabit(w http.ResponseWriter, r *http.Request) {
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

	reqDeps, err := s.BuildRequestDependenciesOrReject(w, r)
	if err != nil {
		return
	}
	app := reqDeps.HabitApp

	err = app.UnShareHabit(habitId, userId)
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
