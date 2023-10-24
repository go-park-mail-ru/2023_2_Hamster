// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	auth "github.com/go-park-mail-ru/2023_2_Hamster/internal/pkg/auth"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *MockUsecase) GenerateAccessToken(ctx context.Context, user models.User) (auth.CookieToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", ctx, user)
	ret0, _ := ret[0].(auth.CookieToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockUsecaseMockRecorder) GenerateAccessToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockUsecase)(nil).GenerateAccessToken), ctx, user)
}

// GetUserByAuthData mocks base method.
func (m *MockUsecase) GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByAuthData", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByAuthData indicates an expected call of GetUserByAuthData.
func (mr *MockUsecaseMockRecorder) GetUserByAuthData(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByAuthData", reflect.TypeOf((*MockUsecase)(nil).GetUserByAuthData), ctx, userID)
}

// GetUserByCreds mocks base method.
func (m *MockUsecase) GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCreds", ctx, username, plainPassword)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCreds indicates an expected call of GetUserByCreds.
func (mr *MockUsecaseMockRecorder) GetUserByCreds(ctx, username, plainPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCreds", reflect.TypeOf((*MockUsecase)(nil).GetUserByCreds), ctx, username, plainPassword)
}

// SignInUser mocks base method.
func (m *MockUsecase) SignInUser(username, plainPassword string) (uuid.UUID, auth.CookieToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInUser", username, plainPassword)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(auth.CookieToken)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignInUser indicates an expected call of SignInUser.
func (mr *MockUsecaseMockRecorder) SignInUser(username, plainPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInUser", reflect.TypeOf((*MockUsecase)(nil).SignInUser), username, plainPassword)
}

// SignUpUser mocks base method.
func (m *MockUsecase) SignUpUser(user models.User) (uuid.UUID, auth.CookieToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUpUser", user)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(auth.CookieToken)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignUpUser indicates an expected call of SignUpUser.
func (mr *MockUsecaseMockRecorder) SignUpUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUpUser", reflect.TypeOf((*MockUsecase)(nil).SignUpUser), user)
}

// ValidateAccessToken mocks base method.
func (m *MockUsecase) ValidateAccessToken(accessToken string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccessToken", accessToken)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAccessToken indicates an expected call of ValidateAccessToken.
func (mr *MockUsecaseMockRecorder) ValidateAccessToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccessToken", reflect.TypeOf((*MockUsecase)(nil).ValidateAccessToken), accessToken)
}

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

// GetUserByAuthData mocks base method.
func (m *MockRepository) GetUserByAuthData(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByAuthData", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByAuthData indicates an expected call of GetUserByAuthData.
func (mr *MockRepositoryMockRecorder) GetUserByAuthData(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByAuthData", reflect.TypeOf((*MockRepository)(nil).GetUserByAuthData), ctx, userID)
}