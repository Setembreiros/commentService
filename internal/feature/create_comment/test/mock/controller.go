// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go

// Package mock_create_comment is a generated GoMock package.
package mock_create_comment

import (
	model "commentservice/internal/model/domain"
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

// AddComment mocks base method.
func (m *MockService) AddComment(comment *model.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddComment indicates an expected call of AddComment.
func (mr *MockServiceMockRecorder) AddComment(comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockService)(nil).AddComment), comment)
}
