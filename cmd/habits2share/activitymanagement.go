package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"internal/auth_http"
	server "internal/baseserver"
	"internal/data"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getActivityHandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	habit := HabitRow{}

	remainingPath, ok := r.Context().Value(server.RemainingPathKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("RemainingPathKey did not contain a string")
	}
	if id, err := uuid.Parse(remainingPath); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request: %s", err)
		return
	} else {
		habit.Id = id
	}

	// pagination
	params := r.URL.Query()
	limit := 7
	offset := 0
	before := time.Now()
	if value, err := strconv.ParseUint(params.Get("limit"), 10, 32); err == nil {
		if value <= 365 {
			limit = int(value)
		}
	}
	if value, err := strconv.ParseUint(params.Get("offset"), 10, 32); err == nil {
		offset = int(value)
	}
	if value, err := time.Parse(time.RFC3339, params.Get("before")); err == nil {
		before = value
	}

	habit, err := getHabitForUser(r.Context(), habit.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Habit wasn't found: %s", err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Scanning row failed: %s", err)
		}
		return
	}

	activities, err := getActivities(r.Context(), habit, before, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query habits: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	type ResponseJson struct {
		Habit      HabitRow      `json:"habit"`
		Activities []ActivityRow `json:"activities"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseJson{Habit: habit, Activities: activities})
}

func postActivityHandler(w resw, r *req) {
	// add or remove days (or the half complete)
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	habit := HabitRow{}
	{
		remainingPath, ok := r.Context().Value(server.RemainingPathKey).(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request: could not parse path")
			return
		}
		if id, err := uuid.Parse(remainingPath); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request: %s", err)
			return
		} else {
			habit.Id = id
		}

		var err error
		habit, err = getHabitOwnedByUser(r.Context(), habit.Id, userId)

		// disallow modification of habits shared to you
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Habit wasn't found: %s", err)
			return
		}
	}

	activity := ActivityRow{}
	{
		headerContentTtype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(headerContentTtype, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprintf(w, "Content Type is not application/json")
			return
		}

		newActivity := struct {
			Logged       time.Time `json:"logged"` // TODO Only accepts RFC3339 time format but make it accept ISO8061
			LoggedStatus string    `json:"logged_status"`
		}{}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&newActivity)
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
		if newActivity.Logged.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request. Fields missing")
			return
		}
		activity.habit = &habit
		activity.Logged = newActivity.Logged
		activity.LoggedStatus = newActivity.LoggedStatus
	}

	{
		// TODO optimisation could be done here
		// (only perform complex calculation when inserting between old values)
		txn, err := data.GetDb().BeginTx(r.Context(), nil)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Transaction failed: %s", err)
			return
		}

		if activity.LoggedStatus != "" {
			_, err = addActivity(r.Context(), txn, activity.habit.Id, activity.Logged, activity.LoggedStatus)
			if err != nil {
				txn.Rollback()
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Insertion failed: %s", err)
				return
			}
		} else {
			_, err = deleteActivity(r.Context(), txn, activity.habit.Id, activity.Logged)
			if err != nil {
				txn.Rollback()
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Failed to delete activity: %s", err)
				return
			}
		}

		_, err = updateStreak(r.Context(), txn, activity.habit.Id)
		if err != nil {
			txn.Rollback()
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Streak update failed: %s", err)
			return
		}

		if err := txn.Commit(); err != nil {
			txn.Rollback()
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Transaction failed: %s", err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// Deprecated
func deleteActivityHandler(w resw, r *req) {
	// add or remove days (or the half complete)
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	remainingPath, ok := r.Context().Value(server.RemainingPathKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("RemainingPathKey did not contain a string")
	}
	urlParts := strings.SplitN(remainingPath, "/", 2)

	habit := HabitRow{}
	{
		if id, err := uuid.Parse(urlParts[0]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request: %s", err)
			return
		} else {
			habit.Id = id
		}

		var err error
		habit, err = getHabitOwnedByUser(r.Context(), habit.Id, userId)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Habit wasn't found: %s", err)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatalf("Scanning row failed: %s", err)
			}
			return
		}
	}

	txn, err := data.GetDb().BeginTx(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Unable to start transaction: %s", err)
		return
	}
	// delete date
	var activityDate time.Time
	if date, err := time.Parse("2006-01-02", urlParts[1]); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request: %s", err)
		return
	} else {
		activityDate = date
	}
	_, err = deleteActivity(r.Context(), txn, habit.Id, activityDate)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to delete activity: %s", err)
		return
	}

	_, err = updateStreak(r.Context(), txn, habit.Id)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Streak update failed: %s", err)
		return
	}

	if err := txn.Commit(); err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Transaction failed: %s", err)
		return
	}
}
