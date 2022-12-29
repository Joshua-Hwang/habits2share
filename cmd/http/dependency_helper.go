package main

import (
	"net/http"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
	"github.com/Joshua-Hwang/habits2share/pkg/auth_file"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share_file"
	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

type key string
type UserIdType string

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
	HabitApp    *habit_share.App
	TodoApp     *todo.App
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
