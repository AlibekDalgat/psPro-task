// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "psPro-task/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// CreateCommand mocks base method.
func (m *MockCommand) CreateCommand(command models.Command) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommand", command)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommand indicates an expected call of CreateCommand.
func (mr *MockCommandMockRecorder) CreateCommand(command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommand", reflect.TypeOf((*MockCommand)(nil).CreateCommand), command)
}

// ExecuteCommand mocks base method.
func (m *MockCommand) ExecuteCommand(arg0 int, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ExecuteCommand", arg0, arg1)
}

// ExecuteCommand indicates an expected call of ExecuteCommand.
func (mr *MockCommandMockRecorder) ExecuteCommand(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteCommand", reflect.TypeOf((*MockCommand)(nil).ExecuteCommand), arg0, arg1)
}

// GetAllCommands mocks base method.
func (m *MockCommand) GetAllCommands() ([]models.CommResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCommands")
	ret0, _ := ret[0].([]models.CommResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCommands indicates an expected call of GetAllCommands.
func (mr *MockCommandMockRecorder) GetAllCommands() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCommands", reflect.TypeOf((*MockCommand)(nil).GetAllCommands))
}

// GetOneCommand mocks base method.
func (m *MockCommand) GetOneCommand(arg0 int) (models.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneCommand", arg0)
	ret0, _ := ret[0].(models.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneCommand indicates an expected call of GetOneCommand.
func (mr *MockCommandMockRecorder) GetOneCommand(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneCommand", reflect.TypeOf((*MockCommand)(nil).GetOneCommand), arg0)
}

// KillCommand mocks base method.
func (m *MockCommand) KillCommand(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KillCommand", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// KillCommand indicates an expected call of KillCommand.
func (mr *MockCommandMockRecorder) KillCommand(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KillCommand", reflect.TypeOf((*MockCommand)(nil).KillCommand), arg0)
}

// StartCommand mocks base method.
func (m *MockCommand) StartCommand(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartCommand", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartCommand indicates an expected call of StartCommand.
func (mr *MockCommandMockRecorder) StartCommand(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartCommand", reflect.TypeOf((*MockCommand)(nil).StartCommand), arg0)
}

// StopCommand mocks base method.
func (m *MockCommand) StopCommand(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopCommand", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopCommand indicates an expected call of StopCommand.
func (mr *MockCommandMockRecorder) StopCommand(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopCommand", reflect.TypeOf((*MockCommand)(nil).StopCommand), arg0)
}
