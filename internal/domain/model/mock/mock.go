// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/model/domain.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	domain "soundproof/internal/domain/model"

	gin "github.com/gin-gonic/gin"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CheckUserCredentials mocks base method.
func (m *MockStorage) CheckUserCredentials(ctx *gin.Context, req domain.LoginRequest) (bool, *domain.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserCredentials", ctx, req)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*domain.LoginResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CheckUserCredentials indicates an expected call of CheckUserCredentials.
func (mr *MockStorageMockRecorder) CheckUserCredentials(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserCredentials", reflect.TypeOf((*MockStorage)(nil).CheckUserCredentials), ctx, req)
}

// GetUserProfile mocks base method.
func (m *MockStorage) GetUserProfile(ctx *gin.Context, req domain.LoginRequest) (*domain.ProfileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, req)
	ret0, _ := ret[0].(*domain.ProfileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockStorageMockRecorder) GetUserProfile(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockStorage)(nil).GetUserProfile), ctx, req)
}

// RegisterUserInDB mocks base method.
func (m *MockStorage) RegisterUserInDB(ctx *gin.Context, req domain.UserRegistrationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUserInDB", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUserInDB indicates an expected call of RegisterUserInDB.
func (mr *MockStorageMockRecorder) RegisterUserInDB(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUserInDB", reflect.TypeOf((*MockStorage)(nil).RegisterUserInDB), ctx, req)
}

// UpdateUserProfile mocks base method.
func (m *MockStorage) UpdateUserProfile(ctx *gin.Context, address, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", ctx, address, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockStorageMockRecorder) UpdateUserProfile(ctx, address, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockStorage)(nil).UpdateUserProfile), ctx, address, email)
}

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

// CheckCredentials mocks base method.
func (m *MockService) CheckCredentials(c *gin.Context, req domain.LoginRequest) (bool, *domain.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCredentials", c, req)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*domain.LoginResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CheckCredentials indicates an expected call of CheckCredentials.
func (mr *MockServiceMockRecorder) CheckCredentials(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCredentials", reflect.TypeOf((*MockService)(nil).CheckCredentials), c, req)
}

// GetUserProfile mocks base method.
func (m *MockService) GetUserProfile(c *gin.Context, req domain.LoginRequest) (*domain.ProfileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", c, req)
	ret0, _ := ret[0].(*domain.ProfileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockServiceMockRecorder) GetUserProfile(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockService)(nil).GetUserProfile), c, req)
}

// RegisterUser mocks base method.
func (m *MockService) RegisterUser(c *gin.Context, req domain.UserRegistrationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", c, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockServiceMockRecorder) RegisterUser(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockService)(nil).RegisterUser), c, req)
}

// UpdateUser mocks base method.
func (m *MockService) UpdateUser(c *gin.Context, req domain.UpdateUserProfileRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", c, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockServiceMockRecorder) UpdateUser(c, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockService)(nil).UpdateUser), c, req)
}
