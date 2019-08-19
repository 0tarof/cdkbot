// Code generated by MockGen. DO NOT EDIT.
// Source: functions/operation/git/git.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	go_git_v4 "gopkg.in/src-d/go-git.v4"
	reflect "reflect"
)

// MockClienter is a mock of Clienter interface
type MockClienter struct {
	ctrl     *gomock.Controller
	recorder *MockClienterMockRecorder
}

// MockClienterMockRecorder is the mock recorder for MockClienter
type MockClienterMockRecorder struct {
	mock *MockClienter
}

// NewMockClienter creates a new mock instance
func NewMockClienter(ctrl *gomock.Controller) *MockClienter {
	mock := &MockClienter{ctrl: ctrl}
	mock.recorder = &MockClienterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClienter) EXPECT() *MockClienterMockRecorder {
	return m.recorder
}

// Clone mocks base method
func (m *MockClienter) Clone(path string, hash *string) (*go_git_v4.Worktree, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", path, hash)
	ret0, _ := ret[0].(*go_git_v4.Worktree)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clone indicates an expected call of Clone
func (mr *MockClienterMockRecorder) Clone(path, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockClienter)(nil).Clone), path, hash)
}

// Merge mocks base method
func (m *MockClienter) Merge(workTree *go_git_v4.Worktree, branch string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Merge", workTree, branch)
	ret0, _ := ret[0].(error)
	return ret0
}

// Merge indicates an expected call of Merge
func (mr *MockClienterMockRecorder) Merge(workTree, branch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Merge", reflect.TypeOf((*MockClienter)(nil).Merge), workTree, branch)
}
