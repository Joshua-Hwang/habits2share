package habit_share

import (
	"errors"
	"fmt"
	"time"
)

type Habit struct {
	Id         string
	Owner      string
	SharedWith map[string]struct{}
	Name       string
	Frequency  int
	Archived   bool
}

type Activity struct {
	Id      string // an Id is technically not needed but we let the database provider decide
	HabitId string
	Logged  time.Time
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

type HabitsDatabase interface {
	CreateHabit(name string, owner string, frequency int) (string, error)
	ShareHabit(habitId string, friend string) error
	UnShareHabit(habitId string, friend string) error
	GetSharedWith(habitId string) (map[string]struct{}, error)
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetMyHabits(owner string, limit int, archived bool) ([]Habit, error)
	// this should not show archived habits
	GetSharedHabits(owner string, limit int) ([]Habit, error)
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetHabit(id string) (Habit, error)
	ChangeName(id string, newName string) error
	// frequency should technically be checked (1-7) in this part prior to sending
	// request to underlying implementation
	ChangeFrequency(id string, newFrequency int) error
	ArchiveHabit(id string) error
	UnarchiveHabit(id string) error
	DeleteHabit(id string) error

	CreateActivity(habitId string, logged time.Time, status string) (string, error)
	GetHabitFromActivity(activityId string) (Habit, error)
	GetActivities(habitId string, after time.Time, before time.Time, limit int) (activities []Activity, hasMore bool, err error)
	DeleteActivity(id string) error

	GetScore(habitId string) (int, error)
}
