// Code generated by MockGen. DO NOT EDIT.
// Source: delivery/grpc/generated/auth_grpc.pb.go

// Package mock___ is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	__ "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockAuthServiceClient is a mock of AuthServiceClient interface.
type MockAuthServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceClientMockRecorder
}

// MockAuthServiceClientMockRecorder is the mock recorder for MockAuthServiceClient.
type MockAuthServiceClientMockRecorder struct {
	mock *MockAuthServiceClient
}

// NewMockAuthServiceClient creates a new mock instance.
func NewMockAuthServiceClient(ctrl *gomock.Controller) *MockAuthServiceClient {
	mock := &MockAuthServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceClient) EXPECT() *MockAuthServiceClientMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockAuthServiceClient) ChangePassword(ctx context.Context, in *__.ChangePasswordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangePassword", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthServiceClientMockRecorder) ChangePassword(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuthServiceClient)(nil).ChangePassword), varargs...)
}

// CheckLoginUnique mocks base method.
func (m *MockAuthServiceClient) CheckLoginUnique(ctx context.Context, in *__.UniqCheckRequest, opts ...grpc.CallOption) (*__.UniqCheckResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckLoginUnique", varargs...)
	ret0, _ := ret[0].(*__.UniqCheckResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLoginUnique indicates an expected call of CheckLoginUnique.
func (mr *MockAuthServiceClientMockRecorder) CheckLoginUnique(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLoginUnique", reflect.TypeOf((*MockAuthServiceClient)(nil).CheckLoginUnique), varargs...)
}

// GetByID mocks base method.
func (m *MockAuthServiceClient) GetByID(ctx context.Context, in *__.UserIdRequest, opts ...grpc.CallOption) (*__.UserResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByID", varargs...)
	ret0, _ := ret[0].(*__.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAuthServiceClientMockRecorder) GetByID(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAuthServiceClient)(nil).GetByID), varargs...)
}

// Login mocks base method.
func (m *MockAuthServiceClient) Login(ctx context.Context, in *__.LoginRequest, opts ...grpc.CallOption) (*__.LoginResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*__.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceClientMockRecorder) Login(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServiceClient)(nil).Login), varargs...)
}

// SignUp mocks base method.
func (m *MockAuthServiceClient) SignUp(ctx context.Context, in *__.SignUpRequest, opts ...grpc.CallOption) (*__.SignUpResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SignUp", varargs...)
	ret0, _ := ret[0].(*__.SignUpResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceClientMockRecorder) SignUp(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthServiceClient)(nil).SignUp), varargs...)
}

// MockAuthServiceServer is a mock of AuthServiceServer interface.
type MockAuthServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceServerMockRecorder
}

// MockAuthServiceServerMockRecorder is the mock recorder for MockAuthServiceServer.
type MockAuthServiceServerMockRecorder struct {
	mock *MockAuthServiceServer
}

// NewMockAuthServiceServer creates a new mock instance.
func NewMockAuthServiceServer(ctrl *gomock.Controller) *MockAuthServiceServer {
	mock := &MockAuthServiceServer{ctrl: ctrl}
	mock.recorder = &MockAuthServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServiceServer) EXPECT() *MockAuthServiceServerMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockAuthServiceServer) ChangePassword(arg0 context.Context, arg1 *__.ChangePasswordRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthServiceServerMockRecorder) ChangePassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuthServiceServer)(nil).ChangePassword), arg0, arg1)
}

// CheckLoginUnique mocks base method.
func (m *MockAuthServiceServer) CheckLoginUnique(arg0 context.Context, arg1 *__.UniqCheckRequest) (*__.UniqCheckResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckLoginUnique", arg0, arg1)
	ret0, _ := ret[0].(*__.UniqCheckResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckLoginUnique indicates an expected call of CheckLoginUnique.
func (mr *MockAuthServiceServerMockRecorder) CheckLoginUnique(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckLoginUnique", reflect.TypeOf((*MockAuthServiceServer)(nil).CheckLoginUnique), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockAuthServiceServer) GetByID(arg0 context.Context, arg1 *__.UserIdRequest) (*__.UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*__.UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAuthServiceServerMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAuthServiceServer)(nil).GetByID), arg0, arg1)
}

// Login mocks base method.
func (m *MockAuthServiceServer) Login(arg0 context.Context, arg1 *__.LoginRequest) (*__.LoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*__.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServiceServer)(nil).Login), arg0, arg1)
}

// SignUp mocks base method.
func (m *MockAuthServiceServer) SignUp(arg0 context.Context, arg1 *__.SignUpRequest) (*__.SignUpResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1)
	ret0, _ := ret[0].(*__.SignUpResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceServerMockRecorder) SignUp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthServiceServer)(nil).SignUp), arg0, arg1)
}

// mustEmbedUnimplementedAuthServiceServer mocks base method.
func (m *MockAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServiceServer")
}

// mustEmbedUnimplementedAuthServiceServer indicates an expected call of mustEmbedUnimplementedAuthServiceServer.
func (mr *MockAuthServiceServerMockRecorder) mustEmbedUnimplementedAuthServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServiceServer", reflect.TypeOf((*MockAuthServiceServer)(nil).mustEmbedUnimplementedAuthServiceServer))
}

// MockUnsafeAuthServiceServer is a mock of UnsafeAuthServiceServer interface.
type MockUnsafeAuthServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthServiceServerMockRecorder
}

// MockUnsafeAuthServiceServerMockRecorder is the mock recorder for MockUnsafeAuthServiceServer.
type MockUnsafeAuthServiceServerMockRecorder struct {
	mock *MockUnsafeAuthServiceServer
}

// NewMockUnsafeAuthServiceServer creates a new mock instance.
func NewMockUnsafeAuthServiceServer(ctrl *gomock.Controller) *MockUnsafeAuthServiceServer {
	mock := &MockUnsafeAuthServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthServiceServer) EXPECT() *MockUnsafeAuthServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthServiceServer mocks base method.
func (m *MockUnsafeAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServiceServer")
}

// mustEmbedUnimplementedAuthServiceServer indicates an expected call of mustEmbedUnimplementedAuthServiceServer.
func (mr *MockUnsafeAuthServiceServerMockRecorder) mustEmbedUnimplementedAuthServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServiceServer", reflect.TypeOf((*MockUnsafeAuthServiceServer)(nil).mustEmbedUnimplementedAuthServiceServer))
}
