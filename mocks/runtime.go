// Code generated by MockGen. DO NOT EDIT.
// Source: runtime/runtime.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	runtime "github.com/srl-labs/containerlab/runtime"
	types "github.com/srl-labs/containerlab/types"
)

// MockContainerRuntime is a mock of ContainerRuntime interface.
type MockContainerRuntime struct {
	ctrl     *gomock.Controller
	recorder *MockContainerRuntimeMockRecorder
}

// MockContainerRuntimeMockRecorder is the mock recorder for MockContainerRuntime.
type MockContainerRuntimeMockRecorder struct {
	mock *MockContainerRuntime
}

// NewMockContainerRuntime creates a new mock instance.
func NewMockContainerRuntime(ctrl *gomock.Controller) *MockContainerRuntime {
	mock := &MockContainerRuntime{ctrl: ctrl}
	mock.recorder = &MockContainerRuntimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerRuntime) EXPECT() *MockContainerRuntimeMockRecorder {
	return m.recorder
}

// Config mocks base method.
func (m *MockContainerRuntime) Config() runtime.RuntimeConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(runtime.RuntimeConfig)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockContainerRuntimeMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockContainerRuntime)(nil).Config))
}

// CreateContainer mocks base method.
func (m *MockContainerRuntime) CreateContainer(arg0 context.Context, arg1 *types.NodeConfig) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContainer", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContainer indicates an expected call of CreateContainer.
func (mr *MockContainerRuntimeMockRecorder) CreateContainer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContainer", reflect.TypeOf((*MockContainerRuntime)(nil).CreateContainer), arg0, arg1)
}

// CreateNet mocks base method.
func (m *MockContainerRuntime) CreateNet(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNet", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNet indicates an expected call of CreateNet.
func (mr *MockContainerRuntimeMockRecorder) CreateNet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNet", reflect.TypeOf((*MockContainerRuntime)(nil).CreateNet), arg0)
}

// DeleteContainer mocks base method.
func (m *MockContainerRuntime) DeleteContainer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContainer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContainer indicates an expected call of DeleteContainer.
func (mr *MockContainerRuntimeMockRecorder) DeleteContainer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContainer", reflect.TypeOf((*MockContainerRuntime)(nil).DeleteContainer), arg0, arg1)
}

// DeleteNet mocks base method.
func (m *MockContainerRuntime) DeleteNet(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNet", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNet indicates an expected call of DeleteNet.
func (mr *MockContainerRuntimeMockRecorder) DeleteNet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNet", reflect.TypeOf((*MockContainerRuntime)(nil).DeleteNet), arg0)
}

// Exec mocks base method.
func (m *MockContainerRuntime) Exec(ctx context.Context, cID string, exec types.ExecExecutor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", ctx, cID, exec)
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec.
func (mr *MockContainerRuntimeMockRecorder) Exec(ctx, cID, exec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockContainerRuntime)(nil).Exec), ctx, cID, exec)
}

// ExecNotWait mocks base method.
func (m *MockContainerRuntime) ExecNotWait(ctx context.Context, cID string, exec types.ExecExecutor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecNotWait", ctx, cID, exec)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecNotWait indicates an expected call of ExecNotWait.
func (mr *MockContainerRuntimeMockRecorder) ExecNotWait(ctx, cID, exec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecNotWait", reflect.TypeOf((*MockContainerRuntime)(nil).ExecNotWait), ctx, cID, exec)
}

// GetContainerStatus mocks base method.
func (m *MockContainerRuntime) GetContainerStatus(ctx context.Context, cID string) runtime.ContainerStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContainerStatus", ctx, cID)
	ret0, _ := ret[0].(runtime.ContainerStatus)
	return ret0
}

// GetContainerStatus indicates an expected call of GetContainerStatus.
func (mr *MockContainerRuntimeMockRecorder) GetContainerStatus(ctx, cID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContainerStatus", reflect.TypeOf((*MockContainerRuntime)(nil).GetContainerStatus), ctx, cID)
}

// GetHostsPath mocks base method.
func (m *MockContainerRuntime) GetHostsPath(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostsPath", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHostsPath indicates an expected call of GetHostsPath.
func (mr *MockContainerRuntimeMockRecorder) GetHostsPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostsPath", reflect.TypeOf((*MockContainerRuntime)(nil).GetHostsPath), arg0, arg1)
}

// GetNSPath mocks base method.
func (m *MockContainerRuntime) GetNSPath(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNSPath", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNSPath indicates an expected call of GetNSPath.
func (mr *MockContainerRuntimeMockRecorder) GetNSPath(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNSPath", reflect.TypeOf((*MockContainerRuntime)(nil).GetNSPath), arg0, arg1)
}

// GetName mocks base method.
func (m *MockContainerRuntime) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockContainerRuntimeMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockContainerRuntime)(nil).GetName))
}

// Init mocks base method.
func (m *MockContainerRuntime) Init(arg0 ...runtime.RuntimeOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Init", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockContainerRuntimeMockRecorder) Init(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockContainerRuntime)(nil).Init), arg0...)
}

// ListContainers mocks base method.
func (m *MockContainerRuntime) ListContainers(arg0 context.Context, arg1 []*types.GenericFilter) ([]types.GenericContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListContainers", arg0, arg1)
	ret0, _ := ret[0].([]types.GenericContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContainers indicates an expected call of ListContainers.
func (mr *MockContainerRuntimeMockRecorder) ListContainers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContainers", reflect.TypeOf((*MockContainerRuntime)(nil).ListContainers), arg0, arg1)
}

// Mgmt mocks base method.
func (m *MockContainerRuntime) Mgmt() *types.MgmtNet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mgmt")
	ret0, _ := ret[0].(*types.MgmtNet)
	return ret0
}

// Mgmt indicates an expected call of Mgmt.
func (mr *MockContainerRuntimeMockRecorder) Mgmt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mgmt", reflect.TypeOf((*MockContainerRuntime)(nil).Mgmt))
}

// PauseContainer mocks base method.
func (m *MockContainerRuntime) PauseContainer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseContainer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PauseContainer indicates an expected call of PauseContainer.
func (mr *MockContainerRuntimeMockRecorder) PauseContainer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseContainer", reflect.TypeOf((*MockContainerRuntime)(nil).PauseContainer), arg0, arg1)
}

// PullImageIfRequired mocks base method.
func (m *MockContainerRuntime) PullImageIfRequired(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PullImageIfRequired", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullImageIfRequired indicates an expected call of PullImageIfRequired.
func (mr *MockContainerRuntimeMockRecorder) PullImageIfRequired(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PullImageIfRequired", reflect.TypeOf((*MockContainerRuntime)(nil).PullImageIfRequired), arg0, arg1)
}

// StartContainer mocks base method.
func (m *MockContainerRuntime) StartContainer(arg0 context.Context, arg1 string, arg2 *types.NodeConfig) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartContainer", arg0, arg1, arg2)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartContainer indicates an expected call of StartContainer.
func (mr *MockContainerRuntimeMockRecorder) StartContainer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartContainer", reflect.TypeOf((*MockContainerRuntime)(nil).StartContainer), arg0, arg1, arg2)
}

// StopContainer mocks base method.
func (m *MockContainerRuntime) StopContainer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopContainer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopContainer indicates an expected call of StopContainer.
func (mr *MockContainerRuntimeMockRecorder) StopContainer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopContainer", reflect.TypeOf((*MockContainerRuntime)(nil).StopContainer), arg0, arg1)
}

// UnpauseContainer mocks base method.
func (m *MockContainerRuntime) UnpauseContainer(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpauseContainer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnpauseContainer indicates an expected call of UnpauseContainer.
func (mr *MockContainerRuntimeMockRecorder) UnpauseContainer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpauseContainer", reflect.TypeOf((*MockContainerRuntime)(nil).UnpauseContainer), arg0, arg1)
}

// WithConfig mocks base method.
func (m *MockContainerRuntime) WithConfig(arg0 *runtime.RuntimeConfig) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithConfig", arg0)
}

// WithConfig indicates an expected call of WithConfig.
func (mr *MockContainerRuntimeMockRecorder) WithConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithConfig", reflect.TypeOf((*MockContainerRuntime)(nil).WithConfig), arg0)
}

// WithKeepMgmtNet mocks base method.
func (m *MockContainerRuntime) WithKeepMgmtNet() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithKeepMgmtNet")
}

// WithKeepMgmtNet indicates an expected call of WithKeepMgmtNet.
func (mr *MockContainerRuntimeMockRecorder) WithKeepMgmtNet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithKeepMgmtNet", reflect.TypeOf((*MockContainerRuntime)(nil).WithKeepMgmtNet))
}

// WithMgmtNet mocks base method.
func (m *MockContainerRuntime) WithMgmtNet(arg0 *types.MgmtNet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithMgmtNet", arg0)
}

// WithMgmtNet indicates an expected call of WithMgmtNet.
func (mr *MockContainerRuntimeMockRecorder) WithMgmtNet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithMgmtNet", reflect.TypeOf((*MockContainerRuntime)(nil).WithMgmtNet), arg0)
}
