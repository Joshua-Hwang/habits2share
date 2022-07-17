package habit_share_file

import (
	"os"
	"testing"
)

var inputJson = "input.json"
var outputJson = "output.json"

func TestIo(t *testing.T) {
	if _, err := os.Stat(outputJson); err == nil {
		err = os.Remove(outputJson)
		if err != nil {
			t.Error("during setup failed to delete output json got: ", err)
		}
	}

	t.Run("should write to a file", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		file, err := os.OpenFile(outputJson, os.O_RDWR | os.O_CREATE, 0600)
		if err != nil {
			t.Error("during setup failed to open file got: ", err)
		}
		defer file.Close()

		err = habitShare.WriteToFile(file)
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

	t.Run("should handle missing description from a file", func(t *testing.T) {
		habitShare, err := HabitShareFromFile(inputJson)
		if err != nil {
			t.Error("Failed to parse or read the input file got: ", err)
		}

		habit, ok := habitShare.Habits["testUser2_habitId1"]
		if !ok {
			t.Errorf("Failed to correctly parse the input json %+v", habitShare)
		}
		if habit.Description != "" {
			t.Errorf("Should be empty string got %+v", habit.Description)
		}
		habit.Description = "new description"
	})

	t.Run("should load nothing but continue if file doesn't exist", func(t *testing.T) {
		_, err := HabitShareFromFile("doesn't exist")
		if err != nil {
			t.Error("Failed to handle a non existent file: ", err)
		}
	})

	t.Run("should read from written file", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		// write again to output json to ensure tests are independent
		orgHabitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		file, err := os.OpenFile(outputJson, os.O_RDWR | os.O_CREATE, 0600)
		if err != nil {
			t.Error("during setup failed to open file got: ", err)
		}
		defer file.Close()

		err = orgHabitShare.WriteToFile(file)
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
