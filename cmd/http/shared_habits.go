package main

import (
	"encoding/json"
	"fmt"
	"internal/habit_share"
	"log"
	"net/http"
)

func GetSharedHabits(w http.ResponseWriter, r *http.Request) {
	app, ok := injectApp(w, r)
	if !ok {
		return
	}

	habits, err := app.GetSharedHabits(10)
	if err != nil && err != habit_share.UserNotFoundError {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "GetMyHabits failed")
		log.Printf("GetMyHabits failed with %v", err)
	}

	res, err := json.Marshal(habits)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshalling failed")
		log.Printf("Marshalling failed with %v", err)
	}
	fmt.Fprint(w, string(res))
}
