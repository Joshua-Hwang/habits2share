package auth_file

import (
	"context"
	"internal/auth"
	"os"
	"strings"
	"testing"
	"time"
)

var testCsv = "test1.csv"
var testJson = "test1.json"

func TestAddSession(t *testing.T) {
	if _, err := os.Stat(testCsv); err == nil {
		err = os.Remove(testCsv)
		if err != nil {
			t.Error("during setup failed to delete test1.csv got: ", err)
		}
	}

	t.Run("should append to file", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: testCsv}

		err := authDatabaseFile.AddSession(context.Background(), "abcde", "test@user.com")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		data, err := os.ReadFile(testCsv)
		if err != nil {
			panic("test failed to read testFile")
		}

		if !strings.HasPrefix(string(data), "abcde,test@user.com,") {
			t.Error("expected row to match expected but got: ", string(data))
		}
	})
}

func TestAccountExists(t *testing.T) {
	err := os.WriteFile(testJson,
		[]byte("[{\"email\":\"test@user.com\"},{\"email\":\"fake@user.com\"},{\"email\":\"correct@answer.com\"}]"),
		0600,
	)
	if err != nil {
		panic("test failed to write testFile")
	}

	t.Run("should return correct answer", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{AccountsFilepath: testJson}

		email, err := authDatabaseFile.AccountExists(context.Background(), "correct@answer.com")

		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if email != "correct@answer.com" {
			t.Error("expected email to be correct@answer.com but got: ", email)
		}
	})
}

func TestGetSession(t *testing.T) {
	err := os.WriteFile(testCsv,
		[]byte("abcde,test@user.com,2022-05-15T23:31:17+00:00"),
		0600,
	)
	if err != nil {
		panic("test failed to write testFile")
	}

	t.Run("should successfully get session", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: testCsv}

		accountDetails, err := authDatabaseFile.GetSession(context.Background(),
			"abcde",
			time.Date(2022, 05, 15, 00, 00, 00, 00, time.UTC),
		)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if accountDetails.Email != "test@user.com" {
			t.Error("expected email to be test@user.com but got: ", accountDetails.Email)
		}
	})
}

func TestFilesDoNotExist(t *testing.T) {
	t.Run("access session", func(t *testing.T) {
		authDatabaseFile := AuthDatabaseFile{SessionsFilepath: "does not exist"}

		_, err := authDatabaseFile.GetSession(context.Background(),
			"abcde",
			time.Date(2022, 05, 15, 00, 00, 00, 00, time.UTC),
		)
		if err != auth.ErrNotFound {
			t.Error("expected error to be ErrNotFound got: ", err)
		}
	})
}
