//go:build dev

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
		DEV MODE LOGIN PAGE
		 <form>
			<label for="email">email</label><br>
			<input type="email" id="email" name="email"><br>
			<input formmethod="post" type="submit" value="Submit">
		</form>
	</body>
</html>
`)
}

func (s Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	db := s.AuthDatabase

	email := r.FormValue("email")

	// Find mapping of email to user id
	userId, err := db.GetUserIdFromEmail(r.Context(), email)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		fmt.Fprintf(w, "User does not exist: %s", err)
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
