package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"log"
	"net/http"
	"strings"
)

func GetMyhabits(w http.ResponseWriter, r *http.Request) {
	app, ok := injectApp(w, r)
	if !ok {
		return
	}

	// TODO remove hardcoded value
	habits, err := app.GetMyHabits(10, false)
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

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(res))
}

func PostMyHabits(w http.ResponseWriter, r *http.Request) {
	// TODO could probably make middleware for this
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content Type is not application/json")
		return
	}

	app, ok := injectApp(w, r)
	if !ok {
		return
	}

	newHabit := struct {
		Name        string
		Description string
		Frequency   int
	}{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newHabit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var unmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalErr) {
			fmt.Fprintf(w, "Bad Request. Wrong Type provided for field: %s", unmarshalErr.Field)
		} else {
			fmt.Fprintf(w, "Bad Request: %s", err)
		}
		return
	}

	habitId, err := app.CreateHabit(newHabit.Name, newHabit.Frequency)
	if err != nil {
		if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Input was not valid, %s", inputError)
			return
		}
		// given we block anonymous requests this error would be some internal issue
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something has gone wrong creating habit")
		log.Printf("Something has gone wrong creating habit: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, habitId)
}
