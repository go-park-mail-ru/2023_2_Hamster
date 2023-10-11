// Code generated by MockGen. DO NOT EDIT.
// Source: request_handler.go

// Package mock_http is a generated GoMock package.
package mock_http

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUUIDParser is a mock of UUIDParser interface.
type MockUUIDParser struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDParserMockRecorder
}

// MockUUIDParserMockRecorder is the mock recorder for MockUUIDParser.
type MockUUIDParserMockRecorder struct {
	mock *MockUUIDParser
}

// NewMockUUIDParser creates a new mock instance.
func NewMockUUIDParser(ctrl *gomock.Controller) *MockUUIDParser {
	mock := &MockUUIDParser{ctrl: ctrl}
	mock.recorder = &MockUUIDParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDParser) EXPECT() *MockUUIDParserMockRecorder {
	return m.recorder
}

// GetIDFromRequest mocks base method.
func (m *MockUUIDParser) GetIDFromRequest(id string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIDFromRequest", id)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIDFromRequest indicates an expected call of GetIDFromRequest.
func (mr *MockUUIDParserMockRecorder) GetIDFromRequest(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIDFromRequest", reflect.TypeOf((*MockUUIDParser)(nil).GetIDFromRequest), id)
}
