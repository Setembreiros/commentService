// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_create_comment is a generated GoMock package.
package mock_create_comment

import (
	model "commentservice/internal/model/domain"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddComment mocks base method.
func (m *MockRepository) AddComment(data *model.Comment) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", data)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddComment indicates an expected call of AddComment.
func (mr *MockRepositoryMockRecorder) AddComment(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockRepository)(nil).AddComment), data)
}

// MockTimeService is a mock of TimeService interface.
type MockTimeService struct {
	ctrl     *gomock.Controller
	recorder *MockTimeServiceMockRecorder
}

// MockTimeServiceMockRecorder is the mock recorder for MockTimeService.
type MockTimeServiceMockRecorder struct {
	mock *MockTimeService
}

// NewMockTimeService creates a new mock instance.
func NewMockTimeService(ctrl *gomock.Controller) *MockTimeService {
	mock := &MockTimeService{ctrl: ctrl}
	mock.recorder = &MockTimeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeService) EXPECT() *MockTimeServiceMockRecorder {
	return m.recorder
}

// GetTimeNowUtc mocks base method.
func (m *MockTimeService) GetTimeNowUtc() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimeNowUtc")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTimeNowUtc indicates an expected call of GetTimeNowUtc.
func (mr *MockTimeServiceMockRecorder) GetTimeNowUtc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimeNowUtc", reflect.TypeOf((*MockTimeService)(nil).GetTimeNowUtc))
}
