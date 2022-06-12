package habit_share

import (
	"fmt"
	"time"
)

type App struct {
	Db   HabitsDatabase
	Auth AuthInterface
}

func (a *App) habitOwnerCheck(habitId string) error {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return err
	}

	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return err
	}

	if habit.Owner != user {
		return PermissionDeniedError
	}

	return nil
}

func (a *App) habitSharedCheck(habitId string) error {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return err
	}

	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return err
	}

	if _, ok := habit.SharedWith[user]; !ok {
		return PermissionDeniedError
	}

	return nil
}

// ArchiveHabit implements HabitsDatabase
func (a *App) ArchiveHabit(id string) error {
	// some options for permissions checks here
	// 1. auth service has access to the current user
	// 2. provide a hook as a parameter (this becomes annoying when more parameters are needed)
	// 3. no permissions for this at all. Rely on main to perform checks.
	// I'm liking option 1. Though the auth service is recreated every time the
	// connection to the file doesn't have to be.
	// I'd rather not do option 2 as it sounds tedious
	// Option 3 sounds like I'm shifting business logic to main

	// "But user now becomes an important parameter in this function which is obfuscated by the auth."
	// Yes but those that depend on this can read the documentation/code to
	// figure out how to get this function working.
	// additionally the main function will have access to the habits database. it
	// should be noted that might go against the "plugin" architecture if the
	// main function manages to put the database in an inconsistent state. I
	// don't believe that is possible with the current functionality the habits
	// database provides

	// note for future ensure each operation does not result in a incomplete
	// state OR document how to fix the state.

	// another way to look at it is that auth is just the state of the headers
	// this state does not change for a while so exists as a property of App
	// instead of a parameter
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	// yes we're grabbing the habit twice but I don't think there's a way to
	// provide the optimisation where the habit object gets shared around then
	// upserted back to the db
	// in sql you would add the owner condition in the where clause making this
	// operation very fast.
	// I return to the argument of orthogonality. These are different actions
	// that happen to share state. The modification of activities does not help
	// verifying against the habit's owner.

	// IF this (or similar operations) needed optmisation I would delegate this
	// entire function (and auth logic) to the underlying detail documenting that
	// this operation needs to be as fast as possible given auth constraints.
	return a.Db.ArchiveHabit(id)
}

// ChangeFrequency implements HabitsDatabase
func (a *App) ChangeFrequency(id string, newFrequency int) error {
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	if newFrequency < 1 || newFrequency > 7 {
		return &InputError{StringToParse: fmt.Sprint(newFrequency)}
	}
	return a.Db.ChangeFrequency(id, newFrequency)
}

// CreateActivity implements HabitsDatabase
func (a *App) CreateActivity(habitId string, logged time.Time, status string) (string, error) {
	if err := a.habitOwnerCheck(habitId); err != nil {
		return "", err
	}

	if status != activitySuccess &&
		status != activityMinimum &&
		status != activityNotDone {
		return "", &InputError{StringToParse: status}
	}

	return a.Db.CreateActivity(habitId, logged, status)
}

// CreateHabit implements HabitsDatabase
func (a *App) CreateHabit(name string, frequency int) (string, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return "", err
	}

	if frequency < 1 || frequency > 7 {
		return "", &InputError{StringToParse: fmt.Sprint(frequency)}
	}

	return a.Db.CreateHabit(name, user, frequency)
}

// DeleteActivity implements HabitsDatabase
func (a *App) DeleteActivity(id string) error {
	habit, err := a.Db.GetHabitFromActivity(id)
	if err != nil {
		return err
	}

	if err := a.habitOwnerCheck(habit.Id); err != nil {
		return err
	}

	return a.Db.DeleteActivity(id)
}

// DeleteHabit implements HabitsDatabase
func (a *App) DeleteHabit(id string) error {
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	return a.Db.DeleteHabit(id)
}

// GetActivities implements HabitsDatabase
func (a *App) GetActivities(
	habitId string,
	after time.Time,
	before time.Time,
	limit int,
) (activities []Activity, hasMore bool, err error) {
	if a.habitOwnerCheck(habitId) != nil && a.habitSharedCheck(habitId) != nil {
		return nil, false, HabitNotFoundError
	}

	return a.Db.GetActivities(habitId, after, before, limit)
}

// GetHabit implements HabitsDatabase
func (a *App) GetHabit(id string) (Habit, error) {
	if a.habitOwnerCheck(id) != nil && a.habitSharedCheck(id) != nil {
		return Habit{}, HabitNotFoundError
	}

	return a.Db.GetHabit(id)
}

// GetMyHabits implements HabitsDatabase
func (a *App) GetMyHabits(limit int, archived bool) ([]Habit, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return a.Db.GetMyHabits(user, limit, archived);
}

// GetScore implements HabitsDatabase
func (a *App) GetScore(habitId string) (int, error) {
	if err := a.habitOwnerCheck(habitId); err != nil {
		if err := a.habitSharedCheck(habitId); err != nil {
			// neither owned nor shared
			return 0, err
		}
		// not owner but shared
	}

	return a.Db.GetScore(habitId);
}

// GetSharedHabits implements HabitsDatabase
func (a *App) GetSharedHabits(limit int) ([]Habit, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return a.Db.GetSharedHabits(user, limit);
}

// GetSharedWith implements HabitsDatabase
func (a *App) GetSharedWith(habitId string) (map[string]struct{}, error) {
	if a.habitOwnerCheck(habitId) != nil && a.habitSharedCheck(habitId) != nil {
		return nil, HabitNotFoundError
	}

	return a.Db.GetSharedWith(habitId);
}

// ChangeName implements HabitsDatabase
func (a *App) ChangeName(id string, newName string) error {
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	// TODO disallow characters like \n for readability
	return a.Db.ChangeName(id, newName)
}

// ShareHabit implements HabitsDatabase
func (a *App) ShareHabit(habitId string, friend string) error {
	if err := a.habitOwnerCheck(habitId); err != nil {
		return err
	}

	return a.Db.ShareHabit(habitId, friend)
}

// UnShareHabit implements HabitsDatabase
func (a *App) UnShareHabit(habitId string, friend string) error {
	if err := a.habitOwnerCheck(habitId); err != nil {
		return err
	}

	return a.Db.UnShareHabit(habitId, friend)
}

// UnarchiveHabit implements HabitsDatabase
func (a *App) UnarchiveHabit(id string) error {
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	return a.Db.UnarchiveHabit(id)
}
