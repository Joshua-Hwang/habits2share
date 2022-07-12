package habit_share_import

import (
	"encoding/csv"
	"io"
	"strings"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
)

func ImportCsv(db habit_share.HabitsDatabase, auth habit_share.AuthInterface, csvReader *csv.Reader) ([]string, error) {
	userId, err := auth.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	return importCsv(db, userId, csvReader)
}

// We assume the csvReader is without heading
func importCsv(db habit_share.HabitsDatabase, owner string, csvReader *csv.Reader) ([]string, error) {
	nameToHabitId := make(map[string]string)

	// Must read header first
	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(row) < 3 {
			return nil, &habit_share.InputError{StringToParse: strings.Join(row, ",")}
		}

		status := habit_share.ActivityNotDone
		if row[2] == "success" {
			status = habit_share.ActivitySuccess
		}

		// NOTE: intentionally different to habit_share.DateFormat as they are not
		// necessarily the same
		activityDate, err := time.Parse("2006-01-02", row[1])
		if err != nil {
			return nil, err
		}

		// check if name is in nameToHabitId
		if _, ok := nameToHabitId[row[0]]; !ok {
			// NOTE: The csv doesn't provide frequency info
			habitId, err := db.CreateHabit(row[0], owner, 7)
			if err != nil {
				return nil, err
			}
			nameToHabitId[row[0]] = habitId
		}
		habitId := nameToHabitId[row[0]]
		db.CreateActivity(habitId, activityDate, status)
	}

	v := make([]string, 0, len(nameToHabitId))
	for _, value := range nameToHabitId {
		v = append(v, value)
	}

	return v, nil
}
