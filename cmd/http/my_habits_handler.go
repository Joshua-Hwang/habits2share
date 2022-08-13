package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share_import"
)

func GetMyhabits(w http.ResponseWriter, r *http.Request) {
	app, ok := injectApp(w, r)
	if !ok {
		return
	}

	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		limitString = "64"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Limit query is in incorrect, must be an integer")
	}

	habits, err := app.GetMyHabits(limit, false)
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

	// TODO messy as habit creation is no longer atomic. Please fix
	err = app.ChangeDescription(habitId, newHabit.Description)
	if err != nil {
		if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Input was not valid, %s", inputError)
			return
		}
		// given we block anonymous requests this error would be some internal issue
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something has gone wrong updating description of habit")
		log.Printf("Something has gone wrong description of habit habit: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, habitId)
}

func PostMyHabitsImport(w http.ResponseWriter, r *http.Request) {
	// TODO prevent someone from uploading too much
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/csv") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content Type is not text/csv")
		return
	}

	db, ok := injectDb(w, r);
	if !ok {
		return
	}

	authService, ok := injectAuth(w, r);
	if !ok {
		return
	}

	csvReader := csv.NewReader(r.Body)

	habitIds, err := habit_share_import.ImportCsv(db, authService, csvReader)
	if err != nil {
		if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unable to import csv")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something has gone wrong importing habits")
		log.Printf("Something has gone wrong importing habits: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	for _, habitId := range habitIds {
		fmt.Fprintf(w, "%s\n", habitId)
	}
}
