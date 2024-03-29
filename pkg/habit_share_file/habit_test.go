package habit_share_file

import (
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"testing"
)

func TestHabit(t *testing.T) {
	t.Run("should create a habit and user", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		habit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(habit)
		if err != nil {
			t.Fatal("CreateHabit returned error unexpectedly: ", err)
		}
		if habitId == "" {
			t.Fatal("habit id came back empty")
		}

		user, ok := habitShare.Users["newOwner"]
		if !ok {
			t.Fatal("Could not find newOwner")
		}
		if _, ok := user.MyHabits[habitId]; !ok {
			t.Fatal("New user does not have the habit")
		}
	})

	t.Run("should throw error if habit doesn't exist", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		habit, err := habitShare.GetHabit("not real")

		if err != habit_share.HabitNotFoundError {
			t.Fatal("expecting habit not found but received", habit, err)
		}
	})

	t.Run("should get habit which exists", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		newHabit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		habit, err := habitShare.GetHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if habit.Name != "new habit" {
			t.Fatal("habit name is not what was expected ", habit.Name)
		}
	})

	t.Run("should archive a habit", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		habit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(habit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		err = habitShare.ArchiveHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}
	})

	t.Run("should rename a habit", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		newHabit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		newHabit.Name = "new name"
		err = habitShare.SetHabit(habitId, newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		habit, err := habitShare.GetHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if habit.Name != "new name" {
			t.Fatal("name was not changed")
		}
	})

	t.Run("should change frequency of habit", func(t *testing.T) {
		habitShare := HabitShareFile{Users: map[string]User{}, Habits: map[string]HabitJson{}}

		newHabit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		newHabit.Frequency = 7
		err = habitShare.SetHabit(habitId, newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		habit, err := habitShare.GetHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if habit.Frequency != 7 {
			t.Fatal("frequency was not changed")
		}
	})

	t.Run("should fail to change frequency when out of range", func(t *testing.T) {
		// TODO
	})

	t.Run("should share habit with another user", func(t *testing.T) {
		habitShare := HabitShareFile{
			Users: map[string]User{
				"oldUser": {
					MyHabits:     map[string]struct{}{},
					SharedHabits: map[string]struct{}{},
				},
			},
			Habits: map[string]HabitJson{},
		}

		newHabit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		err = habitShare.ShareHabit(habitId, "oldUser")
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if _, ok := habitShare.Users["oldUser"].SharedHabits[habitId]; !ok {
			t.Fatal("user did not receive shared habit")
		}

		habit, err := habitShare.GetHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if _, ok := habit.SharedWith["oldUser"]; !ok {
			t.Fatal("could not find user in shared with")
		}
	})

	t.Run("should unshare habit with another user", func(t *testing.T) {
		habitShare := HabitShareFile{
			Users: map[string]User{
				"oldUser": {
					MyHabits:     map[string]struct{}{},
					SharedHabits: map[string]struct{}{},
				},
			},
			Habits: map[string]HabitJson{},
		}

		// this is sharing the app
		newHabit := habit_share.Habit{Name: "new habit", Owner: "newOwner", Frequency: 2}
		habitId, err := habitShare.CreateHabit(newHabit)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		err = habitShare.ShareHabit(habitId, "oldUser")
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		habitShare.UnShareHabit(habitId, "oldUser")

		if _, ok := habitShare.Users["oldUser"].SharedHabits[habitId]; ok {
			t.Fatal("user still has access to habit")
		}

		habit, err := habitShare.GetHabit(habitId)
		if err != nil {
			t.Fatal("expected no error got ", err)
		}

		if _, ok := habit.SharedWith["oldUser"]; ok {
			t.Fatal("habit claims it is still being shared")
		}
	})
}
