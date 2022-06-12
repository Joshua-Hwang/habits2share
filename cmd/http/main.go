package main

import (
	"context"
	"fmt"
	"internal/auth"
	"internal/auth_file"
	"internal/habit_share"
	"internal/habit_share_file"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	webClientId := os.Getenv("GOOGLE_WEB_CLIENT_ID")
	mobileClientId := os.Getenv("GOOGLE_MOBILE_CLIENT_ID")

	log.SetFlags(log.LstdFlags | log.Llongfile)

	authDatabase := &auth_file.AuthDatabaseFile{
		SessionsFilepath: "sessions.csv",
		AccountsFilepath: "accounts.json",
	}
	tokenParser := &auth.TokenParserGoogle{
		WebClientId:    webClientId,
		MobileClientId: mobileClientId,
	}

	habitsDatabase, err := habit_share_file.HabitShareFromFile("habits.json")
	if err != nil {
		panic(err)
	}

	// TODO move this to dependency helper
	buildDependencies := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = context.WithValue(ctx, authDbKey, authDatabase)
			ctx = context.WithValue(ctx, tokenParserKey, tokenParser)

			authService := &auth.AuthService{
				GetUserId: BuildUserIdGetter(r),
			}
			ctx = context.WithValue(ctx, authServiceKey, authService)

			app := &habit_share.App{Db: habitsDatabase, Auth: authService}
			ctx = context.WithValue(ctx, appKey, app)

			next(w, r.WithContext(ctx))
		}
	}

	// session parser requires dependencies to be built first
	sessionParser := BuildSessionParser("/login")

	mux := MuxWrapper{ServeMux: http.NewServeMux(), Middleware: buildDependencies}
	mux.RegisterHandlers("/login", map[string]http.HandlerFunc{
		"GET": BuildGetLogin(webClientId),
		"POST": PostLogin,
	})

	mux.RegisterHandlers("/my/habits", map[string]http.HandlerFunc{
		"GET": sessionParser(BlockAnonymous(GetMyhabits)),
		"POST": sessionParser(BlockAnonymous(PostMyHabits)),
	})
	// TODO if performance is an issue create an /all/habits
	mux.RegisterHandlers("/shared/habits", map[string]http.HandlerFunc{
		"GET": sessionParser(BlockAnonymous(GetSharedHabits)),
	});

	// TODO if performance is an issue return activities in same batch as habits
	// How do you keep all this modular without burdening the client?
	// GraphQL provides one such way. An exposed and powerful querying language
	// the clients are able to use in a single request.
	{
		pathPrefix := "/habit/"
		mux.HandleFunc(pathPrefix, sessionParser(BlockAnonymous(func(w http.ResponseWriter, r *http.Request) {
			// get habitId
			remainingUrl := strings.TrimPrefix(r.URL.EscapedPath(), pathPrefix)
			if remainingUrl == r.URL.EscapedPath() {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Path is malformed")
				log.Printf("Something funky going on with trimming path")
			}
			habitId, remainingUrl, _ := strings.Cut(remainingUrl, "/")
			// the slash is removed during cut
			remainingUrl = fmt.Sprintf("/%s?%s", remainingUrl, r.URL.Query().Encode())
			r.URL, _ = url.Parse(remainingUrl)

			app, ok := injectApp(w, r)
			if !ok {
				return
			}
			// TODO check user is allowed to see habit prior
			habit, err := app.GetHabit(habitId)
			if err != nil {
				if err == habit_share.HabitNotFoundError {
					http.NotFound(w, r)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to retrieve habit %v", err)
				log.Printf("Failed to retrieve habit %v", err)
			}

			habitHandler := BuildHabitHandler(&habit)
			habitHandler.ServeHTTP(w, r)
		})));
	}

	{
		//pathPrefix := "/user/"
	}
	// POST /user/:userId in body provide habit to share
	// DELETE /user/:userId/habit/:habitId to unshare app

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
