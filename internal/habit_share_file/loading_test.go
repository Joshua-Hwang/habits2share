package habit_share_file

import (
	"internal/habit_share"
	"os"
	"testing"
	"time"
)

var inputJson = "input.json"
var outputJson = "output.json"

var testUsers = map[string]User{
	"testUser1": {
		MyHabits:     map[string]struct{}{"testUser1_habitId1":struct{}{}},
		SharedHabits: map[string]struct{}{"testUser2_habitId1":struct{}{}},
	},
	"testUser2": {
		MyHabits:     map[string]struct{}{"testUser2_habitId1":struct{}{}},
		SharedHabits: map[string]struct{}{},
	},
}

var testHabits = map[string]HabitJson{
	"testUser1_habitId1": {
		Habit: habit_share.Habit{
			Id:        "testUser1_habitId1",
			Owner:     "testUser1",
			Name:      "first habit",
			Frequency: 3,
			Archived:  false,
		},
		Activities: []habit_share.Activity{},
	},
	"testUser2_habitId1": {
		Habit: habit_share.Habit{
			Id:        "testUser2_habitId1",
			Owner:     "testUser2",
			Name:      "my first habit",
			Frequency: 7,
			Archived:  true,
		},
		Activities: []habit_share.Activity{
			{Id: "1/1/1", HabitId: "habitId1", Logged: time.Now(), Status: "SUCCESS"},
		},
	},
}

func TestIo(t *testing.T) {
	if _, err := os.Stat(outputJson); err == nil {
		err = os.Remove(outputJson)
		if err != nil {
			t.Error("during setup failed to delete output json got: ", err)
		}
	}

	t.Run("should write to a file", func(t *testing.T) {
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		err := habitShare.WriteToFile(outputJson)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
	})

	t.Run("should read from a file", func(t *testing.T) {
		habitShare, err := HabitShareFromFile(inputJson)
		if err != nil {
			t.Error("Failed to parse or read the input file got: ", err)
		}

		// TODO probably a nicer way of doing this (for another time)
		if _, ok := habitShare.Users["testUser1"].MyHabits["testUser1_habitId1"]; !ok {
			t.Errorf("Failed to correctly parse the input json %+v", habitShare)
		}
	})

	t.Run("should load nothing but continue if file doesn't exist", func(t *testing.T) {
		_, err := HabitShareFromFile("doesn't exist")
		if err != nil {
			t.Error("Failed to handle a non existent file: ", err)
		}
	})

	t.Run("should read from written file", func(t *testing.T) {
		// write again to output json to ensure tests are independent
		orgHabitShare := HabitShareFile{Users: testUsers, Habits: testHabits}
		err := orgHabitShare.WriteToFile(outputJson)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		habitShare, err := HabitShareFromFile(outputJson)
		if err != nil {
			t.Error("Failed to parse or read the input file got: ", err)
		}

		// TODO probably a nicer way of doing this (for another time)
		if _, ok := habitShare.Users["testUser1"].MyHabits["testUser1_habitId1"]; !ok {
			t.Errorf("Failed to correctly parse the input json %+v", habitShare)
		}
	})
}
