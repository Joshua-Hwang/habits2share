package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
	"github.com/Joshua-Hwang/habits2share/pkg/auth_file"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share_file"
	"github.com/Joshua-Hwang/habits2share/pkg/todo"
	"github.com/Joshua-Hwang/habits2share/pkg/todo_file"
)

func main() {
	webClientId := os.Getenv("GOOGLE_WEB_CLIENT_ID")
	mobileClientId := os.Getenv("GOOGLE_MOBILE_CLIENT_ID")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	sessionFilePath := os.Getenv("SESSIONS_FILE")
	if sessionFilePath == "" {
		sessionFilePath = "sessions.csv"
	}
	accountsFilePath := os.Getenv("ACCOUNTS_FILE")
	if accountsFilePath == "" {
		accountsFilePath = "accounts.json"
	}
	habitsFilePath := os.Getenv("HABITS_FILE")
	if habitsFilePath == "" {
		habitsFilePath = "habits.json"
	}
	todoFilePath := os.Getenv("TODO_FILE")
	if todoFilePath == "" {
		todoFilePath = "todo.json"
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	authDatabase := &auth_file.AuthDatabaseFile{
		SessionsFilepath: sessionFilePath,
		SessionsFileLock: &sync.RWMutex{},
		AccountsFilepath: accountsFilePath,
	}
	tokenParser := &auth.TokenParserGoogle{
		WebClientId:    webClientId,
		MobileClientId: mobileClientId,
	}

	habitsDatabase, err := habit_share_file.HabitShareFromFile(habitsFilePath)
	if err != nil {
		panic(err)
	}

	todoDatabase, err := todo_file.TodoFromFile(todoFilePath)
	if err != nil {
		panic(err)
	}

	// TODO move this to dependency helper
	buildDependencies := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = context.WithValue(ctx, dbKey, habitsDatabase)
			ctx = context.WithValue(ctx, todoDbKey, todoDatabase)
			ctx = context.WithValue(ctx, authDbKey, authDatabase)
			ctx = context.WithValue(ctx, tokenParserKey, tokenParser)

			authService := &auth.AuthService{
				GetUserId: BuildUserIdGetter(r),
			}
			ctx = context.WithValue(ctx, authServiceKey, authService)

			app := &habit_share.App{Db: habitsDatabase, Auth: authService}
			ctx = context.WithValue(ctx, appKey, app)

			todoApp := &todo.App{Db: todoDatabase, Auth: authService}
			ctx = context.WithValue(ctx, todoAppKey, todoApp)

			next(w, r.WithContext(ctx))
		}
	}

	// session parser requires dependencies to be built first
	sessionParser := BuildSessionParser("/login")

	mux := MuxWrapper{ServeMux: http.NewServeMux(), Middleware: buildDependencies}

	mux.RegisterHandlers("/healthcheck", map[string]http.HandlerFunc{
		// TODO make this more like a deepcheck
		"HEAD": func(w http.ResponseWriter, r *http.Request) {
			// Uses the HEAD command for uptime service I use but can be changed easily
			return
		},
	})

	// These assets need to skip the credentials check because Firefox doesn't
	// send cookies when requesting manifest.json
	mux.HandleFunc("/web/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/build/manifest.json")
	})
	mux.HandleFunc("/web/asset-manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/build/asset-manifest.json")
	})

	// TODO because session parsing happens here. All static files are checked for login credentials causing unnecessary reads
	mux.Handle("/web/", sessionParser(BlockAnonymous(
		BuildGetLogin(webClientId, "/web/"),
		http.StripPrefix("/web/", http.FileServer(http.Dir("./frontend/build"))).ServeHTTP,
	)))

	mux.RegisterHandlers("/login", map[string]http.HandlerFunc{
		"GET":  BuildGetLogin(webClientId, "/web"),
		"POST": PostLogin,
	})

	mux.RegisterHandlers("/my/habits", map[string]http.HandlerFunc{
		"GET":  sessionParser(BlockAnonymous(nil, GetMyhabits)),
		"POST": sessionParser(BlockAnonymous(nil, PostMyHabits)),
	})
	mux.RegisterHandlers("/my/habits/upload", map[string]http.HandlerFunc{
		"POST": sessionParser(BlockAnonymous(nil, PostMyHabitsImport)),
	})

	// NOTE if performance is an issue consider creating an /all/habits
	mux.RegisterHandlers("/shared/habits", map[string]http.HandlerFunc{
		"GET": sessionParser(BlockAnonymous(nil, GetSharedHabits)),
	})

	// NOTE if performance is an issue return activities in same batch as habits
	// How do you keep all this modular without burdening the client?
	// GraphQL provides one such way. An exposed and powerful querying language
	// the clients are able to use in a single request.
	{
		pathPrefix := "/habit/"
		mux.Handle(pathPrefix, http.StripPrefix(pathPrefix,
			sessionParser(BlockAnonymous(nil, func(w http.ResponseWriter, r *http.Request) {
				// get habitId
				habitId, remainingUrl, _ := strings.Cut(r.URL.EscapedPath(), "/")
				// the slash is removed during cut
				remainingUrl = fmt.Sprintf("/%s?%s", remainingUrl, r.URL.Query().Encode())
				r.URL, _ = url.Parse(remainingUrl)

				app, ok := injectApp(w, r)
				if !ok {
					return
				}

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
			})),
		))
	}

	// NOTE this doesn't work if other endpoints exist on this prefix
	mux.RegisterHandlers("/user/", map[string]http.HandlerFunc{
		"POST":   sessionParser(BlockAnonymous(nil, PostUserHabit)),
		"DELETE": sessionParser(BlockAnonymous(nil, DeleteUserHabit)),
	})

	mux.RegisterHandlers("/my/todos", map[string]http.HandlerFunc{
		"GET":  sessionParser(BlockAnonymous(nil, GetMyTodos)),
		"POST": sessionParser(BlockAnonymous(nil, PostMyTodos)),
	})

	{
		pathPrefix := "/todo/"
		mux.Handle(pathPrefix, http.StripPrefix(pathPrefix,
			sessionParser(BlockAnonymous(nil, func(w http.ResponseWriter, r *http.Request) {
				// get habitId
				todoId, remainingUrl, _ := strings.Cut(r.URL.EscapedPath(), "/")
				// the slash is removed during cut
				remainingUrl = fmt.Sprintf("/%s?%s", remainingUrl, r.URL.Query().Encode())
				r.URL, _ = url.Parse(remainingUrl)

				app, ok := injectTodoApp(w, r)
				if !ok {
					return
				}

				todoItem, err := app.GetTodo(todoId)
				if err != nil {
					if err == todo.TodoNotFoundError {
						http.NotFound(w, r)
						return
					}
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to retrieve habit %v", err)
					log.Printf("Failed to retrieve habit %v", err)
				}

				habitHandler := BuildTodoHandler(&todoItem)
				habitHandler.ServeHTTP(w, r)
			})),
		))
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
