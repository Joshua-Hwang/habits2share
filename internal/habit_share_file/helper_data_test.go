package habit_share_file

import (
	"internal/habit_share"
	"time"
)

func generateTestData() (testUsers map[string]User, testHabits map[string]HabitJson) {
	testUsers = map[string]User{
		"testUser1": {
			MyHabits:     map[string]struct{}{"testUser1_habitId1": {}},
			SharedHabits: map[string]struct{}{"testUser2_habitId1": {}},
		},
		"testUser2": {
			MyHabits:     map[string]struct{}{"testUser2_habitId1": {}},
			SharedHabits: map[string]struct{}{},
		},
	}

	testHabits = map[string]HabitJson{
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
				{
					Id:      "testUser2_habitId1_2001-01-01",
					HabitId: "testUser2_habitId1",
					Logged:  time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC),
					Status:  "SUCCESS",
				},
			},
		},
	}

	return
}
