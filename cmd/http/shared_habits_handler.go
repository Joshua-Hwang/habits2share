package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
)

func (s Server) GetSharedHabits(w http.ResponseWriter, r *http.Request) {
	requestDependencies, err := s.BuildRequestDependenciesOrReject(w, r)
	if err != nil {
		return
	}
	app := requestDependencies.HabitApp

	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		limitString = "64"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Limit query is in incorrect, must be an integer")
	}

	habits, err := app.GetSharedHabits(limit)
	if err != nil && err != habit_share.UserNotFoundError {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "GetMyHabits failed")
		log.Printf("GetMyHabits failed with %v", err)
	}

	// TODO the SharedWith shouldn't be exposed in what is shared
	res, err := json.Marshal(habits)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshalling failed")
		log.Printf("Marshalling failed with %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(res))
}
