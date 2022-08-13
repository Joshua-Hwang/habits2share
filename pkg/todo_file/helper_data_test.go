package todo_file

import (
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/todo"
)

func generateTestData() map[string]UsersTodos {
	testUsers := map[string]UsersTodos{
		"testUser1": {
			MyTodos: map[string]TodoJson{
				"testUser1_todoId1": {Todo: todo.Todo{Id: "testUser1_todoId1", Owner: "testUser1", Name: "hello there", Description: "first todo", DueDate: time.Now(), Completed: false}},
				"testUser1_todoId2": {Todo: todo.Todo{Id: "testUser1_todoId2", Owner: "testUser1", Name: "new todo", Description: "second todo", DueDate: time.Now(), Completed: false}},
			},
		},
		"testUser2": {
			MyTodos: map[string]TodoJson{
				"testUser2_todoId1": {Todo: todo.Todo{Id: "testUser2_todoId1",
					Owner:       "testUser2",
					Name:        "goodbye world",
					Description: "second\nitem",
					DueDate:     time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC),
					Completed:   false,
				}},
			},
		},
	}

	return testUsers
}
