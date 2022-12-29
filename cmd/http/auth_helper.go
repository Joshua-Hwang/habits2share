package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"

	"github.com/google/uuid"
)

const userIdKey = key("USER_ID")
const sessionCookieName = "__Secure-SESSIONID"
const sessionTtl = time.Duration(24 * 60 * 60 * 1000 * 1000 * 1000)
const cookieTtl = 365 * 24 * 60 * 60

func buildSessionKey(sessionId string) string {
	return fmt.Sprintf("session/%s", sessionId)
}

func (s Server) BuildAuthService(r *http.Request) (*auth.AuthService, error) {
	db := s.AuthDatabase
	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("Unexpected error with cookie: %s", err)
		}
		return nil, err
	}
	sessionId, err := uuid.Parse(sessionCookie.Value)
	if err != nil {
		return nil, err
	}
	userId, err := db.GetUserIdFromSession(
		r.Context(),
		sessionId.String(),
		time.Now().AddDate(-1, 0, 0),
	)
	if err != nil {
		return nil, err
	}
	authService := &auth.AuthService{
		UserId: userId,
	}
	return authService, nil
}

// This is duplicate code for convenience sake. Time will tell if this was an appropriate change to the API
func (s Server) BuildAuthServiceOrReject(w http.ResponseWriter, r *http.Request) (*auth.AuthService, error) {
	authService, err := s.BuildAuthService(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Anonymous access forbidden")
		return nil, err
	}
	return authService, nil
}

func BuildUserIdGetter(r *http.Request) func() (string, error) {
	return func() (string, error) {
		val := r.Context().Value(userIdKey)
		if val == nil {
			return "", nil
		}
		userId, ok := val.(string)
		if !ok {
			panic("userIdKey did not return string")
		}
		return userId, nil
	}
}
