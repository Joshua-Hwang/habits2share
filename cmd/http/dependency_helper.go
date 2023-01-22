package main

import (
	"net/http"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
	"github.com/Joshua-Hwang/habits2share/pkg/auth_file"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share_file"
	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

type key string
type UserIdType string

type HabitAppInterface interface {
	ArchiveHabit(id string) error
	ChangeDescription(id string, newDescription string) error
	ChangeFrequency(id string, newFrequency int) error
	ChangeName(id string, newName string) error
	CreateActivity(habitId string, logged habit_share.Time, status string) (string, error)
	CreateHabit(name string, frequency int) (string, error)
	DeleteActivity(habitId string, id string) error
	DeleteHabit(id string) error
	GetActivities(habitId string, after habit_share.Time, before habit_share.Time, limit int) (activities []habit_share.Activity, hasMore bool, err error)
	GetHabit(id string) (habit_share.Habit, error)
	GetMyHabits(limit int, archived bool) ([]habit_share.Habit, error)
	GetScore(habitId string) (int, error)
	GetSharedHabits(limit int) ([]habit_share.Habit, error)
	ShareHabit(habitId string, friend string) error
	UnShareHabit(habitId string, friend string) error
}

type TodoAppInterface interface {
	ChangeDescription(todoId string, newDescription string) error 
	ChangeDueDate(todoId string, newTime time.Time) error 
	ChangeName(todoId string, newName string) error
	CompleteTodo(todoId string) error 
	CreateTodo(name string, dueDate time.Time) (string, error) 
	GetMyTodos(limit int, completed bool) ([]todo.Todo, error) 
	GetTodo(todoId string) (todo.Todo, error)
}

// Contains app scoped dependencies
type Server struct {
	// Nothing else is in this struct. Dependencies is here purely for semantics
	GlobalDependencies
}

type GlobalDependencies struct {
	// For both performance and mutexes these are here
	AuthDatabase   *auth_file.AuthDatabaseFile
	TokenParser    *auth.TokenParserGoogle
	HabitsDatabase *habit_share_file.HabitShareFile
	TodoDatabase   todo.TodoDatabase
}

// TODO probably worth splitting, not very performant
// Maybe results in ball of mud?
type RequestDependencies struct {
	GlobalDependencies
	AuthService *auth.AuthService
	HabitApp    HabitAppInterface
	TodoApp     TodoAppInterface
}

func (s Server) BuildRequestDependenciesOrReject(w http.ResponseWriter, r *http.Request) (*RequestDependencies, error) {
	authService, err := s.BuildAuthServiceOrReject(w, r)
	if err != nil {
		return nil, err
	}
	habitApp := s.BuildHabitApp(authService)
	todoApp := s.BuildTodoApp(authService)

	requestDependencies := RequestDependencies{
		GlobalDependencies: s.GlobalDependencies,
		AuthService:        authService,
		HabitApp:           habitApp,
		TodoApp:            todoApp,
	}

	return &requestDependencies, nil
}

// TODO If InitHabits needs more services we get caught at compile time
// The intention is to eventually use wire
func (s Server) BuildHabitApp(
	authService habit_share.AuthInterface,
) *habit_share.App {
	return &habit_share.App{Db: s.HabitsDatabase, Auth: authService}
}

func (s Server) BuildTodoApp(
	authService todo.AuthInterface,
) *todo.App {
	return &todo.App{Db: s.TodoDatabase, Auth: authService}
}
