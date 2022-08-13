package todo

import (
	"errors"
	"fmt"
	"time"
)

var TodoNotFoundError = errors.New("Todo could not be found")
// TODO habit sharing also has these errors
var UserNotFoundError = errors.New("User could not be found")
var PermissionDeniedError = errors.New("Operation was denied")

type InputError struct {
	Message string
}

var _ error = (*InputError)(nil)

// Error implements error
func (e *InputError) Error() string {
	return fmt.Sprintf("Failed to parse input because: %s", e.Message)
}

type AuthInterface interface {
	GetCurrentUser() (string, error)
}

type Todo struct {
	Id          string
	Owner       string
	Name        string
	Description string
	DueDate     time.Time
	Completed   bool
}

type TodoDatabase interface {
	CreateTodo(name string, owner string, dueDate time.Time) (string, error)
	ChangeName(id string, newName string) error
	ChangeDescription(id string, newDescription string) error
	ChangeDueDate(id string, newTime time.Time) error
	CompleteTodo(id string) error
	GetTodosByOwner(owner string, limit int, completed bool) ([]Todo, error)
	GetTodo(id string) (Todo, error)
}

type App struct {
	Db   TodoDatabase
	Auth AuthInterface
}

func (a *App) ownerCheck(todoId string) error {
	// TODO fetching the todo twice is hard
	todo, err := a.Db.GetTodo(todoId)
	if err != nil {
		return err
	}

	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return err
	}

	if todo.Owner != user {
		return PermissionDeniedError
	}

	return nil
}

func (a *App) CreateTodo(name string, dueDate time.Time) (string, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return "", err
	}

	return a.Db.CreateTodo(name, user, dueDate)
}

func (a *App) ChangeName(todoId string, newName string) error {
	if err := a.ownerCheck(todoId); err != nil {
		return err
	}

	return a.Db.ChangeName(todoId, newName)
}

func (a *App) ChangeDescription(todoId string, newDescription string) error {
	if err := a.ownerCheck(todoId); err != nil {
		return err
	}

	return a.Db.ChangeDescription(todoId, newDescription)
}

func (a *App) ChangeDueDate(todoId string, newTime time.Time) error {
	if err := a.ownerCheck(todoId); err != nil {
		return err
	}

	return a.Db.ChangeDueDate(todoId, newTime)
}

func (a *App) CompleteTodo(todoId string) error {
	if err := a.ownerCheck(todoId); err != nil {
		return err
	}

	return a.Db.CompleteTodo(todoId)
}

func (a *App) GetMyTodos(limit int, completed bool) ([]Todo, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return a.Db.GetTodosByOwner(user, limit, completed)
}

func (a *App) GetTodo(todoId string) (Todo, error) {
	if err := a.ownerCheck(todoId); err != nil {
		return Todo{}, err
	}

	return a.Db.GetTodo(todoId)
}
