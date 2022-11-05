package habit_share_import

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
)

// TODO we don't need to mock all of this right?
type HabitsDatabaseMock struct {
	CreateHabitImplementation func(habit habit_share.Habit) (string, error)
	CreateHabitCalled int
	ShareHabitImplementation func(habitId string, friend string) error
	ShareHabitCalled int
	UnShareHabitImplementation func(habitId string, friend string) error
	UnShareHabitCalled int
	GetSharedWithImplementation func(habitId string) (map[string]struct{}, error)
	GetSharedWithCalled int
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetMyHabitsImplementation func(owner string, limit int, archived bool) ([]habit_share.Habit, error)
	GetMyHabitsCalled int
	// this should not show archived habits
	GetSharedHabitsImplementation func(owner string, limit int) ([]habit_share.Habit, error)
	GetSharedHabitsCalled int
	// the value returned should not be modified in case of an in-memory database
	// avoiding copying
	GetHabitImplementation func(id string) (habit_share.Habit, error)
	GetHabitCalled int
	SetHabitImplementation func(id string, habit habit_share.Habit) error
	SetHabitCalled int
	DeleteHabitImplementation func(id string) error
	DeleteHabitCalled int

	CreateActivityImplementation func(habitId string, logged habit_share.Time, status string) (string, error)
	CreateActivityCalled int
	GetHabitFromActivityImplementation func(activityId string) (habit_share.Habit, error)
	GetHabitFromActivityCalled int
	GetActivitiesImplementation func(habitId string, after habit_share.Time, before habit_share.Time, limit int) (activities []habit_share.Activity, hasMore bool, err error)
	GetActivitiesCalled int
	DeleteActivityImplementation func(habitId string, id string) error
	DeleteActivityCalled int

	GetScoreImplementation func(habitId string) (int, error)
	GetScoreCalled int
}

// ArchiveHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) SetHabit(id string, habit habit_share.Habit) error {
	mock.SetHabitCalled++
	return mock.SetHabitImplementation(id, habit)
}

// CreateActivity implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) CreateActivity(habitId string, logged habit_share.Time, status string) (string, error) {
	mock.CreateActivityCalled++
	return mock.CreateActivityImplementation(habitId, logged, status)
}

// CreateHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) CreateHabit(habit habit_share.Habit) (string, error) {
	mock.CreateHabitCalled++
	return mock.CreateHabitImplementation(habit)
}

// DeleteActivity implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) DeleteActivity(habitId string, id string) error {
	mock.DeleteActivityCalled++
	return mock.DeleteActivityImplementation(habitId, id)
}

// DeleteHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) DeleteHabit(id string) error {
	mock.DeleteHabitCalled++
	return mock.DeleteHabitImplementation(id)
}

// GetActivities implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetActivities(habitId string, after habit_share.Time, before habit_share.Time, limit int) (activities []habit_share.Activity, hasMore bool, err error) {
	mock.GetActivitiesCalled++
	return mock.GetActivitiesImplementation(habitId, after, before, limit)
}

// GetHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetHabit(id string) (habit_share.Habit, error) {
	mock.GetHabitCalled++
	return mock.GetHabitImplementation(id)
}

// GetHabitFromActivity implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetHabitFromActivity(activityId string) (habit_share.Habit, error) {
	mock.GetHabitFromActivityCalled++
	return mock.GetHabitFromActivityImplementation(activityId)
}

// GetMyHabits implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetMyHabits(owner string, limit int, archived bool) ([]habit_share.Habit, error) {
	mock.GetMyHabitsCalled++
	return mock.GetMyHabitsImplementation(owner, limit, archived)
}

// GetScore implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetScore(habitId string) (int, error) {
	mock.GetScoreCalled++
	return mock.GetScoreImplementation(habitId)
}

// GetSharedHabits implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetSharedHabits(owner string, limit int) ([]habit_share.Habit, error) {
	mock.GetSharedHabitsCalled++
	return mock.GetSharedHabitsImplementation(owner, limit)
}

// GetSharedWith implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetSharedWith(habitId string) (map[string]struct{}, error) {
	mock.GetSharedWithCalled++
	return mock.GetSharedWithImplementation(habitId)
}

// ShareHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) ShareHabit(habitId string, friend string) error {
	mock.ShareHabitCalled++
	return mock.ShareHabitImplementation(habitId, friend)
}

// UnShareHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) UnShareHabit(habitId string, friend string) error {
	mock.UnShareHabitCalled++
	return mock.UnShareHabitImplementation(habitId, friend)
}

var _ habit_share.HabitsDatabase = (*HabitsDatabaseMock)(nil)

// TODO this test needs to be better
func TestCsvParser(t *testing.T) {
	t.Run("correctly read simple csv", func(t *testing.T) {
		csvReader := csv.NewReader(strings.NewReader(`Habit,Date,Status,Comment
Ankle Rehab,2021-07-20,success,
Ankle Rehab,2021-07-21,success,
Reading,2021-07-22,success,
Reading,2021-07-23,skip,
Ankle Rehab,2021-07-24,success,
Reading,2021-07-24,skip,
Reading,2021-07-25,skip,
Ankle Rehab,2021-07-26,success,
Reading,2021-07-26,skip,`))
		db := &HabitsDatabaseMock{}
		db.CreateHabitImplementation = func(habit habit_share.Habit) (string, error) {
			return habit.Name, nil
		}
		db.CreateActivityImplementation = func(habitId string, logged habit_share.Time, status string) (string, error) {
			return "useless id", nil
		}

		_, err := importCsv(db, "testUser", csvReader)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if db.CreateHabitCalled != 2 {
			t.Error("Expected 2 habits to be created got: ", db.CreateHabitCalled)
		}
		if db.CreateActivityCalled != 9 {
			t.Error("Expected 9 activities to be created got: ", db.CreateActivityCalled)
		}
	})
}
