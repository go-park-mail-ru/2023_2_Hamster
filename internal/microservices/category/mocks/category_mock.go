// Code generated by MockGen. DO NOT EDIT.
// Source: category.go

// Package mock_category is a generated GoMock package.
package mock_category

import (
	context "context"
	reflect "reflect"

	category "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category"
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

// CreateTag mocks base method.
func (m *MockUsecase) CreateTag(ctx context.Context, tag category.TagInput) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTag", ctx, tag)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTag indicates an expected call of CreateTag.
func (mr *MockUsecaseMockRecorder) CreateTag(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTag", reflect.TypeOf((*MockUsecase)(nil).CreateTag), ctx, tag)
}

// DeleteTag mocks base method.
func (m *MockUsecase) DeleteTag(ctx context.Context, tagId, userId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, tagId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockUsecaseMockRecorder) DeleteTag(ctx, tagId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockUsecase)(nil).DeleteTag), ctx, tagId, userId)
}

// GetTags mocks base method.
func (m *MockUsecase) GetTags(ctx context.Context, userId uuid.UUID) ([]models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", ctx, userId)
	ret0, _ := ret[0].([]models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockUsecaseMockRecorder) GetTags(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockUsecase)(nil).GetTags), ctx, userId)
}

// UpdateTag mocks base method.
func (m *MockUsecase) UpdateTag(ctx context.Context, tag *models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTag", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTag indicates an expected call of UpdateTag.
func (mr *MockUsecaseMockRecorder) UpdateTag(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTag", reflect.TypeOf((*MockUsecase)(nil).UpdateTag), ctx, tag)
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

// CheckExist mocks base method.
func (m *MockRepository) CheckExist(ctx context.Context, userId, tagId uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExist", ctx, userId, tagId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckExist indicates an expected call of CheckExist.
func (mr *MockRepositoryMockRecorder) CheckExist(ctx, userId, tagId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckExist", reflect.TypeOf((*MockRepository)(nil).CheckExist), ctx, userId, tagId)
}

// CheckNameUniq mocks base method.
func (m *MockRepository) CheckNameUniq(ctx context.Context, userId, parentId uuid.UUID, name string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNameUniq", ctx, userId, parentId, name)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckNameUniq indicates an expected call of CheckNameUniq.
func (mr *MockRepositoryMockRecorder) CheckNameUniq(ctx, userId, parentId, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNameUniq", reflect.TypeOf((*MockRepository)(nil).CheckNameUniq), ctx, userId, parentId, name)
}

// CreateTag mocks base method.
func (m *MockRepository) CreateTag(ctx context.Context, category models.Category) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTag", ctx, category)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTag indicates an expected call of CreateTag.
func (mr *MockRepositoryMockRecorder) CreateTag(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTag", reflect.TypeOf((*MockRepository)(nil).CreateTag), ctx, category)
}

// DeleteTag mocks base method.
func (m *MockRepository) DeleteTag(ctx context.Context, tagId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, tagId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockRepositoryMockRecorder) DeleteTag(ctx, tagId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockRepository)(nil).DeleteTag), ctx, tagId)
}

// GetTags mocks base method.
func (m *MockRepository) GetTags(ctx context.Context, userId uuid.UUID) ([]models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", ctx, userId)
	ret0, _ := ret[0].([]models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockRepositoryMockRecorder) GetTags(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockRepository)(nil).GetTags), ctx, userId)
}

// UpdateTag mocks base method.
func (m *MockRepository) UpdateTag(ctx context.Context, tag *models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTag", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTag indicates an expected call of UpdateTag.
func (mr *MockRepositoryMockRecorder) UpdateTag(ctx, tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTag", reflect.TypeOf((*MockRepository)(nil).UpdateTag), ctx, tag)
}
