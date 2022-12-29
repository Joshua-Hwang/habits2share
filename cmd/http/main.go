package main

import (
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
	config := GetGlobalConfig()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	// check necessary files exist before continuing.
	// It's either here or we discover during initial testing in production
	// if necessary files don't exist crash
	if _, err := os.Stat(config.accountsFilePath); err != nil {
		log.Fatalf("accountsFilePath points to \"%s\" which does not exist", config.accountsFilePath)
	}

	authDatabase := &auth_file.AuthDatabaseFile{
		SessionsFilepath: config.sessionFilePath,
		SessionsFileLock: &sync.RWMutex{},
		AccountsFilepath: config.accountsFilePath,
	}
	tokenParser := &auth.TokenParserGoogle{
		WebClientId:    config.webClientId,
		MobileClientId: config.mobileClientId,
	}

	habitsDatabase, err := habit_share_file.HabitShareFromFile(config.habitsFilePath)
	if err != nil {
		panic(err)
	}

	todoDatabase, err := todo_file.TodoFromFile(config.todoFilePath)
	if err != nil {
		panic(err)
	}

	// Hopefully it's sufficiently clear that this isn't all the dependencies
	server := Server{
		GlobalDependencies{
			AuthDatabase:   authDatabase,
			TokenParser:    tokenParser,
			HabitsDatabase: habitsDatabase,
			TodoDatabase:   todoDatabase,
		},
	}

	mux := MuxWrapper{ServeMux: http.NewServeMux()}

	mux.RegisterHandlers("/healthcheck", MethodHandlers{
		// TODO make this more like a deepcheck
		"HEAD": func(w http.ResponseWriter, r *http.Request) {
			// Uses the HEAD command for uptime service I use but can be changed easily
			return
		},
	})

	// TODO because session parsing happens here. All static files are checked for login credentials causing unnecessary reads
	// No authentication is done on the Single Page App. It should be their responsibility to get the user logged in.
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./frontend/build"))))

	mux.RegisterHandlers("/login", MethodHandlers{
		"GET":  server.GetLogin,
		"POST": server.PostLogin,
	})

	mux.RegisterHandlers("/my/habits", MethodHandlers{
		"GET":  server.GetMyHabits,
		"POST": server.PostMyHabits,
	})
	mux.RegisterHandlers("/my/habits/upload", MethodHandlers{
		"POST": server.PostMyHabitsImport,
	})

	// NOTE if performance is an issue consider creating an /all/habits
	mux.RegisterHandlers("/shared/habits", MethodHandlers{
		"GET": server.GetSharedHabits,
	})

	// NOTE if performance is an issue return activities in same batch as habits
	// How do you keep all this modular without burdening the client?
	// GraphQL provides one such way. An exposed and powerful querying language
	// the clients are able to use in a single request.
	{
		pathPrefix := "/habit/"
		mux.Handle(pathPrefix, http.StripPrefix(pathPrefix, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var err error
				// get habitId
				habitId, remainingUrl, _ := strings.Cut(r.URL.EscapedPath(), "/")
				// the slash is removed during cut
				remainingUrl = fmt.Sprintf("/%s?%s", remainingUrl, r.URL.Query().Encode())
				r.URL, _ = url.Parse(remainingUrl)

				reqDeps, err := server.BuildRequestDependenciesOrReject(w, r)
				if err != nil {
					return
				}
				app := reqDeps.HabitApp

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

				habitHandler := reqDeps.BuildHabitHandler(&habit)
				habitHandler.ServeHTTP(w, r)
			}),
		))
	}

	// NOTE this doesn't work if other endpoints exist on this prefix
	mux.RegisterHandlers("/user/", MethodHandlers{
		"POST":   server.PostUserHabit,
		"DELETE": server.DeleteUserHabit,
	})

	mux.RegisterHandlers("/my/todos", MethodHandlers{
		"GET":  server.GetMyTodos,
		"POST": server.PostMyTodos,
	})

	{
		pathPrefix := "/todo/"
		mux.Handle(pathPrefix, http.StripPrefix(pathPrefix, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				var err error
				// get todoId
				todoId, remainingUrl, _ := strings.Cut(r.URL.EscapedPath(), "/")
				// the slash is removed during cut
				remainingUrl = fmt.Sprintf("/%s?%s", remainingUrl, r.URL.Query().Encode())
				r.URL, _ = url.Parse(remainingUrl)

				reqDeps, err := server.BuildRequestDependenciesOrReject(w, r)
				if err != nil {
					return
				}
				app := reqDeps.TodoApp

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

				habitHandler := reqDeps.BuildTodoHandler(&todoItem)
				habitHandler.ServeHTTP(w, r)
			}),
		))
	}

	log.Printf("Listening on port %s", config.port)
	log.Printf("Process ID %d", os.Getpid())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.port), mux))
}
