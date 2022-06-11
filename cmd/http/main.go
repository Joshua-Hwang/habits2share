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
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	webClientId := os.Getenv("GOOGLE_WEB_CLIENT_ID")
	mobileClientId := os.Getenv("GOOGLE_MOBILE_CLIENT_ID")

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

	// POST to /habit/:habitId with status in body to register an activity
	// GET to /habit/:habitId?limit=...&order=... works on the pagination of activities
	// this endpoint also returns description, name and shared with
	// GET /habit/:habitid/streak to get the number

	// TODO need a concept of user Id
	// POST /user/userId

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
