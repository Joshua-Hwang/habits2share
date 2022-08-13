package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

func GetMyTodos(w http.ResponseWriter, r *http.Request) {
	app, ok := injectTodoApp(w, r)
	if !ok {
		return
	}

	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		limitString = "64"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Limit query is in incorrect, must be an integer")
	}

	todos, err := app.GetMyTodos(limit, false)
	if err != nil && err != todo.UserNotFoundError {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "GetMyTodos failed")
		log.Printf("GetMyTodos failed with %v", err)
	}

	res, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshalling failed")
		log.Printf("Marshalling failed with %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(res))
}

func PostMyTodos(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintf(w, "Content Type is not application/json")
		return
	}

	app, ok := injectTodoApp(w, r)
	if !ok {
		return
	}

	newTodo := struct {
		Name        string
		Description string
		DueDate     time.Time
	}{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&newTodo)
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

	todoId, err := app.CreateTodo(newTodo.Name, newTodo.DueDate)
	if err != nil {
		if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Input was not valid, %s", inputError)
			return
		}
		// given we block anonymous requests this error would be some internal issue
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something has gone wrong creating todo")
		log.Printf("Something has gone wrong creating todo: %v", err)
		return
	}

	// TODO messy as habit creation is no longer atomic. Please fix
	err = app.ChangeDescription(todoId, newTodo.Description)
	if err != nil {
		if inputError := (*todo.InputError)(nil); errors.As(err, &inputError) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Input was not valid, %s", inputError)
			return
		}
		// given we block anonymous requests this error would be some internal issue
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something has gone wrong updating description of todo")
		log.Printf("Something has gone wrong description of todo: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, todoId)
}
