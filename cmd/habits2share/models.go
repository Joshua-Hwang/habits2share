/*
This file contains all the queries we make.
Better for all of them to be located in a central file.
In future I might either use an ORM or develop something primitive
*/
package main

import (
	"context"
	"database/sql"
	"internal/data"
	"time"

	"github.com/google/uuid"
)

type AccountRow struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SharedHabitRow struct {
	// this will likely bite me in the ass one day
	// not sure whether to expand this out or not
	Habit uuid.UUID `json:"habit"`
	Email string    `json:"email"`
}

// These reflect the database
type HabitRow struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Owner         string    `json:"owner"`
	CurrentStreak int       `json:"current_streak"`
	Frequency     int       `json:"frequency"`
}

type ActivityRow struct {
	habit        *HabitRow `json:"-"`
	Logged       time.Time `json:"logged"`
	LoggedStatus string    `json:"logged_status"`
}

func getAccounts(ctx context.Context) ([]AccountRow, error) {
	db := data.GetDb()
	accounts := make([]AccountRow, 0)

	rows, err := db.QueryContext(ctx,
		`SELECT email, name FROM accounts`)

	if err != nil {
		return accounts, err
	}
	defer rows.Close()

	for rows.Next() {
		account := AccountRow{}
		if err := rows.Scan(&account.Email, &account.Name); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// TODO archived versions

func getHabit(ctx context.Context, id uuid.UUID) (HabitRow, error) {
	db := data.GetDb()
	habit := HabitRow{}

	row := db.QueryRowContext(ctx,
		`SELECT id, habits.name as habit_name, description, accounts.name as account_name, current_streak, frequency
FROM habits
INNER JOIN accounts ON habits.owner=accounts.email
WHERE id=$1`,
		id)
	if err := row.Scan(&habit.Id, &habit.Name, &habit.Description, &habit.Owner, &habit.CurrentStreak, &habit.Frequency); err != nil {
		return habit, err
	}

	return habit, nil
}

func getHabitOwnedByUser(ctx context.Context, id uuid.UUID, userId string) (HabitRow, error) {
	db := data.GetDb()
	habit := HabitRow{}

	row := db.QueryRowContext(ctx,
		`SELECT id, habits.name as habit_name, description, accounts.name as account_name, current_streak, frequency
FROM habits
INNER JOIN accounts ON habits.owner=accounts.email
WHERE id=$1 AND owner=$2`,
		id, userId)
	if err := row.Scan(&habit.Id, &habit.Name, &habit.Description, &habit.Owner, &habit.CurrentStreak, &habit.Frequency); err != nil {
		return habit, err
	}

	return habit, nil
}

func getHabitForUser(ctx context.Context, id uuid.UUID, userId string) (HabitRow, error) {
	db := data.GetDb()
	habit := HabitRow{}

	row := db.QueryRowContext(ctx,
		`SELECT id, habits.name as habit_name, description, accounts.name as account_name, current_streak, frequency
FROM habits
INNER JOIN accounts ON habits.owner=accounts.email
WHERE id=$1
	AND (owner=$2 OR id IN (SELECT habit FROM shared_habits WHERE shared_with=$2))`,
		id, userId)
	if err := row.Scan(&habit.Id, &habit.Name, &habit.Description, &habit.Owner, &habit.CurrentStreak, &habit.Frequency); err != nil {
		return habit, err
	}

	return habit, nil
}

func getHabitsForUser(ctx context.Context, userId string) ([]HabitRow, error) {
	db := data.GetDb()
	habits := make([]HabitRow, 0)

	rows, err := db.QueryContext(ctx,
		`SELECT id, habits.name as habit_name, description, accounts.name as account_name, current_streak, frequency
FROM habits
INNER JOIN accounts ON habits.owner=accounts.email
WHERE owner=$1 AND archived=FALSE
ORDER BY habit_name`,
		userId)
	if err != nil {
		return habits, err
	}
	defer rows.Close()

	for rows.Next() {
		habit := HabitRow{}
		if err := rows.Scan(&habit.Id, &habit.Name, &habit.Description, &habit.Owner, &habit.CurrentStreak, &habit.Frequency); err != nil {
			return habits, err
		}
		habits = append(habits, habit)
	}

	return habits, nil
}

func getHabitsShared(ctx context.Context, userId string) ([]HabitRow, error) {
	db := data.GetDb()
	habits := make([]HabitRow, 0)

	rows, err := db.QueryContext(ctx,
		`SELECT id, habits.name as habit_name, description, accounts.name as account_name, current_streak, frequency
FROM habits
INNER JOIN accounts ON habits.owner=accounts.email
WHERE id IN (SELECT habit FROM shared_habits WHERE shared_with=$1) AND archived=FALSE
ORDER BY habit_name`,
		userId)
	if err != nil {
		return habits, err
	}
	defer rows.Close()

	for rows.Next() {
		habit := HabitRow{}
		if err := rows.Scan(&habit.Id, &habit.Name, &habit.Description, &habit.Owner, &habit.CurrentStreak, &habit.Frequency); err != nil {
			return habits, err
		}
		habits = append(habits, habit)
	}
	if err := rows.Err(); err != nil {
		return habits, err
	}

	return habits, nil
}

func getHabitsSharedByUser(ctx context.Context, userId string) ([]SharedHabitRow, error) {
	db := data.GetDb()
	sharedHabits := make([]SharedHabitRow, 0)

	rows, err := db.QueryContext(ctx,
		`SELECT habit, shared_with
FROM shared_habits INNER JOIN habits
ON shared_habits.habit = habits.id
WHERE habits.owner=$1`,
		userId)
	if err != nil {
		return sharedHabits, err
	}
	defer rows.Close()

	for rows.Next() {
		sharedHabit := SharedHabitRow{}
		if err := rows.Scan(&sharedHabit.Habit, &sharedHabit.Email); err != nil {
			return sharedHabits, err
		}
		sharedHabits = append(sharedHabits, sharedHabit)
	}

	return sharedHabits, nil
}

func getActivities(ctx context.Context, habit HabitRow, before time.Time, limit int, offset int) ([]ActivityRow, error) {
	db := data.GetDb()
	activities := make([]ActivityRow, 0)
	rows, err := db.QueryContext(ctx,
		`SELECT logged, logged_status
FROM activity
WHERE habit=$1 AND logged <= $2
ORDER BY logged DESC
LIMIT $3 OFFSET $4`,
		habit.Id, before, limit, offset)
	if err != nil {
		return activities, err
	}
	defer rows.Close()

	for rows.Next() {
		activity := ActivityRow{habit: &habit}
		if err := rows.Scan(&activity.Logged, &activity.LoggedStatus); err != nil {
			return activities, err
		}
		activities = append(activities, activity)
	}
	if err := rows.Err(); err != nil {
		return activities, err
	}

	return activities, nil
}

// All updates and insertions are done through transactions

func renameHabitForUser(ctx context.Context, txn *sql.Tx, newName string, id uuid.UUID, owner string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`UPDATE habits
SET name=$1
WHERE id=$2 AND owner=$3`,
		newName, id, owner)
}

func describeHabitForUser(ctx context.Context, txn *sql.Tx, newDescription string, id uuid.UUID, owner string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`UPDATE habits
SET description=$1
WHERE id=$2 AND owner=$3`,
		newDescription, id, owner)
}

func changeFrequencyHabitForUser(ctx context.Context, txn *sql.Tx, newFreq int, id uuid.UUID, owner string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`UPDATE habits
SET frequency=$1
WHERE id=$2 AND owner=$3`,
		newFreq, id, owner)
}

func addHabit(ctx context.Context, txn *sql.Tx, name string, description string, frequency int, owner string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`INSERT INTO habits (name, description, owner, frequency)
VALUES ($1, $2, $3, $4)`,
		name, description, owner, frequency)
}

func shareHabit(ctx context.Context, txn *sql.Tx, habit uuid.UUID, sharedWith string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`INSERT INTO shared_habits (habit, shared_with)
VALUES ($1, $2)`,
		habit, sharedWith)
}

func archiveHabit(ctx context.Context, txn *sql.Tx, habit uuid.UUID) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`UPDATE habits
SET archived=TRUE
WHERE id=$1`,
		habit)
}

func deleteHabit(ctx context.Context, txn *sql.Tx, habit uuid.UUID) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`DELETE
FROM habits
WHERE id=$1`,
		habit)
}

func deleteShare(ctx context.Context, txn *sql.Tx, habit uuid.UUID, sharedWith string) (sql.Result, error) {
	return txn.ExecContext(ctx,
		`DELETE
FROM shared_habits
WHERE habit=$1 and shared_with=$2`,
		habit, sharedWith)
}

func addActivity(ctx context.Context, txn *sql.Tx, habit uuid.UUID, date time.Time, status string) (sql.Result, error) {
	return txn.Exec(
		`INSERT INTO activity (habit, logged, logged_status)
VALUES ($1, $2, $3)
ON CONFLICT(habit, logged) DO UPDATE
SET logged_status=EXCLUDED.logged_status`,
		habit, date, status)
}

func deleteActivity(ctx context.Context, txn *sql.Tx, habit uuid.UUID, date time.Time) (sql.Result, error) {
	return txn.Exec(
		`DELETE
FROM activity
WHERE habit=$1
	AND "logged"=$2`,
		habit, date)
}

func updateStreak(ctx context.Context, txn *sql.Tx, habitId uuid.UUID) (sql.Result, error) {
	// 2000-01-03 is a Monday
	return txn.ExecContext(ctx,
		`WITH
	weeklyActivity AS (
		SELECT
			COUNT(*) activityCount,
			MIN(logged) weekStart,
			(DATE_PART('day', "logged" - TO_TIMESTAMP('2000-01-03', 'YYYY-MM-DD'))::INTEGER)/7 week
		FROM activity
		WHERE habit=$1
		GROUP BY DATE_PART('day', "logged" - TO_TIMESTAMP('2000-01-03', 'YYYY-MM-DD'))::INTEGER/7
		HAVING COUNT(*) >= (SELECT frequency FROM habits WHERE id=$1)
			OR DATE_PART('day', now() - TO_TIMESTAMP('2000-01-03', 'YYYY-MM-DD'))::INTEGER/7 = DATE_PART('day', "logged" - TO_TIMESTAMP('2000-01-03', 'YYYY-MM-DD'))::INTEGER/7
	),
	streaks AS (
		SELECT activityCount, weekStart, week,
			week - ROW_NUMBER() OVER (ORDER BY week) streakNumber
		FROM weeklyActivity
	),
	latestWeek AS (
		SELECT streakNumber
		FROM streaks
		ORDER BY weekStart DESC
		LIMIT 1
	),
	streak AS (
		SELECT MIN(weekStart) startDate
		FROM streaks
		WHERE streakNumber = (SELECT streakNumber FROM latestWeek)
		GROUP BY streakNumber
		ORDER BY startDate DESC
		LIMIT 1
	)
UPDATE habits
SET current_streak=subquery.c
FROM (
	SELECT count(*) AS c
	FROM activity
	WHERE habit=$1
		AND ("logged" >= (SELECT startDate FROM streak) OR "logged" >= DATE_TRUNC('week', now()))
		AND logged_status='SUCCESS'
) AS subquery
WHERE habits.id=$1`,
		habitId)
}
