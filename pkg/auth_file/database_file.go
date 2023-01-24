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
	"io"
	"os"
	"sync"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
)

// TTL in seconds
const cacheTtl = 10

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
	SessionsFilepath    string
	SessionsFileLock    *sync.RWMutex
	sessionCache        map[string]sessionInfo // map session id to user id
	sessionCacheCreated time.Time
	AccountsFilepath    string
}

var _ auth.AuthDatabase = (*AuthDatabaseFile)(nil)

// ExpireSession implements auth.AuthDatabase
func (a *AuthDatabaseFile) ExpireSession(ctx context.Context, sessionId string) error {
	a.SessionsFileLock.Lock()
	defer a.SessionsFileLock.Unlock()

	sessionsFile, err := os.OpenFile(a.SessionsFilepath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer sessionsFile.Close()

	sessionFileReader := csv.NewReader(sessionsFile)
	records, err := sessionFileReader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range(records) {
		if record[0] == sessionId {
			record[2] = time.Time{}.Format(time.RFC3339)
			records[i] = record
		}
	}

	sessionsFileWriter := csv.NewWriter(sessionsFile)
	defer sessionsFileWriter.Flush()

	sessionsFileWriter.WriteAll(records)

	return nil
}

// UserExists implements auth.AuthDatabase
func (a *AuthDatabaseFile) UserExists(ctx context.Context, userId string) (bool, error) {
	accountsFile, err := os.Open(a.AccountsFilepath)
	if err != nil {
		return false, err
	}
	defer accountsFile.Close()

	accountsDecoder := json.NewDecoder(accountsFile)

	accounts := []struct {
		Id    string
		Email string
	}{}

	err = accountsDecoder.Decode(&accounts)
	if err != nil {
		return false, err
	}

	for _, account := range accounts {
		if account.Id == userId {
			return true, nil
		}
	}

	return false, nil
}

// TODO This doesn't scale, we're iterating over a list to find the email
// We're also parsing the whole file
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
	a.SessionsFileLock.Lock()
	defer a.SessionsFileLock.Unlock()

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

// TODO need to perform some clean up of expired sessions for performance reasons. Can be done external in a cronjob.
// NOTE the cache means we can't provide an instant session reset
func (a *AuthDatabaseFile) GetUserIdFromSession(ctx context.Context, sessionId string, since time.Time) (string, error) {
	a.SessionsFileLock.RLock()
	defer a.SessionsFileLock.RUnlock()

	// TODO the reason the locks are at the top here instead of near where we open the file is because
	// this isn't thread safe. If this was done in a constructor this would be fine.
	if a.sessionCache == nil || time.Since(a.sessionCacheCreated) > time.Duration(cacheTtl*float64(time.Second)) {
		a.sessionCache = make(map[string]sessionInfo, 0)
		a.sessionCacheCreated = time.Now()
	}
	// check cache first
	session, ok := a.sessionCache[sessionId]
	if ok && session.sessionId == sessionId && session.createdAt.After(since) {
		return session.userId, nil
	}

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

		a.sessionCache[session.sessionId] = session

		if session.sessionId == sessionId && session.createdAt.After(since) {
			return session.userId, nil
		}
	}

	return "", auth.ErrNotFound
}
