/*
This version will use a CSV to store the sessions

It shouldn't be that bad given new sessions will be appended.

Will have to write some cleanup job in future once the file gets too big (a
problem for another time)

We want to go from session id to account.
It should also contain createdAt to calculate expiry

The account information will be stored in a separate file (JSON) with at least
AccountDetails

For simplicity every request we will perform IO (because I don't care and it
can be easily improved)
*/

package auth_file

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"internal/auth"
	"io"
	"os"
	"time"
)

type sessionInfo struct {
	sessionId string
	userId    string
	createdAt time.Time
}

func csvRowToSessionInfo(row []string) (sessionInfo, error) {
	if len(row) < 3 {
		return sessionInfo{}, errors.New("Not enough columns in this CSV")
	}

	createdAt, err := time.Parse(time.RFC3339, row[2])
	if err != nil {
		return sessionInfo{}, err
	}

	return sessionInfo{sessionId: row[0], userId: row[1], createdAt: createdAt}, nil
}

func sessionInfoToCsv(session sessionInfo) [3]string {
	return [3]string{
		session.sessionId,
		session.userId,
		session.createdAt.Format(time.RFC3339),
	}
}

type AuthDatabaseFile struct {
	SessionsFilepath string
	AccountsFilepath string
}

var _ auth.AuthDatabase = (*AuthDatabaseFile)(nil)

// TODO This is so inefficient as we're iterating over a list to find the email
func (a *AuthDatabaseFile) GetUserIdFromEmail(ctx context.Context, email string) (string, error) {
	accountsFile, err := os.Open(a.AccountsFilepath)
	if err != nil {
		return "", err
	}
	defer accountsFile.Close()

	accountsDecoder := json.NewDecoder(accountsFile)

	accounts := []struct {
		Id    string
		Email string
	}{}

	err = accountsDecoder.Decode(&accounts)
	if err != nil {
		return "", err
	}

	for _, account := range accounts {
		if account.Email == email {
			return account.Id, nil
		}
	}

	return "", auth.ErrNotFound
}

func (a *AuthDatabaseFile) AddSession(ctx context.Context, sessionId string, userId string) error {
	sessionsFile, err := os.OpenFile(a.SessionsFilepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer sessionsFile.Close()

	sessionsFileWriter := csv.NewWriter(sessionsFile)
	defer sessionsFileWriter.Flush()

	sessionInfo := sessionInfoToCsv(sessionInfo{sessionId: sessionId, userId: userId, createdAt: time.Now()})

	sessionsFileWriter.Write(sessionInfo[:])

	return nil
}

func (a *AuthDatabaseFile) GetUserIdFromSession(ctx context.Context, sessionId string, since time.Time) (string, error) {
	sessionsFile, err := os.Open(a.SessionsFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", auth.ErrNotFound
		}
		return "", err
	}
	defer sessionsFile.Close()

	sessionsFileReader := csv.NewReader(sessionsFile)
	for {
		sessionCsvRow, err := sessionsFileReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		session, err := csvRowToSessionInfo(sessionCsvRow)
		if err != nil {
			return "", err
		}

		if session.sessionId == sessionId && session.createdAt.After(since) {
			return session.userId, nil
		}
	}

	return "", auth.ErrNotFound
}
