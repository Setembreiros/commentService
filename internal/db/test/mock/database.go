// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	model "commentservice/internal/model/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDatabaseClient is a mock of DatabaseClient interface.
type MockDatabaseClient struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseClientMockRecorder
}

// MockDatabaseClientMockRecorder is the mock recorder for MockDatabaseClient.
type MockDatabaseClientMockRecorder struct {
	mock *MockDatabaseClient
}

// NewMockDatabaseClient creates a new mock instance.
func NewMockDatabaseClient(ctrl *gomock.Controller) *MockDatabaseClient {
	mock := &MockDatabaseClient{ctrl: ctrl}
	mock.recorder = &MockDatabaseClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseClient) EXPECT() *MockDatabaseClientMockRecorder {
	return m.recorder
}

// CreateComment mocks base method.
func (m *MockDatabaseClient) CreateComment(data *model.Comment) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", data)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockDatabaseClientMockRecorder) CreateComment(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockDatabaseClient)(nil).CreateComment), data)
}
