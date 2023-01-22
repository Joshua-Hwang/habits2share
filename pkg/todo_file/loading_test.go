package todo_file

import (
	"sync"
	"testing"
)

var inputJson = "input.json"

func TestIo(t *testing.T) {
	t.Run("should write to a file", func(t *testing.T) {
		tempDir := t.TempDir()
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: tempDir + "/output.json",
			fileLock: &sync.Mutex{},
		}

		err := todoFile.write()
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
	})

	t.Run("should read from file that doesn't exist", func(t *testing.T) {
		todoFile := TodoFile {
			filename: "doesn't exist",
			fileLock: &sync.Mutex{},
		}

		err := todoFile.read()
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
	})

	t.Run("should read from a file", func(t *testing.T) {
		todoFile := TodoFile {
			filename: inputJson,
			fileLock: &sync.Mutex{},
		}

		err := todoFile.read()
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
	})

	t.Run("should read from file that is written", func(t *testing.T) {
		tempDir := t.TempDir()
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: tempDir + "/output.json",
			fileLock: &sync.Mutex{},
		}

		err := todoFile.write()
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		todoFile2 := TodoFile {
			filename: tempDir + "/output.json",
			fileLock: &sync.Mutex{},
		}

		err = todoFile2.read()
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		// TODO probably a nicer way of doing this (for another time)
		if _, ok := todoFile.UsersTodos["testUser1"]; !ok {
			t.Errorf("Failed to correctly parse the input json %+v", todoFile)
		}
	})
}
