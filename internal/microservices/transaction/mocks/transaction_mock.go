// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock_transaction is a generated GoMock package.
package mock_transaction

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

// CreateTransaction mocks base method.
func (m *MockUsecase) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, transaction)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockUsecaseMockRecorder) CreateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockUsecase)(nil).CreateTransaction), ctx, transaction)
}

// DeleteTransaction mocks base method.
func (m *MockUsecase) DeleteTransaction(ctx context.Context, transactionID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", ctx, transactionID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockUsecaseMockRecorder) DeleteTransaction(ctx, transactionID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockUsecase)(nil).DeleteTransaction), ctx, transactionID, userID)
}

// GetFeed mocks base method.
func (m *MockUsecase) GetFeed(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Transaction, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, userID, page, pageSize)
	ret0, _ := ret[0].([]models.Transaction)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockUsecaseMockRecorder) GetFeed(ctx, userID, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockUsecase)(nil).GetFeed), ctx, userID, page, pageSize)
}

// UpdateTransaction mocks base method.
func (m *MockUsecase) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockUsecaseMockRecorder) UpdateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockUsecase)(nil).UpdateTransaction), ctx, transaction)
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

// Check mocks base method.
func (m *MockRepository) Check(ctx context.Context, transactionID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, transactionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRepositoryMockRecorder) Check(ctx, transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), ctx, transactionID)
}

// CheckForbidden mocks base method.
func (m *MockRepository) CheckForbidden(ctx context.Context, transactinID uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckForbidden", ctx, transactinID)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckForbidden indicates an expected call of CheckForbidden.
func (mr *MockRepositoryMockRecorder) CheckForbidden(ctx, transactinID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckForbidden", reflect.TypeOf((*MockRepository)(nil).CheckForbidden), ctx, transactinID)
}

// CreateTransaction mocks base method.
func (m *MockRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, transaction)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockRepositoryMockRecorder) CreateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockRepository)(nil).CreateTransaction), ctx, transaction)
}

// DeleteTransaction mocks base method.
func (m *MockRepository) DeleteTransaction(ctx context.Context, transactionID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", ctx, transactionID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockRepositoryMockRecorder) DeleteTransaction(ctx, transactionID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockRepository)(nil).DeleteTransaction), ctx, transactionID, userID)
}

// GetFeed mocks base method.
func (m *MockRepository) GetFeed(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Transaction, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, userID, page, pageSize)
	ret0, _ := ret[0].([]models.Transaction)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockRepositoryMockRecorder) GetFeed(ctx, userID, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockRepository)(nil).GetFeed), ctx, userID, page, pageSize)
}

// UpdateTransaction mocks base method.
func (m *MockRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", ctx, transaction)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockRepositoryMockRecorder) UpdateTransaction(ctx, transaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockRepository)(nil).UpdateTransaction), ctx, transaction)
}