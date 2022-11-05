package habit_share

import (
	"errors"
	"fmt"
)

type Habit struct {
	Id          string
	Owner       string
	SharedWith  map[string]struct{}
	Name        string
	Description string
	Frequency   int
	Archived    bool
}

type Activity struct {
	Id      string // an Id is technically not needed but we let the database provider decide
	HabitId string
	Logged  Time
	Status  string
}

var ActivityNotFoundError = errors.New("Activity could not be found")
var HabitNotFoundError = errors.New("Habit could not be found")
var UserNotFoundError = errors.New("User could not be found")

var PermissionDeniedError = errors.New("Operation was denied")

type InputError struct {
	StringToParse string
}

var _ error = (*InputError)(nil)

// Error implements error
func (e *InputError) Error() string {
	return fmt.Sprintf("Failed to parse input, input was %s", e.StringToParse)
}

const (
	ActivitySuccess = "SUCCESS"
	ActivityMinimum = "MINIMUM"
	ActivityNotDone = "NOT_DONE"
)

/*
Thoughts on the current API.
So I made this reflecting on Clean Architecture. I defined, from the habit
share's perspective, the actions necessary to run.
This API is quite verbose though.

The operations to change different attributes of the Habit could probably be reduced.
Don't forget habit_share defines the struct. The database providers should know about the struct.
Additionally we're replicating the data on either side of this API boundary (not a huge deal given how ephemeral the habit_share side is).
*/
type HabitsDatabase interface {
	// Not sure this is a good idea. Instead to create a habit struct and the habit id is populated for you and also returned
	CreateHabit(newHabit Habit) (string, error)
	ShareHabit(habitId string, friend string) error
	UnShareHabit(habitId string, friend string) error
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetMyHabits(owner string, limit int, archived bool) ([]Habit, error)
	// this should not show archived habits
	GetSharedHabits(owner string, limit int) ([]Habit, error)
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetHabit(id string) (Habit, error)

	SetHabit(habitId string, updatedHabit Habit) error;

	DeleteHabit(id string) error

	CreateActivity(habitId string, logged Time, status string) (string, error)
	GetActivities(habitId string, after Time, before Time, limit int) (activities []Activity, hasMore bool, err error)
	DeleteActivity(habitId, id string) error

	GetScore(habitId string) (int, error)
}
