package todo_file

import (
	"sync"
	"testing"
	"time"
)

func TestTodo(t *testing.T) {
	t.Run("should create todo", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		todoId, err := todoFile.CreateTodo("new todo", "testUser1", time.Now())
		if err != nil || todoId == "" {
			t.Error("expected error to be nil got: ", err)
		}
	})

	t.Run("should return todo", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		todo, err := todoFile.GetTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		// TODO think of better way to do equality for the object
		// This is assumed based on generateTestData
		if todo.Name != "hello there" {
			t.Error("expected todo object to have name \"hello there\" got: ", todo.Name)
		}
	})

	t.Run("should get todos by owner", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		// limit is more than number of todos as a test
		todos, err := todoFile.GetTodosByOwner("testUser2", 10, false)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}
		if len(todos) != 1 {
			t.Error("expected the number of todos to be 1 got: ", len(todos))
		}
	})

	t.Run("should change name", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		// limit is more than number of todos as a test
		err := todoFile.ChangeName("testUser1_todoId1", "new name")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		todo, err := todoFile.GetTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if todo.Name != "new name" {
			t.Error("expected todo object change name to \"new name\" got: ", todo.Name)
		}
	})

	t.Run("should change description", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		// limit is more than number of todos as a test
		err := todoFile.ChangeDescription("testUser1_todoId1", "new description")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		todo, err := todoFile.GetTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if todo.Description != "new description" {
			t.Error("expected todo object change name to \"new description\" got: ", todo.Description)
		}
	})

	t.Run("should change due dates", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		// limit is more than number of todos as a test
		newDate := time.Date(2022, time.August, 13, 0, 0, 0, 0, time.UTC);
		err := todoFile.ChangeDueDate("testUser1_todoId1", newDate)
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		todo, err := todoFile.GetTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if todo.DueDate != newDate {
			t.Error("expected todo object change due date to ", newDate, " got: ", todo.DueDate)
		}
	})

	t.Run("should complete todo", func(t *testing.T) {
		testData := generateTestData()
		todoFile := TodoFile {
			UsersTodos: testData,
			filename: outputJson,
			fileLock: &sync.Mutex{},
		}

		// limit is more than number of todos as a test
		err := todoFile.CompleteTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		todo, err := todoFile.GetTodo("testUser1_todoId1")
		if err != nil {
			t.Error("expected error to be nil got: ", err)
		}

		if todo.Completed != true {
			t.Error("expected todo object change completion status but was not completed")
		}
	})
}
