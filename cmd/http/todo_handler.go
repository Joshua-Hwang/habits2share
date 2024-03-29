package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

// We build this way because the mux can't accept the todoItem beyond the two parameters
// Ideally the mux should take any function and perform redirection based on path.
// A more advanced version of a map.
// Is it a failure of the language to prevent me from doing the simpler thing?
// Under the hood it's just pointers and it could be made that the mux doesn't run the operation.
// mux accepts anonymous functions which makes closures. Can't hot swap the todoItem
// Maybe it's a failure of the library that I can separate the routing logic from the mux further.
// What if we raise all this to the main mux? won't work because todoId is present in the URL
// I think the correct approach is to abandon the mux provided by the http package and run our own
// Ideas are welcome
func (reqDeps RequestDependencies) BuildTodoHandler(todoItem *todo.Todo) http.Handler {
	mux := MuxWrapper{ServeMux: http.NewServeMux()}
	mux.RegisterHandlers("/", map[string]http.HandlerFunc{
		"POST": func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				fmt.Fprintf(w, "Content Type is not application/json")
				return
			}

			app := reqDeps.TodoApp

			updatePayload := struct {
				Name        string
				Description string
				DueDate     time.Time
				Completed   bool
			}{}
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&updatePayload)
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

			if updatePayload.Name != "" {
				err = app.ChangeName(todoItem.Id, updatePayload.Name)
			}
			if err != nil {
				if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Name was invalid. Name, description, due date and completion not modified")
				} else if errors.Is(err, todo.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this todo")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed rename todo")
					log.Printf("Something has gone wrong renaming todo: %v", err)
				}
				return
			}

			if updatePayload.Description != "" {
				err = app.ChangeDescription(todoItem.Id, updatePayload.Description)
			}
			if err != nil {
				if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Description was invalid. Description, due date and completion not modified")
				} else if errors.Is(err, todo.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this todo")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to change description")
					log.Printf("Something has gone wrong changing description: %v", err)
				}
				return
			}

			if !updatePayload.DueDate.IsZero() {
				err = app.ChangeDueDate(todoItem.Id, updatePayload.DueDate)
			}
			if err != nil {
				if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, DueDate was invalid. Due date and completion not modified")
				} else if errors.Is(err, todo.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this todo")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to change due date")
					log.Printf("Something has gone wrong changing due date: %v", err)
				}
				return
			}

			// TODO since you can't tell if a a boolean is not specific or false this is the best we got
			if updatePayload.Completed {
				err = app.CompleteTodo(todoItem.Id)
			}
			if err != nil {
				if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Completion was invalid. Completion was not modified")
				} else if errors.Is(err, todo.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this todo")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed completing todo")
					log.Printf("Something has gone wrong completing todo: %v", err)
				}
				return
			}

			w.WriteHeader(http.StatusCreated)
		},
	})

	return mux
}
