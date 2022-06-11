package auth_sql

import (
	"context"
	"database/sql"
	"time"
	"internal/auth"
)

type AuthDatabaseSql struct {
	db *sql.DB
}

var _ auth.AuthDatabase = (*AuthDatabaseSql)(nil)

func (a *AuthDatabaseSql) AccountExists(ctx context.Context, email string) (string, error) {
	var orgEmail string
	row := a.db.QueryRowContext(ctx, "SELECT email FROM accounts WHERE email=$1", email)
	if err := row.Scan(&orgEmail); err != nil {
		if err == sql.ErrNoRows {
			return "", auth.ErrNotFound
		} else {
			return "", err
		}
	}

	return orgEmail, nil
}

func (a *AuthDatabaseSql) AddSession(ctx context.Context, sessionId string, email string) error {
	_, err := a.db.ExecContext(ctx, `
INSERT INTO sessions (session, account)
VALUES ($1, $2)`,
		sessionId, email)
	return err
}

func (a *AuthDatabaseSql) GetSession(ctx context.Context, sessionId string, since time.Time) (auth.AccountDetails, error) {
	var accountDetails auth.AccountDetails
	row := a.db.QueryRowContext(ctx, `
SELECT account
FROM sessions
WHERE session=$1 AND created > $2`,
		sessionId, since)

	if err := row.Scan(&accountDetails.Email); err != nil {
		if err == sql.ErrNoRows {
			return accountDetails, auth.ErrNotFound
		} else {
			return accountDetails, err
		}
	}

	return accountDetails, nil
}
