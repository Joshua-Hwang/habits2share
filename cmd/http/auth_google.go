//go:build !dev

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (s Server) GetLogin(w http.ResponseWriter, r *http.Request) {
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
`, s.TokenParser.WebClientId, "/web")
	// hardcoded redirect. Think about more general when we need to
	// I'll justify saying it's really the frontends job to properly log in the user.
}

func (s Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	db := s.AuthDatabase
	tokenParser := s.TokenParser

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
