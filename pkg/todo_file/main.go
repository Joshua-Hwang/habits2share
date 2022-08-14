package todo_file

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Joshua-Hwang/habits2share/pkg/todo"
	todo_app "github.com/Joshua-Hwang/habits2share/pkg/todo"
	"github.com/google/uuid"
)

// TTL in seconds
const cacheTtl = 10

type TodoJson struct {
	todo_app.Todo
}

type UsersTodos struct {
	MyTodos map[string]TodoJson
}

// TODO could probably split this into a file per user
type TodoFile struct {
	UsersTodos map[string]UsersTodos
	filename   string
	fileLock   *sync.Mutex // This can't be a rw mutex as you're always "writing" the parsed file to the struct
	lastRead   time.Time
}

var _ todo_app.TodoDatabase = (*TodoFile)(nil)

func ConstructTodoId(owner string, postfix string) string {
	return fmt.Sprintf("%s_%s", owner, postfix)
}

func parseTodoId(todoId string) (owner string, postfix string, err error) {
	lastIndex := strings.LastIndex(todoId, "_")
	if lastIndex == -1 {
		err = &todo_app.InputError{Message: fmt.Sprintf("Todo id could not be parsed. Todo id: %s", todoId)}
		return
	}

	owner = todoId[:lastIndex]
	postfix = todoId[lastIndex+1:]
	return
}

func TodoFromFile(filename string) (*TodoFile, error) {
	var todoFile TodoFile
	todoFile.filename = filename
	todoFile.fileLock = &sync.Mutex{}

	err := todoFile.read()

	if err != nil {
		return nil, err
	}

	return &todoFile, nil
}

func (a *TodoFile) read() error {
	if a.filename != "" && time.Since(a.lastRead) > time.Duration(cacheTtl*float64(time.Second)) {
		a.fileLock.Lock()
		defer a.fileLock.Unlock()

		content, err := os.ReadFile(a.filename)
		a.lastRead = time.Now()
		if err != nil || len(content) == 0 {
			if !os.IsNotExist(err) {
				return err
			}
			// file does not exist or got removed
			a.UsersTodos = make(map[string]UsersTodos, 0)
			return nil
		}
		err = json.Unmarshal(content, a)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (a *TodoFile) write() error {
	if a.filename != "" {
		a.fileLock.Lock()
		defer a.fileLock.Unlock()

		file, err := os.OpenFile(a.filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		// TODO generating the marshalled value should be separate from writing the file (testing purposes)
		jsonString, err := json.MarshalIndent(a, "", " ")
		if err != nil {
			return err
		}
		_, err = file.Write(jsonString)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

// ChangeDescription implements todo.TodoDatabase
func (a *TodoFile) ChangeDescription(id string, newDescription string) error {
	if err := a.read(); err != nil {
		return err
	}

	userId, _, err := parseTodoId(id)
	if err != nil {
		return err
	}

	if user, ok := a.UsersTodos[userId]; !ok {
		return todo_app.UserNotFoundError
	} else if todo, ok := user.MyTodos[id]; !ok {
		return todo_app.TodoNotFoundError
	} else {
		todo.Description = newDescription
		user.MyTodos[id] = todo
		a.UsersTodos[userId] = user // TODO this is a concurrent write to the map. Please fix
	}

	err = a.write()
	if err != nil {
		return err
	}

	return nil
}

// ChangeDueDate implements todo.TodoDatabase
func (a *TodoFile) ChangeDueDate(id string, newTime time.Time) error {
	if err := a.read(); err != nil {
		return err
	}

	userId, _, err := parseTodoId(id)
	if err != nil {
		return err
	}

	if user, ok := a.UsersTodos[userId]; !ok {
		return todo_app.UserNotFoundError
	} else if todo, ok := user.MyTodos[id]; !ok {
		return todo_app.TodoNotFoundError
	} else {
		todo.DueDate = newTime
		user.MyTodos[id] = todo
		a.UsersTodos[userId] = user
	}

	err = a.write()
	if err != nil {
		return err
	}

	return nil
}

// ChangeName implements todo.TodoDatabase
func (a *TodoFile) ChangeName(id string, newName string) error {
	if err := a.read(); err != nil {
		return err
	}

	userId, _, err := parseTodoId(id)
	if err != nil {
		return err
	}

	if user, ok := a.UsersTodos[userId]; !ok {
		return todo_app.UserNotFoundError
	} else if todo, ok := user.MyTodos[id]; !ok {
		return todo_app.TodoNotFoundError
	} else {
		todo.Name = newName
		user.MyTodos[id] = todo
		a.UsersTodos[userId] = user
	}

	err = a.write()
	if err != nil {
		return err
	}

	return nil
}

// CompleteTodo implements todo.TodoDatabase
func (a *TodoFile) CompleteTodo(id string) error {
	if err := a.read(); err != nil {
		return err
	}

	userId, _, err := parseTodoId(id)
	if err != nil {
		return err
	}

	if user, ok := a.UsersTodos[userId]; !ok {
		return todo_app.UserNotFoundError
	} else if todo, ok := user.MyTodos[id]; !ok {
		return todo_app.TodoNotFoundError
	} else {
		todo.Completed = true
		user.MyTodos[id] = todo
		a.UsersTodos[userId] = user
	}

	err = a.write()
	if err != nil {
		return err
	}

	return nil
}

// CreateTodo implements todo.TodoDatabase
func (a *TodoFile) CreateTodo(name string, owner string, dueDate time.Time) (string, error) {
	if err := a.read(); err != nil {
		return "", err
	}

	user, ok := a.UsersTodos[owner]
	if !ok {
		// if user doesn't exist create user
		user = UsersTodos{MyTodos: make(map[string]TodoJson, 0)}
		a.UsersTodos[owner] = user
	}

	todo := todo.Todo{
		Id:        ConstructTodoId(owner, uuid.NewString()),
		Owner:     owner,
		Name:      name,
		DueDate:   dueDate,
		Completed: false,
	}

	a.UsersTodos[owner].MyTodos[todo.Id] = TodoJson{todo}

	err := a.write()
	if err != nil {
		return todo.Id, err
	}

	return todo.Id, nil
}

// GetTodo implements todo.TodoDatabase
func (a *TodoFile) GetTodo(id string) (todo_app.Todo, error) {
	ownerId, _, err := parseTodoId(id)
	if err != nil {
		return todo_app.Todo{}, err
	}
	if err := a.read(); err != nil {
		return todo_app.Todo{}, err
	}

	user, ok := a.UsersTodos[ownerId]
	if !ok {
		return todo_app.Todo{}, todo_app.TodoNotFoundError
	}

	todo, ok := user.MyTodos[id]
	if !ok {
		return todo_app.Todo{}, todo_app.TodoNotFoundError
	}

	return todo.Todo, nil
}

// GetTodosByOwner implements todo.TodoDatabase
func (a *TodoFile) GetTodosByOwner(ownerId string, limit int, completed bool) ([]todo_app.Todo, error) {
	if err := a.read(); err != nil {
		return nil, err
	}

	user, ok := a.UsersTodos[ownerId]
	if !ok {
		return make([]todo_app.Todo, 0), nil
	}

	myTodos := make([]todo_app.Todo, 0, len(user.MyTodos))
	for todoId := range user.MyTodos {
		todo := user.MyTodos[todoId]

		if !todo.Completed || completed {
			myTodos = append(myTodos, todo.Todo)
		}
	}

	// map does not guaratee this is in order
	sort.Slice(myTodos[:], func(i, j int) bool {
		return myTodos[i].Name < myTodos[j].Name
	})

	return myTodos, nil
}
