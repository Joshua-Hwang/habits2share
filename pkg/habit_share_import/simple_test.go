package habit_share_import

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	"github.com/Joshua-Hwang/habits2share/pkg/habit_share/mock"
	"github.com/golang/mock/gomock"
)

// TODO this test needs to be better
func TestCsvParser(t *testing.T) {
	t.Run("correctly read simple csv", func(t *testing.T) {
		csvReader := csv.NewReader(strings.NewReader(`Habit,Date,Status,Comment
Ankle Rehab,2021-07-20,success,
Ankle Rehab,2021-07-21,success,
Reading,2021-07-22,success,
Reading,2021-07-23,skip,
Ankle Rehab,2021-07-24,success,
Reading,2021-07-24,skip,
Reading,2021-07-25,skip,
Ankle Rehab,2021-07-26,success,
Reading,2021-07-26,skip,`))
		ctrl := gomock.NewController(t)
		db := mock_habit_share.NewMockHabitsDatabase(ctrl)
		db.EXPECT().CreateHabit(gomock.Any()).Times(2).DoAndReturn(
			func(habit habit_share.Habit) (string, error) {
				return habit.Name, nil
			})
		db.EXPECT().CreateActivity(gomock.Any(), gomock.Any(), gomock.Any()).Times(9).Return("useless id", nil)

		_, err := importCsv(db, "testUser", csvReader)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
	})
}
