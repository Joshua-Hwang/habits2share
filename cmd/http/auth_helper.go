package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"

	"github.com/google/uuid"
)

const userIdKey = key("USER_ID")
const sessionCookieName = "__Secure-SESSIONID"
const sessionTtl = time.Duration(24 * 60 * 60 * 1000 * 1000 * 1000)
const cookieTtl = 365 * 24 * 60 * 60

func BuildGetLogin(webClientId string, redirectUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<html>
	<body>
		<script src="https://accounts.google.com/gsi/client" async defer></script>
		<div id="g_id_onload"
			data-client_id="%s"
			data-auto_prompt="false"
			data-login_uri="/login?redirect_url=%s"
			data-context="signin"
			data-ux_mode="redirect"
		></div>
		<div class="g_id_signin"
			data-type="standard"
			data-size="large"
			data-theme="outline"
			data-text="sign_in_with"
			data-shape="rectangular"
			data-logo_alignment="left"
		></div>
	</body>
</html>
`, webClientId, url.QueryEscape(redirectUrl))
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(authDbKey).(auth.AuthDatabase)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for AuthDatabase")
		return
	}

	tokenParser, ok := r.Context().Value(tokenParserKey).(auth.TokenParser)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Dependency injection failed")
		log.Printf("Dependency injection failed for TokenParser")
		return
	}

	const csrf_token = "g_csrf_token"
	csrfTokenBody := r.FormValue(csrf_token)
	csrfTokenCookie, err := r.Cookie(csrf_token)
	if csrfTokenBody == "" || err != nil || csrfTokenBody != csrfTokenCookie.Value {
		// status code from https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to verify double submit cookie")
		return
	}

	tokenString := r.FormValue("credential")
	claims, err := tokenParser.ParseToken(r.Context(), tokenString)
	// TODO throw single error message as security measure
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Authorization failed")
		return
	}

	// Find mapping of email to user id
	userId, err := db.GetUserIdFromEmail(r.Context(), claims.Email)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		fmt.Fprintf(w, "Failed to query habits: %s", err)
		log.Printf("Failed to make query: %s", err)
		return
	}

	sessionId := uuid.NewString()

	err = db.AddSession(r.Context(), sessionId, userId)
	if err != nil {
		log.Printf("Failed to add session to database: %s", err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionId,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   cookieTtl,
	})

	redirectUrl := r.URL.Query().Get("redirect_url")
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func buildSessionKey(sessionId string) string {
	return fmt.Sprintf("session/%s", sessionId)
}

// Parses session if it exists
func BuildSessionParser(redirectUrl string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			db, ok := r.Context().Value(authDbKey).(auth.AuthDatabase)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected error occurred")
				log.Printf("Dependency injection failed for AuthDatabase")
				return
			}

			authService, ok := r.Context().Value(authServiceKey).(*auth.AuthService)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected error occurred")
				log.Printf("Dependency injection failed for AuthService")
				return
			}

			sessionCookie, err := r.Cookie(sessionCookieName)

			if err != nil {
				if err != http.ErrNoCookie {
					log.Printf("Unexpected error with cookie: %s", err)
				}
				// no cookie thus anonymous
				next(w, r)
				return
			} else {
				sessionId, err := uuid.Parse(sessionCookie.Value)
				if err != nil {
					// parsing the cookie failed assume their anonymous
					next(w, r)
					return
				}
				userId, err := db.GetUserIdFromSession(
					r.Context(),
					sessionId.String(),
					time.Now().AddDate(-1, 0, 0),
				)

				if err != nil {
					if err == auth.ErrNotFound {
						// session could not be found
						// either doesn't exist or is past the due date
						// assume anonymous
						next(w, r)
						return
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						log.Fatalf("Scanning row failed: %s", err)
					}
				}

				ctx := context.WithValue(r.Context(), userIdKey, userId)
				r = r.WithContext(ctx)
				// hack as this is a new r and the existing closure doesn't capture it
				//authService.GetUserId = BuildUserIdGetter(r)
				authService.GetUserId =
					func() (string, error) {
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
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Fatalf("Storing user id: %s", err)
				}
			}

			next(w, r)
		}
	}
}

func BlockAnonymous(redirectPage http.HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	var blockingResponse http.HandlerFunc = redirectPage
	if blockingResponse == nil {
		blockingResponse = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Anonymous access forbidden")
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authService, ok := injectAuth(w, r)
		if !ok {
			return
		}

		if _, err := authService.GetCurrentUser(); err != nil {
			if err == habit_share.UserNotFoundError {
				blockingResponse(w, r)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected error accessing session info")
			}
		}

		app, ok := r.Context().Value(appKey).(*habit_share.App)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Dependency injection failed")
			log.Printf("Dependency injection failed for app")
			return
		}
		if _, err := app.Auth.GetCurrentUser(); err != nil {
			if err == habit_share.UserNotFoundError {
				blockingResponse(w, r)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected error accessing session info")
			}
		}

		next(w, r)
	}
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
