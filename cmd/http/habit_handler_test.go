package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Joshua-Hwang/habits2share/cmd/http/mock"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/golang/mock/gomock"
)

func TestHandleHabitEndpoints(t *testing.T) {
	t.Run("GET / returns habit info", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		habitApp := mock_main.NewMockHabitAppInterface(ctrl)
		reqDeps := RequestDependencies{HabitApp: habitApp}
		habit := habit_share.Habit{
			Id:          "mock id",
			Owner:       "mock owner",
			SharedWith:  map[string]struct{}{},
			Name:        "mock name",
			Description: "mock desc",
			Frequency:   4,
			Archived:    false,
		}
		habitHandler := reqDeps.BuildHabitHandler(&habit)

		habitApp.EXPECT().GetActivities("mock id", gomock.Any(), gomock.Any(), 7).
			Return([]habit_share.Activity{}, false, nil)
		habitApp.EXPECT().GetScore("mock id").Return(20, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		habitHandler.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()

		resPayload := struct {
			*habit_share.Habit
		}{}
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resPayload)
		if err != nil {
			t.Error("expected err to be nil got:", err)
		}
		if resPayload.Id != habit.Id || resPayload.Owner != habit.Owner || resPayload.Name != habit.Name {
			t.Error("expected equality between", habit, "and", *resPayload.Habit)
		}
	})

	t.Run("GET / returns activity info", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		habitApp := mock_main.NewMockHabitAppInterface(ctrl)
		reqDeps := RequestDependencies{HabitApp: habitApp}
		habit := habit_share.Habit{
			Id:          "mock id",
			Owner:       "mock owner",
			SharedWith:  map[string]struct{}{},
			Name:        "mock name",
			Description: "mock desc",
			Frequency:   4,
			Archived:    false,
		}
		habitHandler := reqDeps.BuildHabitHandler(&habit)

		activities := []habit_share.Activity{
			{Id: "fake id 1",
				Logged: habit_share.Time{Time: time.Now().AddDate(0, 0, -2)}},
			{Id: "fake id 2",
				Logged: habit_share.Time{Time: time.Now().AddDate(0, 0, 0)}},
			{Id: "fake id 3",
				Logged: habit_share.Time{Time: time.Now().AddDate(0, 0, -1)}},
		}
		habitApp.EXPECT().GetActivities("mock id", gomock.Any(), gomock.Any(), 7).
			Return(activities, false, nil)
		habitApp.EXPECT().GetScore("mock id").Return(20, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		habitHandler.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()

		resPayload := struct {
			Activities []habit_share.Activity
		}{}
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resPayload)
		if err != nil {
			t.Error("expected err to be nil got:", err)
		}
		if len(resPayload.Activities) != len(activities) ||
			resPayload.Activities[0].Id != activities[0].Id {
			t.Error("expected equality between", activities, "and", resPayload.Activities)
		}
	})

	t.Run("POST / updates habit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		habitApp := mock_main.NewMockHabitAppInterface(ctrl)
		reqDeps := RequestDependencies{HabitApp: habitApp}
		habit := habit_share.Habit{
			Id:          "mock id",
			Owner:       "mock owner",
			SharedWith:  map[string]struct{}{},
			Name:        "mock name",
			Description: "mock desc",
			Frequency:   4,
			Archived:    false,
		}
		habitHandler := reqDeps.BuildHabitHandler(&habit)

		habitApp.EXPECT().ChangeName("mock id", "new name")

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{\"Name\": \"new name\"}"))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		habitHandler.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Error("expected status code to be", http.StatusCreated, "got", res.StatusCode)
		}
	})
}
