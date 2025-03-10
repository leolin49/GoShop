// Code generated by MockGen. DO NOT EDIT.
// Source: redis.go

// Package mock_redis is a generated GoMock package.
package redis

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	proto "google.golang.org/protobuf/proto"
)

// MockIRdb is a mock of IRdb interface.
type MockIRdb struct {
	ctrl     *gomock.Controller
	recorder *MockIRdbMockRecorder
}

// MockIRdbMockRecorder is the mock recorder for MockIRdb.
type MockIRdbMockRecorder struct {
	mock *MockIRdb
}

// NewMockIRdb creates a new mock instance.
func NewMockIRdb(ctrl *gomock.Controller) *MockIRdb {
	mock := &MockIRdb{ctrl: ctrl}
	mock.recorder = &MockIRdbMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRdb) EXPECT() *MockIRdbMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockIRdb) Del(k string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", k)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockIRdbMockRecorder) Del(k interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockIRdb)(nil).Del), k)
}

// Exist mocks base method.
func (m *MockIRdb) Exist(k string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exist", k)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exist indicates an expected call of Exist.
func (mr *MockIRdbMockRecorder) Exist(k interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exist", reflect.TypeOf((*MockIRdb)(nil).Exist), k)
}

// Get mocks base method.
func (m *MockIRdb) Get(k string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", k)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIRdbMockRecorder) Get(k interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIRdb)(nil).Get), k)
}

// GetInt mocks base method.
func (m *MockIRdb) GetInt(k string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", k)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInt indicates an expected call of GetInt.
func (mr *MockIRdbMockRecorder) GetInt(k interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockIRdb)(nil).GetInt), k)
}

// GetProto mocks base method.
func (m *MockIRdb) GetProto(k string, v proto.Message) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProto", k, v)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProto indicates an expected call of GetProto.
func (mr *MockIRdbMockRecorder) GetProto(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProto", reflect.TypeOf((*MockIRdb)(nil).GetProto), k, v)
}

// Ping mocks base method.
func (m *MockIRdb) Ping() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockIRdbMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockIRdb)(nil).Ping))
}

// Set mocks base method.
func (m *MockIRdb) Set(k, v string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", k, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockIRdbMockRecorder) Set(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockIRdb)(nil).Set), k, v)
}

// SetInt mocks base method.
func (m *MockIRdb) SetInt(k string, v int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInt", k, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInt indicates an expected call of SetInt.
func (mr *MockIRdbMockRecorder) SetInt(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInt", reflect.TypeOf((*MockIRdb)(nil).SetInt), k, v)
}

// SetProto mocks base method.
func (m *MockIRdb) SetProto(k string, v proto.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetProto", k, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetProto indicates an expected call of SetProto.
func (mr *MockIRdbMockRecorder) SetProto(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetProto", reflect.TypeOf((*MockIRdb)(nil).SetProto), k, v)
}
