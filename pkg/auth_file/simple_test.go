package auth_file

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/auth"
)

var testCsv_ = "test1.csv"
var testJson_ = "test1.json"
var lock sync.RWMutex

func TestAddSession(t *testing.T) {
	t.Run("should append to file", func(t *testing.T) {
		tempDir := t.TempDir()
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: tempDir + "/test1.csv", SessionsFileLock: &lock}

		err := authDatabaseFile.AddSession(context.Background(), "abcde", "test@user.com")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		data, err := os.ReadFile(tempDir + "/test1.csv")
		if err != nil {
			panic("test failed to read testFile")
		}

		if !strings.HasPrefix(string(data), "abcde,test@user.com,") {
			t.Error("expected row to match expected but got: ", string(data))
		}
	})
}

func TestGetUserIdFromEmail(t *testing.T) {
	tempDir := t.TempDir()
	err := os.WriteFile(tempDir + "/test1.json",
	[]byte("[{\"Id\":\"123\",\"Email\":\"test@user.com\"},{\"Id\":\"987\",\"Email\":\"fake@user.com\"},{\"Id\":\"432\",\"Email\":\"correct@answer.com\"}]"),
		0600,
	)
	if err != nil {
		panic("test failed to write testFile")
	}

	t.Run("should return correct answer", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{AccountsFilepath: tempDir + "/test1.json", SessionsFileLock: &lock}

		userId, err := authDatabaseFile.GetUserIdFromEmail(context.Background(), "correct@answer.com")

		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if userId != "432" {
			t.Error("expected email to be 432 but got: ", userId)
		}
	})
}

func TestGetSession(t *testing.T) {
	tempDir := t.TempDir()
	err := os.WriteFile(tempDir + "/test1.csv",
		[]byte("abcde,userId,2022-05-15T23:31:17+00:00"),
		0600,
	)
	if err != nil {
		panic("test failed to write testFile")
	}

	t.Run("should successfully get session", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: tempDir + "/test1.csv", SessionsFileLock: &lock}

		userId, err := authDatabaseFile.GetUserIdFromSession(context.Background(),
			"abcde",
			time.Date(2022, 05, 15, 00, 00, 00, 00, time.UTC),
		)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if userId != "userId" {
			t.Error("expected email to be userId but got: ", userId)
		}
	})
}

func TestFilesDoNotExist(t *testing.T) {
	t.Run("access session", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: "does not exist", SessionsFileLock: &lock}

		_, err := authDatabaseFile.GetUserIdFromSession(context.Background(),
			"abcde",
			time.Date(2022, 05, 15, 00, 00, 00, 00, time.UTC),
		)
		if err != auth.ErrNotFound {
			t.Error("expected error to be ErrNotFound got: ", err)
		}
	})
}
