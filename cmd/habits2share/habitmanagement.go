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

	"github.com/google/uuid"
)

func getHabitsHandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	habits, err := getHabitsForUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query habits: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habits)
}

func getSharedHabitsHandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	habits, err := getHabitsShared(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query habits: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(habits)
}

func postHabitsHandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(headerContentTtype, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content Type is not application/json")
		return
	}

	newHabit := struct {
		Id          uuid.NullUUID `json:"id"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Frequency   int           `json:"frequency"`
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

	// Frequency 0 is considered sentinel
	if !(newHabit.Frequency == 0 || (1 <= newHabit.Frequency && newHabit.Frequency <= 7)) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request: Frequency must be 1 to 7 inclusive not %d", newHabit.Frequency)
		return
	}

	txn, err := data.GetDb().BeginTx(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Unable to start transaction: %s", err)
		return
	}

	if newHabit.Id.Valid {
		if newHabit.Name != "" {
			_, err = renameHabitForUser(r.Context(), txn, newHabit.Name, newHabit.Id.UUID, userId)
			if err != nil {
				txn.Rollback()
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Update failed: %s", err)
				return
			}
		}

		// This implies every description must contain data. You can't delete the description.
		if newHabit.Description != "" {
			_, err = describeHabitForUser(r.Context(), txn, newHabit.Description, newHabit.Id.UUID, userId)
			if err != nil {
				txn.Rollback()
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Update failed: %s", err)
				return
			}
		}

		if newHabit.Frequency != 0 {
			_, err = changeFrequencyHabitForUser(r.Context(), txn, newHabit.Frequency, newHabit.Id.UUID, userId)
			if err != nil {
				txn.Rollback()
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprintf(w, "Update failed: %s", err)
				return
			}
		}
	} else {
		if newHabit.Frequency == 0 {
			newHabit.Frequency = 7
		}
		_, err = addHabit(r.Context(), txn, newHabit.Name, newHabit.Description, newHabit.Frequency, userId)
		if err != nil {
			txn.Rollback()
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Insertion failed: %s", err)
			return
		}
	}

	if err := txn.Commit(); err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Transaction failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func postShareHandler(w resw, r *req) {
	// json with habit (you own) and who you're sharing with
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(headerContentTtype, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content Type is not application/json")
		return
	}

	newShare := struct {
		Habit      uuid.NullUUID `json:"habit"`
		SharedWith string        `json:"shared_with"`
	}{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newShare)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var unmarshalErr *json.UnmarshalTypeError
		if errors.As(err, &unmarshalErr) {
			fmt.Fprintf(w, "Bad Request: Wrong Type provided for field: %s", unmarshalErr.Field)
		} else {
			fmt.Fprintf(w, "Bad Request: %s", err)
		}
		return
	}
	if !newShare.Habit.Valid {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request: habit is not valid UUID")
		return
	}

	// not in a transaction so risk of race conditions
	habit, err := getHabitForUser(r.Context(), newShare.Habit.UUID, userId)
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

	txn, err := data.GetDb().BeginTx(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Unable to start transaction: %s", err)
		return
	}

	_, err = shareHabit(r.Context(), txn, habit.Id, newShare.SharedWith)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Insertion failed: %s", err)
		return
	}

	if err := txn.Commit(); err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Transaction failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func deleteShareHandler(w resw, r *req) {
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
	habitId := urlParts[0]
	sharedWith := urlParts[1]

	habit := HabitRow{}
	{
		if id, err := uuid.Parse(habitId); err != nil {
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

	_, err = deleteShare(r.Context(), txn, habit.Id, sharedWith)
	if err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to delete habit: %s", err)
		return
	}

	if err := txn.Commit(); err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Transaction failed: %s", err)
		return
	}
}

func deleteHabitHandler(w resw, r *req) {
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

	habit := HabitRow{}
	{
		if id, err := uuid.Parse(remainingPath); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request: %s", err)
			return
		} else {
			habit.Id = id
		}

		var err error
		habit, err = getHabitForUser(r.Context(), habit.Id, userId)
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

	// delete (archive) habit
	params := r.URL.Query()
	archive := true
	if value, err := strconv.ParseBool(params.Get("archive")); err == nil {
		// since delete is meant to be idempotent we need to specify permanent
		// deletion with a query parameter
		archive = value
	}

	txn, err := data.GetDb().BeginTx(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Unable to start transaction: %s", err)
		return
	}

	if archive {
		_, err := archiveHabit(r.Context(), txn, habit.Id)
		if err != nil {
			txn.Rollback()
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Failed to archive habit: %s", err)
			return
		}
	} else {
		_, err := deleteHabit(r.Context(), txn, habit.Id)
		if err != nil {
			txn.Rollback()
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "Failed to delete habit: %s", err)
			return
		}
	}

	if err := txn.Commit(); err != nil {
		txn.Rollback()
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Transaction failed: %s", err)
		return
	}
}
