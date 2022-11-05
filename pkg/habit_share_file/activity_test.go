package habit_share_file

import (
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"testing"
	"time"
)

func TestActivity(t *testing.T) {
	t.Run("should register activity", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		activityId, err := habitShare.CreateActivity(
			"testUser1_habitId1",
			habit_share.Time{Time: time.Date(1, time.January, 1, 1, 1, 1, 1, time.Local)},
			"SUCCESS",
		)

		if err != nil {
			t.Fatal("CreateActivity returned error unexpectedly:", err)
		}

		if activityId != "testUser1_habitId1_0001-01-01" {
			t.Fatal("CreateActivity did not return correct id got: ", activityId)
		}

		if len(habitShare.Habits["testUser1_habitId1"].Activities) != 1 {
			t.Fatal("Activity was not added to list")
		}
	})
	t.Run("should delete activity", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		err := habitShare.DeleteActivity("testUser2_habitId1", "testUser2_habitId1_2001-01-01")

		if err != nil {
			t.Fatal("DeleteActivity returned error unexpectedly:", err)
		}

		if len(habitShare.Habits["testUser2_habitId1"].Activities) != 0 {
			t.Fatal("Activity was not deleted")
		}
	})

	t.Run("should calcuate score is 0", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		score, err := habitShare.GetScore("testUser1_habitId1")

		if err != nil {
			t.Fatal("GetScore returned error unexpectedly:", err)
		}

		if score != 0 {
			t.Fatal("Score not calculated properly expected 0 got:", score)
		}
	})

	t.Run("should calcuate score is 1", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}
		_, err := habitShare.CreateActivity("testUser2_habitId1", habit_share.Time{Time: time.Now()}, "SUCCESS")

		score, err := habitShare.GetScore("testUser2_habitId1")

		if err != nil {
			t.Fatal("GetScore returned error unexpectedly:", err)
		}

		if score != 1 {
			t.Fatal("Score not calculated properly expected 1 got:", score)
		}
	})

	t.Run("should get activities", func(t *testing.T) {
		testUsers, testHabits := generateTestData()
		habitShare := HabitShareFile{Users: testUsers, Habits: testHabits}

		after := habit_share.Time{Time: testHabits["testUser2_habitId1"].Activities[0].Logged.AddDate(0, 0, -1)}
		before := habit_share.Time{Time: testHabits["testUser2_habitId1"].Activities[0].Logged.AddDate(0, 0, 1)}
		activities, hasMore, err := habitShare.GetActivities("testUser2_habitId1", after, before, 1)

		if len(activities) != 1 {
			t.Fatal("activity was not returned correctly", activities)
		}

		if hasMore {
			t.Fatal("implying more activities when there shouldn't be")
		}

		if err != nil {
			t.Fatal("unexpected error occurred", err)
		}
	})
}
