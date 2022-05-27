// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sortednet/statuschecker/internal/statuschecker (interfaces: DbQuery)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	store "github.com/sortednet/statuschecker/internal/store"
)

// MockDbQuery is a mock of DbQuery interface.
type MockDbQuery struct {
	ctrl     *gomock.Controller
	recorder *MockDbQueryMockRecorder
}

// MockDbQueryMockRecorder is the mock recorder for MockDbQuery.
type MockDbQueryMockRecorder struct {
	mock *MockDbQuery
}

// NewMockDbQuery creates a new mock instance.
func NewMockDbQuery(ctrl *gomock.Controller) *MockDbQuery {
	mock := &MockDbQuery{ctrl: ctrl}
	mock.recorder = &MockDbQueryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDbQuery) EXPECT() *MockDbQueryMockRecorder {
	return m.recorder
}

// GetServices mocks base method.
func (m *MockDbQuery) GetServices(arg0 context.Context) ([]store.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServices", arg0)
	ret0, _ := ret[0].([]store.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServices indicates an expected call of GetServices.
func (mr *MockDbQueryMockRecorder) GetServices(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServices", reflect.TypeOf((*MockDbQuery)(nil).GetServices), arg0)
}

// RegisterService mocks base method.
func (m *MockDbQuery) RegisterService(arg0 context.Context, arg1 store.RegisterServiceParams) (store.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterService", arg0, arg1)
	ret0, _ := ret[0].(store.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterService indicates an expected call of RegisterService.
func (mr *MockDbQueryMockRecorder) RegisterService(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterService", reflect.TypeOf((*MockDbQuery)(nil).RegisterService), arg0, arg1)
}

// UnregisterService mocks base method.
func (m *MockDbQuery) UnregisterService(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterService", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterService indicates an expected call of UnregisterService.
func (mr *MockDbQueryMockRecorder) UnregisterService(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterService", reflect.TypeOf((*MockDbQuery)(nil).UnregisterService), arg0, arg1)
}
