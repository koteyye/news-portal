// Code generated by MockGen. DO NOT EDIT.
// Source: proto/user_grpc.pb.go

// Package mock_proto is a generated GoMock package.
package mock_proto

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	proto "github.com/koteyye/news-portal/proto"
	grpc "google.golang.org/grpc"
)

// MockUserClient is a mock of UserClient interface.
type MockUserClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserClientMockRecorder
}

// MockUserClientMockRecorder is the mock recorder for MockUserClient.
type MockUserClientMockRecorder struct {
	mock *MockUserClient
}

// NewMockUserClient creates a new mock instance.
func NewMockUserClient(ctrl *gomock.Controller) *MockUserClient {
	mock := &MockUserClient{ctrl: ctrl}
	mock.recorder = &MockUserClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserClient) EXPECT() *MockUserClientMockRecorder {
	return m.recorder
}

// GetUserByIDs mocks base method.
func (m *MockUserClient) GetUserByIDs(ctx context.Context, in *proto.UserByIDsRequest, opts ...grpc.CallOption) (*proto.UserByIDsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserByIDs", varargs...)
	ret0, _ := ret[0].(*proto.UserByIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByIDs indicates an expected call of GetUserByIDs.
func (mr *MockUserClientMockRecorder) GetUserByIDs(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByIDs", reflect.TypeOf((*MockUserClient)(nil).GetUserByIDs), varargs...)
}

// MockUserServer is a mock of UserServer interface.
type MockUserServer struct {
	ctrl     *gomock.Controller
	recorder *MockUserServerMockRecorder
}

// MockUserServerMockRecorder is the mock recorder for MockUserServer.
type MockUserServerMockRecorder struct {
	mock *MockUserServer
}

// NewMockUserServer creates a new mock instance.
func NewMockUserServer(ctrl *gomock.Controller) *MockUserServer {
	mock := &MockUserServer{ctrl: ctrl}
	mock.recorder = &MockUserServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServer) EXPECT() *MockUserServerMockRecorder {
	return m.recorder
}

// GetUserByIDs mocks base method.
func (m *MockUserServer) GetUserByIDs(arg0 context.Context, arg1 *proto.UserByIDsRequest) (*proto.UserByIDsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByIDs", arg0, arg1)
	ret0, _ := ret[0].(*proto.UserByIDsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByIDs indicates an expected call of GetUserByIDs.
func (mr *MockUserServerMockRecorder) GetUserByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByIDs", reflect.TypeOf((*MockUserServer)(nil).GetUserByIDs), arg0, arg1)
}

// mustEmbedUnimplementedUserServer mocks base method.
func (m *MockUserServer) mustEmbedUnimplementedUserServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedUserServer")
}

// mustEmbedUnimplementedUserServer indicates an expected call of mustEmbedUnimplementedUserServer.
func (mr *MockUserServerMockRecorder) mustEmbedUnimplementedUserServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedUserServer", reflect.TypeOf((*MockUserServer)(nil).mustEmbedUnimplementedUserServer))
}

// MockUnsafeUserServer is a mock of UnsafeUserServer interface.
type MockUnsafeUserServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeUserServerMockRecorder
}

// MockUnsafeUserServerMockRecorder is the mock recorder for MockUnsafeUserServer.
type MockUnsafeUserServerMockRecorder struct {
	mock *MockUnsafeUserServer
}

// NewMockUnsafeUserServer creates a new mock instance.
func NewMockUnsafeUserServer(ctrl *gomock.Controller) *MockUnsafeUserServer {
	mock := &MockUnsafeUserServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeUserServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeUserServer) EXPECT() *MockUnsafeUserServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedUserServer mocks base method.
func (m *MockUnsafeUserServer) mustEmbedUnimplementedUserServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedUserServer")
}

// mustEmbedUnimplementedUserServer indicates an expected call of mustEmbedUnimplementedUserServer.
func (mr *MockUnsafeUserServerMockRecorder) mustEmbedUnimplementedUserServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedUserServer", reflect.TypeOf((*MockUnsafeUserServer)(nil).mustEmbedUnimplementedUserServer))
}
