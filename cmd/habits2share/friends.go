package main

import (
	"encoding/json"
	"fmt"
	"internal/auth_http"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// This could be made more complicated in the future
func getFriendshandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	accounts, err := getAccounts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query accounts: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	friends := make(map[string][]uuid.UUID, 0)
	for _, account := range accounts {
		friends[account.Email] = make([]uuid.UUID, 0)
	}

	sharedRows, err := getHabitsSharedByUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query shared habits: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	for _, shared := range sharedRows {
		friends[shared.Email] = append(friends[shared.Email], shared.Habit)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}

func getAccountshandler(w resw, r *req) {
	userId := auth.ReadUserId(r.Context())
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Account required")
		return
	}

	accounts, err := getAccounts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Failed to query accounts: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	emailToName := make(map[string]string, 0)
	for _, account := range accounts {
		emailToName[account.Email] = account.Name
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emailToName)
}
