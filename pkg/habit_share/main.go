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

	if status != ActivitySuccess &&
		status != ActivityMinimum &&
		status != ActivityNotDone {
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

// ChangeDescription
func (a *App) ChangeDescription(id string, newDescription string) error {
	if err := a.habitOwnerCheck(id); err != nil {
		return err
	}

	return a.Db.ChangeDescription(id, newDescription)
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
