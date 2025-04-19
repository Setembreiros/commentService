// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go

// Package mock_delete_comment is a generated GoMock package.
package mock_delete_comment

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// RemoveComment mocks base method.
func (m *MockService) RemoveComment(commentId uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveComment", commentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveComment indicates an expected call of RemoveComment.
func (mr *MockServiceMockRecorder) RemoveComment(commentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveComment", reflect.TypeOf((*MockService)(nil).RemoveComment), commentId)
}
