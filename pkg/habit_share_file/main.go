package habit_share_file

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"

	"github.com/google/uuid"
)

// TTL in seconds
const cacheTtl = 10

const DateFormat = "2006-01-02"

/*
habit_share requires a single id to identity the activity in order for our
system to work the activity id contains the habitId
*/
func ConstructActivityId(habitId string, logged time.Time) string {
	return fmt.Sprintf("%s_%s", habitId, logged.Format(DateFormat))
}

func parseActivityId(activityId string) (habitId string, date time.Time, err error) {
	lastIndex := strings.LastIndex(activityId, "_")
	if lastIndex == -1 {
		err = &habit_share.InputError{StringToParse: activityId}
		return
	}

	habitId = activityId[:lastIndex]
	dateString := activityId[lastIndex+1:]
	date, err = time.Parse(DateFormat, dateString)
	if err != nil {
		err = &habit_share.InputError{StringToParse: activityId}
		return
	}

	return
}

/*
Entire operation is stored is a list of JSONs stored in a single file and
loaded into memory.

After every setting we save it to the file.
*/

type HabitJson struct {
	habit_share.Habit
	Activities []habit_share.Activity
}

type User struct {
	MyHabits     map[string]struct{}
	SharedHabits map[string]struct{}
	// an in memory solution would use pointers but JSONs can't parse pointers
}

type HabitShareFile struct {
	Users    map[string]User
	Habits   map[string]HabitJson
	filename string
	fileLock *sync.Mutex // This can't be a rw mutex as you're always "writing" the parsed file to the struct
	lastRead time.Time
}

var _ habit_share.HabitsDatabase = (*HabitShareFile)(nil)

func HabitShareFromFile(filename string) (*HabitShareFile, error) {
	var habitShareFile HabitShareFile
	habitShareFile.filename = filename
	habitShareFile.fileLock = &sync.Mutex{}

	err := habitShareFile.read()

	if err != nil {
		return nil, err
	}

	return &habitShareFile, nil
}

// it's the responsibility of the server to WriteToFile
// it might be better to use a filename instead of os.File
// oh well
func (a *HabitShareFile) WriteToFile(file *os.File) error {
	jsonString, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonString)
	if err != nil {
		return err
	}

	return nil
}

func (a *HabitShareFile) write() error {
	if a.filename != "" {
		a.fileLock.Lock()
		defer a.fileLock.Unlock()

		file, err := os.OpenFile(a.filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		err = a.WriteToFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *HabitShareFile) read() error {
	if a.filename != "" && time.Since(a.lastRead) > time.Duration(cacheTtl*float64(time.Second)) {
		a.fileLock.Lock()
		defer a.fileLock.Unlock()

		content, err := os.ReadFile(a.filename)
		a.lastRead = time.Now()
		if err != nil || len(content) == 0 {
			if !os.IsNotExist(err) {
				return err
			}
			// file does not exist or got removed
			a.Habits = make(map[string]HabitJson, 0)
			a.Users = make(map[string]User, 0)
			return nil
		}
		err = json.Unmarshal(content, a)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

// ShareHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) ShareHabit(habitId string, friend string) error {
	if err := a.read(); err != nil {
		return err
	}
	// don't forget to write back to the file if it exists
	user, ok := a.Users[friend]
	if !ok {
		return habit_share.UserNotFoundError
	}

	habit, ok := a.Habits[habitId]
	if !ok {
		return habit_share.HabitNotFoundError
	}

	user.SharedHabits[habitId] = struct{}{}
	a.Users[friend] = user

	habit.SharedWith[friend] = struct{}{}

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// GetSharedWith implements habit_share.HabitsDatabase
func (a *HabitShareFile) GetSharedWith(habitId string) (map[string]struct{}, error) {
	if err := a.read(); err != nil {
		return nil, err
	}
	habit, ok := a.Habits[habitId]
	if !ok {
		return nil, habit_share.HabitNotFoundError
	}

	return habit.SharedWith, nil
}

// UnSharehabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) UnShareHabit(habitId string, friend string) error {
	if err := a.read(); err != nil {
		return err
	}
	// don't forget to write back to the file if it exists
	user, ok := a.Users[friend]
	if !ok {
		return habit_share.UserNotFoundError
	}

	habit, ok := a.Habits[habitId]
	if !ok {
		return habit_share.HabitNotFoundError
	}

	delete(user.SharedHabits, habitId)
	a.Users[friend] = user

	delete(habit.SharedWith, friend)

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// ArchiveHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) ArchiveHabit(id string) error {
	if err := a.read(); err != nil {
		return err
	}

	habit, ok := a.Habits[id]
	if !ok {
		return habit_share.HabitNotFoundError
	}
	habit.Archived = true
	a.Habits[id] = habit

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// ChangeFrequency implements habit_share.HabitsDatabase
func (a *HabitShareFile) ChangeFrequency(id string, newFrequency int) error {
	if err := a.read(); err != nil {
		return err
	}

	if habit, ok := a.Habits[id]; !ok {
		return habit_share.HabitNotFoundError
	} else {
		habit.Frequency = newFrequency
		a.Habits[id] = habit
	}

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// ChangeDescription implements habit_share.HabitsDatabase
func (a *HabitShareFile) ChangeDescription(id string, newDescription string) error {
	if err := a.read(); err != nil {
		return err
	}

	if habit, ok := a.Habits[id]; !ok {
		return habit_share.HabitNotFoundError
	} else {
		habit.Description = newDescription
		a.Habits[id] = habit
	}

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// CreateHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) CreateHabit(name string, owner string, frequency int) (string, error) {
	if err := a.read(); err != nil {
		return "", err
	}

	user, ok := a.Users[owner]
	if !ok {
		// if user doesn't exist create user
		user = User{MyHabits: map[string]struct{}{}, SharedHabits: map[string]struct{}{}}
		a.Users[owner] = user
	}

	// Create new habit to ensure we don't modify newHabit parameter
	habit := habit_share.Habit{
		Id:         fmt.Sprintf("%s_%s", owner, uuid.NewString()),
		Owner:      owner, // strings are immutable so we're fine adding this without copying
		Name:       name,
		Frequency:  frequency,
		Archived:   false,
		SharedWith: make(map[string]struct{}, 0),
	}

	// We could probably perform a collision check
	a.Habits[habit.Id] = HabitJson{Habit: habit, Activities: make([]habit_share.Activity, 0)}

	user.MyHabits[habit.Id] = struct{}{}

	err := a.write()
	if err != nil {
		return habit.Id, err
	}

	return habit.Id, nil
}

// CreateActivity implements habit_share.HabitsDatabase
// logged should be the first moments of the day under UTC. If not we transform it anyway.
func (a *HabitShareFile) CreateActivity(habitId string, logged time.Time, status string) (string, error) {
	if err := a.read(); err != nil {
		return "", err
	}
	// activity id will be defined as habit_date
	habit, ok := a.Habits[habitId]
	if !ok {
		return "", habit_share.HabitNotFoundError
	}

	activityId := ConstructActivityId(habitId, logged)
	// check if activity with that id already exists
	// TODO this doesn't scale
	toAppend := true
	for i, activity := range habit.Activities {
		if activity.Id == activityId {
			if activity.Status == status {
				return activityId, nil
			}
			// update status
			habit.Activities[i].Status = status
			toAppend = false
		}
	}
	if toAppend {
		appended := append(habit.Activities, habit_share.Activity{
			Id:      activityId,
			HabitId: habitId,
			Logged:  time.Date(logged.Year(), logged.Month(), logged.Day(), 0, 0, 0, 0, time.UTC),
			Status:  status,
		})
		// sort is ascending so later times are further down the array
		sort.Slice(appended[:], func(i, j int) bool {
			return appended[i].Logged.Before(appended[j].Logged)
		})
		habit.Activities = appended
	}

	a.Habits[habitId] = habit

	err := a.write()
	if err != nil {
		return activityId, err
	}

	return activityId, nil
}

func (a *HabitShareFile) GetHabitFromActivity(activityId string) (habit_share.Habit, error) {
	if err := a.read(); err != nil {
		return habit_share.Habit{}, err
	}

	habitId, _, err := parseActivityId(activityId)
	if err != nil {
		return habit_share.Habit{}, err
	}

	habit, ok := a.Habits[habitId]
	if !ok {
		return habit_share.Habit{}, habit_share.HabitNotFoundError
	}

	return habit.Habit, nil
}

// DeleteActivity implements habit_share.HabitsDatabase
func (a *HabitShareFile) DeleteActivity(id string) error {
	// I think this is the most efficient as I expect usage to be near the most
	// recent (end of array)
	habitId, date, err := parseActivityId(id)
	if err != nil {
		return err
	}

	habit, ok := a.Habits[habitId]
	if !ok {
		return habit_share.HabitNotFoundError
	}

	n := len(habit.Activities)

	// activities should always be sorted
	index := sort.Search(n, func(i int) bool {
		return habit.Activities[i].Logged.After(date) || habit.Activities[i].Logged.Equal(date)
	})

	if index == n || habit.Activities[index].Id != id {
		return habit_share.ActivityNotFoundError
	}

	// shift everything down
	for i := index + 1; i < n; i++ {
		habit.Activities[i-1] = habit.Activities[i]
	}
	habit.Activities = habit.Activities[:n-1]

	a.Habits[habitId] = habit

	err = a.write()
	if err != nil {
		return err
	}

	return nil
}

// DeleteHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) DeleteHabit(id string) error {
	habit, ok := a.Habits[id]
	if !ok {
		return habit_share.HabitNotFoundError
	}

	for userId := range habit.SharedWith {
		user, ok := a.Users[userId]
		if !ok {
			panic("Habit shared with user that doesn't exist")
		}
		delete(user.SharedHabits, id)
	}

	ownerId := habit.Owner
	owner, ok := a.Users[ownerId]
	if !ok {
		panic("Habit exists but owner does not")
	}
	delete(owner.MyHabits, id)

	delete(a.Habits, id)

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// GetActivities implements habit_share.HabitsDatabase
func (a *HabitShareFile) GetActivities(
	habitId string,
	after time.Time,
	before time.Time,
	limit int,
) (activities []habit_share.Activity, hasMore bool, err error) {
	if err := a.read(); err != nil {
		return nil, false, err
	}

	if before.Before(after) {
		return nil, false, &habit_share.InputError{
			StringToParse: fmt.Sprintf("before=%s after=%s", before, after),
		}
	}
	habit, ok := a.Habits[habitId]
	if !ok {
		return nil, false, habit_share.HabitNotFoundError
	}

	n := len(habit.Activities)

	l := sort.Search(n, func(i int) bool {
		return habit.Activities[i].Logged.After(after) ||
			habit.Activities[i].Logged.Equal(after)
	})
	r := sort.Search(n, func(i int) bool {
		// doesn't seem right to use "After" but trust me bro :)
		return habit.Activities[i].Logged.After(before) ||
			habit.Activities[i].Logged.Equal(before)
	})

	hasMore = false
	if l+limit < r {
		hasMore = true
		r = l + limit
	}

	activities = habit.Activities[l:r]

	return activities, hasMore, nil
}

// GetHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) GetHabit(id string) (habit_share.Habit, error) {
	if err := a.read(); err != nil {
		return habit_share.Habit{}, err
	}
	habit, ok := a.Habits[id]
	if !ok {
		return habit_share.Habit{}, habit_share.HabitNotFoundError
	}
	return habit.Habit, nil
}

// GetMyHabits implements habit_share.HabitsDatabase
func (a *HabitShareFile) GetMyHabits(owner string, limit int, archived bool) ([]habit_share.Habit, error) {
	if err := a.read(); err != nil {
		return nil, err
	}
	// assumes the user is valid thus it must not be in the database yet
	user, ok := a.Users[owner]
	if !ok {
		return make([]habit_share.Habit, 0), nil
	}

	myHabits := make([]habit_share.Habit, 0, len(user.MyHabits))
	for habitId := range user.MyHabits {
		habit, ok := a.Habits[habitId]
		if !ok {
			// this is really bad
			return nil, habit_share.HabitNotFoundError
		}

		// either the habit isnt archived or we're looking for archived habits
		if !habit.Archived || archived {
			myHabits = append(myHabits, habit.Habit)
		}
	}

	// map does not guarantee this is in order
	sort.Slice(myHabits[:], func(i, j int) bool {
		return myHabits[i].Name < myHabits[j].Name
	})

	return myHabits, nil
}

// GetSharedHabits implements habit_share.HabitsDatabase
func (a *HabitShareFile) GetSharedHabits(owner string, limit int) ([]habit_share.Habit, error) {
	if err := a.read(); err != nil {
		return nil, err
	}
	// assumes the user is valid thus it must not be in the database yet
	user, ok := a.Users[owner]
	if !ok {
		return make([]habit_share.Habit, 0), nil
	}

	sharedHabits := make([]habit_share.Habit, 0, len(user.SharedHabits))
	for habitId := range user.SharedHabits {
		habit, ok := a.Habits[habitId]
		if !ok {
			panic("Habit that is shared does not exist")
		}
		if !habit.Archived {
			// we could remove archived habits from shared habits during archival but
			// cleanup would be harder to reason about as hanging path are everywhere
			sharedHabits = append(sharedHabits, habit.Habit)
		}
	}

	// map does not guarantee this is in order
	sort.Slice(sharedHabits[:], func(i, j int) bool {
		return sharedHabits[i].Name < sharedHabits[j].Name
	})

	return sharedHabits, nil
}

// GetStreak implements habit_share.HabitsDatabase
// TODO for performance this should be calculated on activity entry and cached
func (a *HabitShareFile) GetScore(habitId string) (int, error) {
	if err := a.read(); err != nil {
		return 0, err
	}
	habit, ok := a.Habits[habitId]
	if !ok {
		return 0, habit_share.HabitNotFoundError
	}

	// check frequency
	// streak_counter for streak
	// ignore current week and find last sunday (last sunday if sunday is today)
	// start counting from that day till monday
	// if count >= frequency increment streak_counter
	// if count < frequency stop and return
	// if i == 0 return
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	totalScore := 0  // count for successes
	weeklyCount := 0 // threshold for frequency (counts minimum and success)
	// values outside the normal range are normalised day -1 goes to the previous month
	// weekStart is at the END of Sunday. The first second of Monday hence the -1 == +6 mod 7
	// start of this week
	weekStart := today.AddDate(0, 0, -((int(today.Weekday()) + 6) % 7))

	index := len(habit.Activities) - 1
	// loop for current week doesn't matter what the score is this week assume it's part of the streak
	for ; index >= 0; index-- {
		if habit.Activities[index].Status == "NOT_DONE" {
			continue
		}

		if habit.Activities[index].Logged.Before(weekStart) {
			break
		}
		// assume this week is part of a streak
		// find first activity that is before or on the startDate
		// This is probably faster than binary search given distribution of requests
		if habit.Activities[index].Status == "SUCCESS" {
			totalScore++
		}
	}

	weekStart = weekStart.AddDate(0, 0, -7)
	for ; index >= 0; index-- {
		// TODO don't store NOT_DONE just delete them
		if habit.Activities[index].Status == "NOT_DONE" {
			continue
		}

		if habit.Activities[index].Logged.Before(weekStart) {
			// the final week (even if incomplete) is considered part of the streak
			if weeklyCount < habit.Frequency {
				return totalScore, nil
			}
			weeklyCount = 0
			weekStart = weekStart.AddDate(0, 0, -7)
		}

		weeklyCount++
		if habit.Activities[index].Status == "SUCCESS" {
			totalScore++
		}
	}

	return totalScore, nil
}

// ChangeName implements habit_share.HabitsDatabase
func (a *HabitShareFile) ChangeName(id string, newName string) error {
	habit, ok := a.Habits[id]
	if !ok {
		return habit_share.HabitNotFoundError
	}

	habit.Name = newName
	a.Habits[id] = habit

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}

// UnarchiveHabit implements habit_share.HabitsDatabase
func (a *HabitShareFile) UnarchiveHabit(id string) error {
	if habit, ok := a.Habits[id]; !ok {
		return habit_share.HabitNotFoundError
	} else {
		habit.Archived = false
		a.Habits[id] = habit
	}

	err := a.write()
	if err != nil {
		return err
	}

	return nil
}
