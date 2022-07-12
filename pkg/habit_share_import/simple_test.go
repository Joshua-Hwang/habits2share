package habit_share_import

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
)

// TODO we don't need to mock all of this right?
type HabitsDatabaseMock struct {
	CreateHabitImplementation func(name string, owner string, frequency int) (string, error)
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
	ChangeNameImplementation func(id string, newName string) error
	ChangeNameCalled int
	// frequency should technically be checked (1-7) in this part prior to sending
	// request to underlying implementation
	ChangeFrequencyImplementation func(id string, newFrequency int) error
	ChangeFrequencyCalled int
	ArchiveHabitImplementation func(id string) error
	ArchiveHabitCalled int
	UnarchiveHabitImplementation func(id string) error
	UnarchiveHabitCalled int
	DeleteHabitImplementation func(id string) error
	DeleteHabitCalled int

	CreateActivityImplementation func(habitId string, logged time.Time, status string) (string, error)
	CreateActivityCalled int
	GetHabitFromActivityImplementation func(activityId string) (habit_share.Habit, error)
	GetHabitFromActivityCalled int
	GetActivitiesImplementation func(habitId string, after time.Time, before time.Time, limit int) (activities []habit_share.Activity, hasMore bool, err error)
	GetActivitiesCalled int
	DeleteActivityImplementation func(id string) error
	DeleteActivityCalled int

	GetScoreImplementation func(habitId string) (int, error)
	GetScoreCalled int
}

// ArchiveHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) ArchiveHabit(id string) error {
	mock.ArchiveHabitCalled++
	return mock.ArchiveHabitImplementation(id)
}

// ChangeFrequency implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) ChangeFrequency(id string, newFrequency int) error {
	mock.ChangeFrequencyCalled++
	return mock.ChangeFrequencyImplementation(id, newFrequency)
}

// ChangeName implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) ChangeName(id string, newName string) error {
	mock.ChangeNameCalled++
	return mock.ChangeNameImplementation(id, newName)
}

// CreateActivity implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) CreateActivity(habitId string, logged time.Time, status string) (string, error) {
	mock.CreateActivityCalled++
	return mock.CreateActivityImplementation(habitId, logged, status)
}

// CreateHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) CreateHabit(name string, owner string, frequency int) (string, error) {
	mock.CreateHabitCalled++
	return mock.CreateHabitImplementation(name, owner, frequency)
}

// DeleteActivity implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) DeleteActivity(id string) error {
	mock.DeleteActivityCalled++
	return mock.DeleteActivityImplementation(id)
}

// DeleteHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) DeleteHabit(id string) error {
	mock.DeleteHabitCalled++
	return mock.DeleteHabitImplementation(id)
}

// GetActivities implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) GetActivities(habitId string, after time.Time, before time.Time, limit int) (activities []habit_share.Activity, hasMore bool, err error) {
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

// UnarchiveHabit implements habit_share.HabitsDatabase
func (mock *HabitsDatabaseMock) UnarchiveHabit(id string) error {
	mock.UnarchiveHabitCalled++
	return mock.UnarchiveHabitImplementation(id)
}

var _ habit_share.HabitsDatabase = (*HabitsDatabaseMock)(nil)

// TODO this test needs to be better
func TestCsvParser(t *testing.T) {
	t.Run("correctly read simple csv", func(t *testing.T) {
		csvReader := csv.NewReader(strings.NewReader(`Ankle Rehab,2021-07-20,success,
Ankle Rehab,2021-07-21,success,
Reading,2021-07-22,success,
Reading,2021-07-23,skip,
Ankle Rehab,2021-07-24,success,
Reading,2021-07-24,skip,
Reading,2021-07-25,skip,
Ankle Rehab,2021-07-26,success,
Reading,2021-07-26,skip,`))
		db := &HabitsDatabaseMock{}
		db.CreateHabitImplementation = func(name string, owner string, frequency int) (string, error) {
			return name, nil
		}
		db.CreateActivityImplementation = func(habitId string, logged time.Time, status string) (string, error) {
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
			t.Error("Expected 5 activities to be created got: ", db.CreateActivityCalled)
		}
	})
}
