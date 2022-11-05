package habit_share

import (
	"fmt"
)

type App struct {
	Db   HabitsDatabase
	Auth AuthInterface
}

func (a *App) habitOwnerCheck(habit Habit) error {
	// Instead of fetching the habit here let the developer provide the habit
	// themselves Technically I'm the only developer so not a huge deal API is
	// now coupled but it's plainly revealed in the parameters so also not a huge
	// deal Another annoying point is the developer is now faced with two
	// potential points of failure getting the habit fails THEN if the auth check
	// fails Maybe this is fine since it's possible to create the habit
	// themselves (unlikely use case) Often it seems the habit gets used
	// immediately afterwards to there are few scenarios where purely the auth
	// check needs the habit

	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return err
	}

	if habit.Owner != user {
		return PermissionDeniedError
	}

	return nil
}

func (a *App) habitIdOwnerCheck(habitId string) error {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return err
	}

	return a.habitOwnerCheck(habit)
}

func (a *App) habitSharedCheck(habit Habit) error {
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
	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return err
	}
	if err := a.habitOwnerCheck(habit); err != nil {
		return err
	}

	habit.Archived = true
	return a.Db.SetHabit(id, habit)
}

// ChangeFrequency implements HabitsDatabase
func (a *App) ChangeFrequency(id string, newFrequency int) error {

	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return err
	}
	if err := a.habitOwnerCheck(habit); err != nil {
		return err
	}

	if newFrequency < 1 || newFrequency > 7 {
		return &InputError{StringToParse: fmt.Sprint(newFrequency)}
	}
	habit.Frequency = newFrequency
	return a.Db.SetHabit(id, habit)
}

// CreateActivity implements HabitsDatabase
func (a *App) CreateActivity(habitId string, logged Time, status string) (string, error) {
	if err := a.habitIdOwnerCheck(habitId); err != nil {
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

	habit := Habit{Owner: user, Name: name, Frequency: frequency}
	return a.Db.CreateHabit(habit)
}

// DeleteActivity implements HabitsDatabase
func (a *App) DeleteActivity(habitId string, id string) error {
	if err := a.habitIdOwnerCheck(habitId); err != nil {
		return err
	}

	return a.Db.DeleteActivity(habitId, id)
}

// DeleteHabit implements HabitsDatabase
func (a *App) DeleteHabit(id string) error {
	if err := a.habitIdOwnerCheck(id); err != nil {
		return err
	}

	return a.Db.DeleteHabit(id)
}

// GetActivities implements HabitsDatabase
func (a *App) GetActivities(
	habitId string,
	after Time,
	before Time,
	limit int,
) (activities []Activity, hasMore bool, err error) {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return nil, false, err
	}
	if a.habitOwnerCheck(habit) != nil && a.habitSharedCheck(habit) != nil {
		return nil, false, HabitNotFoundError
	}

	return a.Db.GetActivities(habitId, after, before, limit)
}

// GetHabit implements HabitsDatabase
func (a *App) GetHabit(id string) (Habit, error) {
	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return Habit{}, err
	}
	if a.habitOwnerCheck(habit) != nil && a.habitSharedCheck(habit) != nil {
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

	return a.Db.GetMyHabits(user, limit, archived)
}

// GetScore implements HabitsDatabase
func (a *App) GetScore(habitId string) (int, error) {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return 0, err
	}

	if err := a.habitOwnerCheck(habit); err != nil {
		if err := a.habitSharedCheck(habit); err != nil {
			// neither owned nor shared
			return 0, err
		}
		// not owner but shared
	}

	return a.Db.GetScore(habitId)
}

// GetSharedHabits implements HabitsDatabase
func (a *App) GetSharedHabits(limit int) ([]Habit, error) {
	user, err := a.Auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return a.Db.GetSharedHabits(user, limit)
}

// GetSharedWith implements HabitsDatabase
func (a *App) GetSharedWith(habitId string) (map[string]struct{}, error) {
	habit, err := a.Db.GetHabit(habitId)
	if err != nil {
		return nil, err
	}
	if a.habitOwnerCheck(habit) != nil && a.habitSharedCheck(habit) != nil {
		return nil, HabitNotFoundError
	}

	return habit.SharedWith, nil
}

// ChangeName implements HabitsDatabase
func (a *App) ChangeName(id string, newName string) error {
	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return err
	}

	if err := a.habitOwnerCheck(habit); err != nil {
		return err
	}

	// TODO disallow characters like \n for readability
	habit.Name = newName;
	return a.Db.SetHabit(id, habit)
}

// ChangeDescription
func (a *App) ChangeDescription(id string, newDescription string) error {
	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return err
	}

	if err := a.habitOwnerCheck(habit); err != nil {
		return err
	}

	habit.Description = newDescription
	return a.Db.SetHabit(id, habit)
}

// ShareHabit implements HabitsDatabase
func (a *App) ShareHabit(habitId string, friend string) error {
	if err := a.habitIdOwnerCheck(habitId); err != nil {
		return err
	}

	return a.Db.ShareHabit(habitId, friend)
}

// UnShareHabit implements HabitsDatabase
func (a *App) UnShareHabit(habitId string, friend string) error {
	if err := a.habitIdOwnerCheck(habitId); err != nil {
		return err
	}

	return a.Db.UnShareHabit(habitId, friend)
}

// UnarchiveHabit implements HabitsDatabase
func (a *App) UnarchiveHabit(id string) error {
	habit, err := a.Db.GetHabit(id)
	if err != nil {
		return err
	}

	if err := a.habitOwnerCheck(habit); err != nil {
		return err
	}

	habit.Archived = false;
	return a.Db.SetHabit(id, habit)
}
