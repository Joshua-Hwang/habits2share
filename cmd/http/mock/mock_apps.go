// Code generated by MockGen. DO NOT EDIT.
// Source: dependency_helper.go

// Package mock_main is a generated GoMock package.
package mock_main

import (
	reflect "reflect"
	time "time"

	habit_share "github.com/Joshua-Hwang/habits2share/pkg/habit_share"
	todo "github.com/Joshua-Hwang/habits2share/pkg/todo"
	gomock "github.com/golang/mock/gomock"
)

// MockHabitAppInterface is a mock of HabitAppInterface interface.
type MockHabitAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHabitAppInterfaceMockRecorder
}

// MockHabitAppInterfaceMockRecorder is the mock recorder for MockHabitAppInterface.
type MockHabitAppInterfaceMockRecorder struct {
	mock *MockHabitAppInterface
}

// NewMockHabitAppInterface creates a new mock instance.
func NewMockHabitAppInterface(ctrl *gomock.Controller) *MockHabitAppInterface {
	mock := &MockHabitAppInterface{ctrl: ctrl}
	mock.recorder = &MockHabitAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHabitAppInterface) EXPECT() *MockHabitAppInterfaceMockRecorder {
	return m.recorder
}

// ArchiveHabit mocks base method.
func (m *MockHabitAppInterface) ArchiveHabit(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveHabit", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ArchiveHabit indicates an expected call of ArchiveHabit.
func (mr *MockHabitAppInterfaceMockRecorder) ArchiveHabit(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).ArchiveHabit), id)
}

// ChangeDescription mocks base method.
func (m *MockHabitAppInterface) ChangeDescription(id, newDescription string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDescription", id, newDescription)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeDescription indicates an expected call of ChangeDescription.
func (mr *MockHabitAppInterfaceMockRecorder) ChangeDescription(id, newDescription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDescription", reflect.TypeOf((*MockHabitAppInterface)(nil).ChangeDescription), id, newDescription)
}

// ChangeFrequency mocks base method.
func (m *MockHabitAppInterface) ChangeFrequency(id string, newFrequency int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeFrequency", id, newFrequency)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeFrequency indicates an expected call of ChangeFrequency.
func (mr *MockHabitAppInterfaceMockRecorder) ChangeFrequency(id, newFrequency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeFrequency", reflect.TypeOf((*MockHabitAppInterface)(nil).ChangeFrequency), id, newFrequency)
}

// ChangeName mocks base method.
func (m *MockHabitAppInterface) ChangeName(id, newName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeName", id, newName)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeName indicates an expected call of ChangeName.
func (mr *MockHabitAppInterfaceMockRecorder) ChangeName(id, newName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeName", reflect.TypeOf((*MockHabitAppInterface)(nil).ChangeName), id, newName)
}

// CreateActivity mocks base method.
func (m *MockHabitAppInterface) CreateActivity(habitId string, logged habit_share.Time, status string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActivity", habitId, logged, status)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActivity indicates an expected call of CreateActivity.
func (mr *MockHabitAppInterfaceMockRecorder) CreateActivity(habitId, logged, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActivity", reflect.TypeOf((*MockHabitAppInterface)(nil).CreateActivity), habitId, logged, status)
}

// CreateHabit mocks base method.
func (m *MockHabitAppInterface) CreateHabit(name string, frequency int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHabit", name, frequency)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHabit indicates an expected call of CreateHabit.
func (mr *MockHabitAppInterfaceMockRecorder) CreateHabit(name, frequency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).CreateHabit), name, frequency)
}

// DeleteActivity mocks base method.
func (m *MockHabitAppInterface) DeleteActivity(habitId, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActivity", habitId, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActivity indicates an expected call of DeleteActivity.
func (mr *MockHabitAppInterfaceMockRecorder) DeleteActivity(habitId, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActivity", reflect.TypeOf((*MockHabitAppInterface)(nil).DeleteActivity), habitId, id)
}

// DeleteHabit mocks base method.
func (m *MockHabitAppInterface) DeleteHabit(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHabit", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHabit indicates an expected call of DeleteHabit.
func (mr *MockHabitAppInterfaceMockRecorder) DeleteHabit(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).DeleteHabit), id)
}

// GetActivities mocks base method.
func (m *MockHabitAppInterface) GetActivities(habitId string, after, before habit_share.Time, limit int) ([]habit_share.Activity, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivities", habitId, after, before, limit)
	ret0, _ := ret[0].([]habit_share.Activity)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetActivities indicates an expected call of GetActivities.
func (mr *MockHabitAppInterfaceMockRecorder) GetActivities(habitId, after, before, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivities", reflect.TypeOf((*MockHabitAppInterface)(nil).GetActivities), habitId, after, before, limit)
}

// GetHabit mocks base method.
func (m *MockHabitAppInterface) GetHabit(id string) (habit_share.Habit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHabit", id)
	ret0, _ := ret[0].(habit_share.Habit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHabit indicates an expected call of GetHabit.
func (mr *MockHabitAppInterfaceMockRecorder) GetHabit(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).GetHabit), id)
}

// GetMyHabits mocks base method.
func (m *MockHabitAppInterface) GetMyHabits(limit int, archived bool) ([]habit_share.Habit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyHabits", limit, archived)
	ret0, _ := ret[0].([]habit_share.Habit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyHabits indicates an expected call of GetMyHabits.
func (mr *MockHabitAppInterfaceMockRecorder) GetMyHabits(limit, archived interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyHabits", reflect.TypeOf((*MockHabitAppInterface)(nil).GetMyHabits), limit, archived)
}

// GetScore mocks base method.
func (m *MockHabitAppInterface) GetScore(habitId string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScore", habitId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScore indicates an expected call of GetScore.
func (mr *MockHabitAppInterfaceMockRecorder) GetScore(habitId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScore", reflect.TypeOf((*MockHabitAppInterface)(nil).GetScore), habitId)
}

// GetSharedHabits mocks base method.
func (m *MockHabitAppInterface) GetSharedHabits(limit int) ([]habit_share.Habit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSharedHabits", limit)
	ret0, _ := ret[0].([]habit_share.Habit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSharedHabits indicates an expected call of GetSharedHabits.
func (mr *MockHabitAppInterfaceMockRecorder) GetSharedHabits(limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSharedHabits", reflect.TypeOf((*MockHabitAppInterface)(nil).GetSharedHabits), limit)
}

// ShareHabit mocks base method.
func (m *MockHabitAppInterface) ShareHabit(habitId, friend string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShareHabit", habitId, friend)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShareHabit indicates an expected call of ShareHabit.
func (mr *MockHabitAppInterfaceMockRecorder) ShareHabit(habitId, friend interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShareHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).ShareHabit), habitId, friend)
}

// UnShareHabit mocks base method.
func (m *MockHabitAppInterface) UnShareHabit(habitId, friend string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnShareHabit", habitId, friend)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnShareHabit indicates an expected call of UnShareHabit.
func (mr *MockHabitAppInterfaceMockRecorder) UnShareHabit(habitId, friend interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnShareHabit", reflect.TypeOf((*MockHabitAppInterface)(nil).UnShareHabit), habitId, friend)
}

// MockTodoAppInterface is a mock of TodoAppInterface interface.
type MockTodoAppInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTodoAppInterfaceMockRecorder
}

// MockTodoAppInterfaceMockRecorder is the mock recorder for MockTodoAppInterface.
type MockTodoAppInterfaceMockRecorder struct {
	mock *MockTodoAppInterface
}

// NewMockTodoAppInterface creates a new mock instance.
func NewMockTodoAppInterface(ctrl *gomock.Controller) *MockTodoAppInterface {
	mock := &MockTodoAppInterface{ctrl: ctrl}
	mock.recorder = &MockTodoAppInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoAppInterface) EXPECT() *MockTodoAppInterfaceMockRecorder {
	return m.recorder
}

// ChangeDescription mocks base method.
func (m *MockTodoAppInterface) ChangeDescription(todoId, newDescription string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDescription", todoId, newDescription)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeDescription indicates an expected call of ChangeDescription.
func (mr *MockTodoAppInterfaceMockRecorder) ChangeDescription(todoId, newDescription interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDescription", reflect.TypeOf((*MockTodoAppInterface)(nil).ChangeDescription), todoId, newDescription)
}

// ChangeDueDate mocks base method.
func (m *MockTodoAppInterface) ChangeDueDate(todoId string, newTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDueDate", todoId, newTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeDueDate indicates an expected call of ChangeDueDate.
func (mr *MockTodoAppInterfaceMockRecorder) ChangeDueDate(todoId, newTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDueDate", reflect.TypeOf((*MockTodoAppInterface)(nil).ChangeDueDate), todoId, newTime)
}

// ChangeName mocks base method.
func (m *MockTodoAppInterface) ChangeName(todoId, newName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeName", todoId, newName)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeName indicates an expected call of ChangeName.
func (mr *MockTodoAppInterfaceMockRecorder) ChangeName(todoId, newName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeName", reflect.TypeOf((*MockTodoAppInterface)(nil).ChangeName), todoId, newName)
}

// CompleteTodo mocks base method.
func (m *MockTodoAppInterface) CompleteTodo(todoId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteTodo", todoId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteTodo indicates an expected call of CompleteTodo.
func (mr *MockTodoAppInterfaceMockRecorder) CompleteTodo(todoId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteTodo", reflect.TypeOf((*MockTodoAppInterface)(nil).CompleteTodo), todoId)
}

// CreateTodo mocks base method.
func (m *MockTodoAppInterface) CreateTodo(name string, dueDate time.Time) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTodo", name, dueDate)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTodo indicates an expected call of CreateTodo.
func (mr *MockTodoAppInterfaceMockRecorder) CreateTodo(name, dueDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTodo", reflect.TypeOf((*MockTodoAppInterface)(nil).CreateTodo), name, dueDate)
}

// GetMyTodos mocks base method.
func (m *MockTodoAppInterface) GetMyTodos(limit int, completed bool) ([]todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyTodos", limit, completed)
	ret0, _ := ret[0].([]todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyTodos indicates an expected call of GetMyTodos.
func (mr *MockTodoAppInterfaceMockRecorder) GetMyTodos(limit, completed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyTodos", reflect.TypeOf((*MockTodoAppInterface)(nil).GetMyTodos), limit, completed)
}

// GetTodo mocks base method.
func (m *MockTodoAppInterface) GetTodo(todoId string) (todo.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodo", todoId)
	ret0, _ := ret[0].(todo.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodo indicates an expected call of GetTodo.
func (mr *MockTodoAppInterfaceMockRecorder) GetTodo(todoId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodo", reflect.TypeOf((*MockTodoAppInterface)(nil).GetTodo), todoId)
}