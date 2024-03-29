package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share_file"
)

func (reqDeps RequestDependencies) BuildHabitHandler(habit *habit_share.Habit) http.Handler {
	mux := MuxWrapper{ServeMux: http.NewServeMux()}
	mux.RegisterHandlers("/", map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			// optimised endpoint
			app := reqDeps.HabitApp

			// taken from the activities endpoint
			activities, _, err := app.GetActivities(habit.Id,
				habit_share.Time{Time: time.Now().AddDate(0, 0, -7)},
				habit_share.Time{Time: time.Now().AddDate(0, 0, 1)},
				7,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Something has gone wrong getting activities")
				log.Printf("Something has gone wrong getting activities: %v", err)
				return
			}

			score, err := app.GetScore(habit.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to calculate streak")
				log.Printf("Failed to calculate streak: %v", err)
				return
			}

			response := struct {
				*habit_share.Habit
				Activities []habit_share.Activity
				Score      int
			}{Habit: habit, Activities: activities, Score: score}

			bytes, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Something has gone wrong writing habit to json")
				log.Printf("Something has gone wrong writing habit to json: %v", err)
			}

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", string(bytes))
		},
		"POST": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				fmt.Fprintf(w, "Content Type is not application/json")
				return
			}

			app := reqDeps.HabitApp

			updatePayload := struct {
				Name        string
				Frequency   int
				Description string
			}{}
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&updatePayload)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				var unmarshalErr *json.UnmarshalTypeError
				if errors.As(err, &unmarshalErr) {
					fmt.Fprintf(w, "Bad Request. Wrong Type provided for field: %s", unmarshalErr.Field)
				} else {
					fmt.Fprintf(w, "Bad Request: %s", err)
				}
				return
			}

			if updatePayload.Name != "" {
				err = app.ChangeName(habit.Id, updatePayload.Name)
			}
			if err != nil {
				if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Name was invalid, frequency not modified")
				} else if errors.Is(err, habit_share.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this habit")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed rename habit")
					log.Printf("Something has gone wrong renaming habit: %v", err)
				}
				return
			}

			if updatePayload.Frequency != 0 {
				err = app.ChangeFrequency(habit.Id, updatePayload.Frequency)
			}
			if err != nil {
				if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Frequency was invalid")
				} else if errors.Is(err, habit_share.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this habit")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to change frequency")
					log.Printf("Something has gone wrong changing frequency: %v", err)
				}
				return
			}

			if updatePayload.Description != "" {
				err = app.ChangeDescription(habit.Id, updatePayload.Description)
			}
			if err != nil {
				if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Description was invalid, description not modified")
				} else if errors.Is(err, habit_share.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this habit")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to change description")
					log.Printf("Something has gone wrong changing description: %v", err)
				}
				return
			}

			w.WriteHeader(http.StatusCreated)
		},
		"DELETE": func(w http.ResponseWriter, r *http.Request) {
			app := reqDeps.HabitApp

			permanent := r.URL.Query().Get("permanent")
			if permanent == "true" {
				err := app.DeleteHabit(habit.Id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to delete habit")
					log.Printf("Something has gone wrong deleting activities: %v", err)
				}
			} else {
				err := app.ArchiveHabit(habit.Id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to archive habit")
					log.Printf("Something has gone wrong archiving habit: %v", err)
				}
			}

			w.WriteHeader(http.StatusNoContent)
		},
	})
	mux.RegisterHandlers("/name", map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", habit.Name)
		},
		"POST": func(w http.ResponseWriter, r *http.Request) {
			app := reqDeps.HabitApp

			if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				fmt.Fprintf(w, "Content Type is not text/html")
				return
			}
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to read body of request")
				log.Printf("Something has gone wrong changing name, reading: %v", err)
				return
			}
			// As a security measure take our Id instead of the user input
			err = app.ChangeName(habit.Id, string(b))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to write name to habit")
				log.Printf("Something has gone wrong changing name, writing: %v", err)
				return
			}
			w.WriteHeader(http.StatusCreated)
		},
	})
	mux.RegisterHandlers("/frequency", map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%d", habit.Frequency)
		},
	})
	mux.RegisterHandlers("/score", map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			app := reqDeps.HabitApp

			score, err := app.GetScore(habit.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to calculate streak")
				log.Printf("Failed to calculate streak: %v", err)
				return
			}
			fmt.Fprintf(w, "%d", score)
		},
	})
	// POST to /habit/:habitId/activities with status in body to register an activity
	// GET to /habit/:habitId/activities?limit=...&order=... works on the pagination of activities
	mux.RegisterHandlers("/activities", map[string]http.HandlerFunc{
		"GET": func(w http.ResponseWriter, r *http.Request) {
			app := reqDeps.HabitApp

			beforeString := r.URL.Query().Get("before")
			if beforeString == "" {
				// before is not inclusive but we'd like today to be displayed.
				// A second day has to be added because timezones
				beforeString = time.Now().AddDate(0, 0, 2).Format(habit_share.DateFormat)
			}
			before, err := time.Parse(habit_share.DateFormat, beforeString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Before query is in incorrect, must be in YYYY-mm-dd format")
			}

			afterString := r.URL.Query().Get("after")
			if afterString == "" {
				// A day is removed because timezones
				afterString = time.Now().AddDate(0, 0, -8).Format(habit_share.DateFormat)
			}

			after, err := time.Parse(habit_share.DateFormat, afterString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "After query is in incorrect, must be in YYYY-mm-dd format")
			}

			limitString := r.URL.Query().Get("limit")
			if limitString == "" {
				limitString = "7"
			}
			limit, err := strconv.Atoi(limitString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Limit query is in incorrect, must be an integer")
			}

			if before.Before(after) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, before date is not before after date")
				return
			}

			activities, hasMore, err := app.GetActivities(
				habit.Id,
				habit_share.Time{Time: after},
				habit_share.Time{Time: before},
				limit,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Something has gone wrong getting activities")
				log.Printf("Something has gone wrong getting activities: %v", err)
				return
			}
			// TODO change from RFC3339 to own date format
			response := struct {
				Activities []habit_share.Activity
				HasMore    bool
			}{Activities: activities, HasMore: hasMore}

			bytes, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Something has gone wrong writing habit to json")
				log.Printf("Something has gone wrong writing habit to json: %v", err)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", string(bytes))
		},
		"POST": func(w http.ResponseWriter, r *http.Request) {
			var err error
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				fmt.Fprintf(w, "Content Type is not application/json")
				return
			}

			app := reqDeps.HabitApp

			newActivity := struct {
				Logged string
				Status string
			}{}
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&newActivity)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				var unmarshalErr *json.UnmarshalTypeError
				if errors.As(err, &unmarshalErr) {
					fmt.Fprintf(w, "Bad Request. Wrong Type provided for field: %s", unmarshalErr.Field)
				} else {
					fmt.Fprintf(w, "Bad Request: %s", err)
				}
				return
			}

			parsedLog, err := time.Parse(habit_share.DateFormat, newActivity.Logged)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Bad Request, Logged must be in YYYY-mm-dd format")
				return
			}

			if newActivity.Status == "NOT_DONE" {
				err = app.DeleteActivity(habit.Id, habit_share_file.ConstructActivityId(habit.Id, habit_share.Time{Time: parsedLog}))
			} else {
				_, err = app.CreateActivity(habit.Id, habit_share.Time{Time: parsedLog}, newActivity.Status)
			}

			if err != nil {
				if inputError := (*habit_share.InputError)(nil); errors.As(err, &inputError) {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Bad Request, Status was not one of the defined enum values")
				} else if errors.Is(err, habit_share.PermissionDeniedError) {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprintf(w, "You do not have permissions for this habit")
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Failed to create activity")
					log.Printf("Something has gone wrong creating an activities: %v", err)
				}
				return
			}

			w.WriteHeader(http.StatusCreated)
		},
	})

	return mux
}
