// Code generated by MockGen. DO NOT EDIT.
// Source: nodes/node.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	nodes "github.com/srl-labs/containerlab/nodes"
	runtime "github.com/srl-labs/containerlab/runtime"
	types "github.com/srl-labs/containerlab/types"
)

// MockNode is a mock of Node interface.
type MockNode struct {
	ctrl     *gomock.Controller
	recorder *MockNodeMockRecorder
}

// MockNodeMockRecorder is the mock recorder for MockNode.
type MockNodeMockRecorder struct {
	mock *MockNode
}

// NewMockNode creates a new mock instance.
func NewMockNode(ctrl *gomock.Controller) *MockNode {
	mock := &MockNode{ctrl: ctrl}
	mock.recorder = &MockNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNode) EXPECT() *MockNodeMockRecorder {
	return m.recorder
}

// CheckInterfaceNamingConvention mocks base method.
func (m *MockNode) CheckInterfaceNamingConvention() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInterfaceNamingConvention")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckInterfaceNamingConvention indicates an expected call of CheckInterfaceNamingConvention.
func (mr *MockNodeMockRecorder) CheckInterfaceNamingConvention() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInterfaceNamingConvention", reflect.TypeOf((*MockNode)(nil).CheckInterfaceNamingConvention))
}

// Config mocks base method.
func (m *MockNode) Config() *types.NodeConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*types.NodeConfig)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockNodeMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockNode)(nil).Config))
}

// Delete mocks base method.
func (m *MockNode) Delete(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockNodeMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNode)(nil).Delete), arg0)
}

// DeleteNetnsSymlink mocks base method.
func (m *MockNode) DeleteNetnsSymlink() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNetnsSymlink")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNetnsSymlink indicates an expected call of DeleteNetnsSymlink.
func (mr *MockNodeMockRecorder) DeleteNetnsSymlink() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNetnsSymlink", reflect.TypeOf((*MockNode)(nil).DeleteNetnsSymlink))
}

// Deploy mocks base method.
func (m *MockNode) Deploy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockNodeMockRecorder) Deploy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockNode)(nil).Deploy), arg0)
}

// GenerateConfig mocks base method.
func (m *MockNode) GenerateConfig(dst, templ string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateConfig", dst, templ)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateConfig indicates an expected call of GenerateConfig.
func (mr *MockNodeMockRecorder) GenerateConfig(dst, templ interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateConfig", reflect.TypeOf((*MockNode)(nil).GenerateConfig), dst, templ)
}

// GetImages mocks base method.
func (m *MockNode) GetImages(arg0 context.Context) map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImages", arg0)
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetImages indicates an expected call of GetImages.
func (mr *MockNodeMockRecorder) GetImages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImages", reflect.TypeOf((*MockNode)(nil).GetImages), arg0)
}

// GetRuntime mocks base method.
func (m *MockNode) GetRuntime() runtime.ContainerRuntime {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRuntime")
	ret0, _ := ret[0].(runtime.ContainerRuntime)
	return ret0
}

// GetRuntime indicates an expected call of GetRuntime.
func (mr *MockNodeMockRecorder) GetRuntime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRuntime", reflect.TypeOf((*MockNode)(nil).GetRuntime))
}

// GetRuntimeInformation mocks base method.
func (m *MockNode) GetRuntimeInformation(ctx context.Context) ([]types.GenericContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRuntimeInformation", ctx)
	ret0, _ := ret[0].([]types.GenericContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRuntimeInformation indicates an expected call of GetRuntimeInformation.
func (mr *MockNodeMockRecorder) GetRuntimeInformation(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRuntimeInformation", reflect.TypeOf((*MockNode)(nil).GetRuntimeInformation), ctx)
}

// Init mocks base method.
func (m *MockNode) Init(arg0 *types.NodeConfig, arg1 ...nodes.NodeOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Init", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockNodeMockRecorder) Init(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockNode)(nil).Init), varargs...)
}

// PostDeploy mocks base method.
func (m *MockNode) PostDeploy(arg0 context.Context, arg1 map[string]nodes.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostDeploy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostDeploy indicates an expected call of PostDeploy.
func (mr *MockNodeMockRecorder) PostDeploy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostDeploy", reflect.TypeOf((*MockNode)(nil).PostDeploy), arg0, arg1)
}

// PreCheckDeploymentConditionsMeet mocks base method.
func (m *MockNode) PreCheckDeploymentConditionsMeet(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreCheckDeploymentConditionsMeet", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PreCheckDeploymentConditionsMeet indicates an expected call of PreCheckDeploymentConditionsMeet.
func (mr *MockNodeMockRecorder) PreCheckDeploymentConditionsMeet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreCheckDeploymentConditionsMeet", reflect.TypeOf((*MockNode)(nil).PreCheckDeploymentConditionsMeet), arg0)
}

// PreDeploy mocks base method.
func (m *MockNode) PreDeploy(ctx context.Context, configName, labCADir, labCARoot string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreDeploy", ctx, configName, labCADir, labCARoot)
	ret0, _ := ret[0].(error)
	return ret0
}

// PreDeploy indicates an expected call of PreDeploy.
func (mr *MockNodeMockRecorder) PreDeploy(ctx, configName, labCADir, labCARoot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreDeploy", reflect.TypeOf((*MockNode)(nil).PreDeploy), ctx, configName, labCADir, labCARoot)
}

// RunExecConfig mocks base method.
func (m *MockNode) RunExecConfig(ctx context.Context) ([]types.ExecReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunExecConfig", ctx)
	ret0, _ := ret[0].([]types.ExecReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunExecConfig indicates an expected call of RunExecConfig.
func (mr *MockNodeMockRecorder) RunExecConfig(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunExecConfig", reflect.TypeOf((*MockNode)(nil).RunExecConfig), ctx)
}

// RunExecType mocks base method.
func (m *MockNode) RunExecType(ctx context.Context, exec *types.Exec) (types.ExecReader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunExecType", ctx, exec)
	ret0, _ := ret[0].(types.ExecReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunExecType indicates an expected call of RunExecType.
func (mr *MockNodeMockRecorder) RunExecType(ctx, exec interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunExecType", reflect.TypeOf((*MockNode)(nil).RunExecType), ctx, exec)
}

// SaveConfig mocks base method.
func (m *MockNode) SaveConfig(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveConfig indicates an expected call of SaveConfig.
func (mr *MockNodeMockRecorder) SaveConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveConfig", reflect.TypeOf((*MockNode)(nil).SaveConfig), arg0)
}

// VerifyStartupConfig mocks base method.
func (m *MockNode) VerifyStartupConfig(topoDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyStartupConfig", topoDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyStartupConfig indicates an expected call of VerifyStartupConfig.
func (mr *MockNodeMockRecorder) VerifyStartupConfig(topoDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyStartupConfig", reflect.TypeOf((*MockNode)(nil).VerifyStartupConfig), topoDir)
}

// WithMgmtNet mocks base method.
func (m *MockNode) WithMgmtNet(arg0 *types.MgmtNet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithMgmtNet", arg0)
}

// WithMgmtNet indicates an expected call of WithMgmtNet.
func (mr *MockNodeMockRecorder) WithMgmtNet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithMgmtNet", reflect.TypeOf((*MockNode)(nil).WithMgmtNet), arg0)
}

// WithRuntime mocks base method.
func (m *MockNode) WithRuntime(arg0 runtime.ContainerRuntime) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithRuntime", arg0)
}

// WithRuntime indicates an expected call of WithRuntime.
func (mr *MockNodeMockRecorder) WithRuntime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRuntime", reflect.TypeOf((*MockNode)(nil).WithRuntime), arg0)
}
