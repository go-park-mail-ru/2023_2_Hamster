// Code generated by MockGen. DO NOT EDIT.
// Source: account.go

// Package mock_account is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
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

// CreateAccount mocks base method.
func (m *MockUsecase) CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, userID, account)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockUsecaseMockRecorder) CreateAccount(ctx, userID, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockUsecase)(nil).CreateAccount), ctx, userID, account)
}

// DeleteAccount mocks base method.
func (m *MockUsecase) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockUsecaseMockRecorder) DeleteAccount(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockUsecase)(nil).DeleteAccount), ctx, userID, accountID)
}

// UpdateAccount mocks base method.
func (m *MockUsecase) UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, userID, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockUsecaseMockRecorder) UpdateAccount(ctx, userID, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockUsecase)(nil).UpdateAccount), ctx, userID, account)
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

// AddUserInAccount mocks base method.
func (m *MockRepository) AddUserInAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserInAccount", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserInAccount indicates an expected call of AddUserInAccount.
func (mr *MockRepositoryMockRecorder) AddUserInAccount(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserInAccount", reflect.TypeOf((*MockRepository)(nil).AddUserInAccount), ctx, userID, accountID)
}

// CheckDuplicate mocks base method.
func (m *MockRepository) CheckDuplicate(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicate", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckDuplicate indicates an expected call of CheckDuplicate.
func (mr *MockRepositoryMockRecorder) CheckDuplicate(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicate", reflect.TypeOf((*MockRepository)(nil).CheckDuplicate), ctx, userID, accountID)
}

// CheckForbidden mocks base method.
func (m *MockRepository) CheckForbidden(ctx context.Context, accountID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckForbidden", ctx, accountID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckForbidden indicates an expected call of CheckForbidden.
func (mr *MockRepositoryMockRecorder) CheckForbidden(ctx, accountID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckForbidden", reflect.TypeOf((*MockRepository)(nil).CheckForbidden), ctx, accountID, userID)
}

// CreateAccount mocks base method.
func (m *MockRepository) CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, userID, account)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockRepositoryMockRecorder) CreateAccount(ctx, userID, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockRepository)(nil).CreateAccount), ctx, userID, account)
}

// DeleteAccount mocks base method.
func (m *MockRepository) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockRepositoryMockRecorder) DeleteAccount(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockRepository)(nil).DeleteAccount), ctx, userID, accountID)
}

// DeleteUserInAccount mocks base method.
func (m *MockRepository) DeleteUserInAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserInAccount", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserInAccount indicates an expected call of DeleteUserInAccount.
func (mr *MockRepositoryMockRecorder) DeleteUserInAccount(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserInAccount", reflect.TypeOf((*MockRepository)(nil).DeleteUserInAccount), ctx, userID, accountID)
}

// SharingCheck mocks base method.
func (m *MockRepository) SharingCheck(ctx context.Context, accountID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SharingCheck", ctx, accountID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SharingCheck indicates an expected call of SharingCheck.
func (mr *MockRepositoryMockRecorder) SharingCheck(ctx, accountID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SharingCheck", reflect.TypeOf((*MockRepository)(nil).SharingCheck), ctx, accountID, userID)
}

// Unsubscribe mocks base method.
func (m *MockRepository) Unsubscribe(ctx context.Context, userID, accountID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", ctx, userID, accountID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockRepositoryMockRecorder) Unsubscribe(ctx, userID, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockRepository)(nil).Unsubscribe), ctx, userID, accountID)
}

// UpdateAccount mocks base method.
func (m *MockRepository) UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, userID, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockRepositoryMockRecorder) UpdateAccount(ctx, userID, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockRepository)(nil).UpdateAccount), ctx, userID, account)
}
