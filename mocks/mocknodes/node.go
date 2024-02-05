// Code generated by MockGen. DO NOT EDIT.
// Source: nodes/node.go
//
// Generated by this command:
//
//	mockgen -package=mocknodes -source=nodes/node.go -destination=./mocks/mocknodes/node.go
//

// Package mocknodes is a generated GoMock package.
package mocknodes

import (
	context "context"
	reflect "reflect"

	ns "github.com/containernetworking/plugins/pkg/ns"
	exec "github.com/srl-labs/containerlab/clab/exec"
	links "github.com/srl-labs/containerlab/links"
	nodes "github.com/srl-labs/containerlab/nodes"
	state "github.com/srl-labs/containerlab/nodes/state"
	runtime "github.com/srl-labs/containerlab/runtime"
	types "github.com/srl-labs/containerlab/types"
	netlink "github.com/vishvananda/netlink"
	gomock "go.uber.org/mock/gomock"
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

// AddEndpoint mocks base method.
func (m *MockNode) AddEndpoint(e links.Endpoint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddEndpoint", e)
}

// AddEndpoint indicates an expected call of AddEndpoint.
func (mr *MockNodeMockRecorder) AddEndpoint(e any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEndpoint", reflect.TypeOf((*MockNode)(nil).AddEndpoint), e)
}

// AddLink mocks base method.
func (m *MockNode) AddLink(l links.Link) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddLink", l)
}

// AddLink indicates an expected call of AddLink.
func (mr *MockNodeMockRecorder) AddLink(l any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLink", reflect.TypeOf((*MockNode)(nil).AddLink), l)
}

// AddLinkToContainer mocks base method.
func (m *MockNode) AddLinkToContainer(ctx context.Context, link netlink.Link, f func(ns.NetNS) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLinkToContainer", ctx, link, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLinkToContainer indicates an expected call of AddLinkToContainer.
func (mr *MockNodeMockRecorder) AddLinkToContainer(ctx, link, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLinkToContainer", reflect.TypeOf((*MockNode)(nil).AddLinkToContainer), ctx, link, f)
}

// CheckDeploymentConditions mocks base method.
func (m *MockNode) CheckDeploymentConditions(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDeploymentConditions", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckDeploymentConditions indicates an expected call of CheckDeploymentConditions.
func (mr *MockNodeMockRecorder) CheckDeploymentConditions(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDeploymentConditions", reflect.TypeOf((*MockNode)(nil).CheckDeploymentConditions), arg0)
}

// CheckInterfaceName mocks base method.
func (m *MockNode) CheckInterfaceName() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInterfaceName")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckInterfaceName indicates an expected call of CheckInterfaceName.
func (mr *MockNodeMockRecorder) CheckInterfaceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInterfaceName", reflect.TypeOf((*MockNode)(nil).CheckInterfaceName))
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
func (mr *MockNodeMockRecorder) Delete(arg0 any) *gomock.Call {
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
func (m *MockNode) Deploy(arg0 context.Context, arg1 *nodes.DeployParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockNodeMockRecorder) Deploy(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockNode)(nil).Deploy), arg0, arg1)
}

// DeployLinks mocks base method.
func (m *MockNode) DeployLinks(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployLinks", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployLinks indicates an expected call of DeployLinks.
func (mr *MockNodeMockRecorder) DeployLinks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployLinks", reflect.TypeOf((*MockNode)(nil).DeployLinks), ctx)
}

// ExecFunction mocks base method.
func (m *MockNode) ExecFunction(arg0 func(ns.NetNS) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecFunction", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecFunction indicates an expected call of ExecFunction.
func (mr *MockNodeMockRecorder) ExecFunction(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecFunction", reflect.TypeOf((*MockNode)(nil).ExecFunction), arg0)
}

// GenerateConfig mocks base method.
func (m *MockNode) GenerateConfig(dst, templ string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateConfig", dst, templ)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateConfig indicates an expected call of GenerateConfig.
func (mr *MockNodeMockRecorder) GenerateConfig(dst, templ any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateConfig", reflect.TypeOf((*MockNode)(nil).GenerateConfig), dst, templ)
}

// GetContainers mocks base method.
func (m *MockNode) GetContainers(ctx context.Context) ([]runtime.GenericContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContainers", ctx)
	ret0, _ := ret[0].([]runtime.GenericContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContainers indicates an expected call of GetContainers.
func (mr *MockNodeMockRecorder) GetContainers(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContainers", reflect.TypeOf((*MockNode)(nil).GetContainers), ctx)
}

// GetEndpoints mocks base method.
func (m *MockNode) GetEndpoints() []links.Endpoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEndpoints")
	ret0, _ := ret[0].([]links.Endpoint)
	return ret0
}

// GetEndpoints indicates an expected call of GetEndpoints.
func (mr *MockNodeMockRecorder) GetEndpoints() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEndpoints", reflect.TypeOf((*MockNode)(nil).GetEndpoints))
}

// GetImages mocks base method.
func (m *MockNode) GetImages(arg0 context.Context) map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImages", arg0)
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetImages indicates an expected call of GetImages.
func (mr *MockNodeMockRecorder) GetImages(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImages", reflect.TypeOf((*MockNode)(nil).GetImages), arg0)
}

// GetLinkEndpointType mocks base method.
func (m *MockNode) GetLinkEndpointType() links.LinkEndpointType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkEndpointType")
	ret0, _ := ret[0].(links.LinkEndpointType)
	return ret0
}

// GetLinkEndpointType indicates an expected call of GetLinkEndpointType.
func (mr *MockNodeMockRecorder) GetLinkEndpointType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkEndpointType", reflect.TypeOf((*MockNode)(nil).GetLinkEndpointType))
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

// GetSSHConfig mocks base method.
func (m *MockNode) GetSSHConfig() *types.SSHConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSSHConfig")
	ret0, _ := ret[0].(*types.SSHConfig)
	return ret0
}

// GetSSHConfig indicates an expected call of GetSSHConfig.
func (mr *MockNodeMockRecorder) GetSSHConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSSHConfig", reflect.TypeOf((*MockNode)(nil).GetSSHConfig))
}

// GetShortName mocks base method.
func (m *MockNode) GetShortName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetShortName indicates an expected call of GetShortName.
func (mr *MockNodeMockRecorder) GetShortName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortName", reflect.TypeOf((*MockNode)(nil).GetShortName))
}

// GetState mocks base method.
func (m *MockNode) GetState() state.NodeState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState")
	ret0, _ := ret[0].(state.NodeState)
	return ret0
}

// GetState indicates an expected call of GetState.
func (mr *MockNodeMockRecorder) GetState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockNode)(nil).GetState))
}

// Init mocks base method.
func (m *MockNode) Init(arg0 *types.NodeConfig, arg1 ...nodes.NodeOption) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Init", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockNodeMockRecorder) Init(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockNode)(nil).Init), varargs...)
}

// PostDeploy mocks base method.
func (m *MockNode) PostDeploy(ctx context.Context, params *nodes.PostDeployParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostDeploy", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostDeploy indicates an expected call of PostDeploy.
func (mr *MockNodeMockRecorder) PostDeploy(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostDeploy", reflect.TypeOf((*MockNode)(nil).PostDeploy), ctx, params)
}

// PreDeploy mocks base method.
func (m *MockNode) PreDeploy(ctx context.Context, params *nodes.PreDeployParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreDeploy", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// PreDeploy indicates an expected call of PreDeploy.
func (mr *MockNodeMockRecorder) PreDeploy(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreDeploy", reflect.TypeOf((*MockNode)(nil).PreDeploy), ctx, params)
}

// RunExec mocks base method.
func (m *MockNode) RunExec(ctx context.Context, execCmd *exec.ExecCmd) (*exec.ExecResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunExec", ctx, execCmd)
	ret0, _ := ret[0].(*exec.ExecResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunExec indicates an expected call of RunExec.
func (mr *MockNodeMockRecorder) RunExec(ctx, execCmd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunExec", reflect.TypeOf((*MockNode)(nil).RunExec), ctx, execCmd)
}

// RunExecFromConfig mocks base method.
func (m *MockNode) RunExecFromConfig(arg0 context.Context, arg1 *exec.ExecCollection) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunExecFromConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunExecFromConfig indicates an expected call of RunExecFromConfig.
func (mr *MockNodeMockRecorder) RunExecFromConfig(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunExecFromConfig", reflect.TypeOf((*MockNode)(nil).RunExecFromConfig), arg0, arg1)
}

// SaveConfig mocks base method.
func (m *MockNode) SaveConfig(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveConfig indicates an expected call of SaveConfig.
func (mr *MockNodeMockRecorder) SaveConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveConfig", reflect.TypeOf((*MockNode)(nil).SaveConfig), arg0)
}

// SetState mocks base method.
func (m *MockNode) SetState(arg0 state.NodeState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetState", arg0)
}

// SetState indicates an expected call of SetState.
func (mr *MockNodeMockRecorder) SetState(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockNode)(nil).SetState), arg0)
}

// UpdateConfigWithRuntimeInfo mocks base method.
func (m *MockNode) UpdateConfigWithRuntimeInfo(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfigWithRuntimeInfo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfigWithRuntimeInfo indicates an expected call of UpdateConfigWithRuntimeInfo.
func (mr *MockNodeMockRecorder) UpdateConfigWithRuntimeInfo(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfigWithRuntimeInfo", reflect.TypeOf((*MockNode)(nil).UpdateConfigWithRuntimeInfo), arg0)
}

// VerifyStartupConfig mocks base method.
func (m *MockNode) VerifyStartupConfig(topoDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyStartupConfig", topoDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyStartupConfig indicates an expected call of VerifyStartupConfig.
func (mr *MockNodeMockRecorder) VerifyStartupConfig(topoDir any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyStartupConfig", reflect.TypeOf((*MockNode)(nil).VerifyStartupConfig), topoDir)
}

// WaitForAllLinksCreated mocks base method.
func (m *MockNode) WaitForAllLinksCreated() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitForAllLinksCreated")
}

// WaitForAllLinksCreated indicates an expected call of WaitForAllLinksCreated.
func (mr *MockNodeMockRecorder) WaitForAllLinksCreated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForAllLinksCreated", reflect.TypeOf((*MockNode)(nil).WaitForAllLinksCreated))
}

// WithMgmtNet mocks base method.
func (m *MockNode) WithMgmtNet(arg0 *types.MgmtNet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithMgmtNet", arg0)
}

// WithMgmtNet indicates an expected call of WithMgmtNet.
func (mr *MockNodeMockRecorder) WithMgmtNet(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithMgmtNet", reflect.TypeOf((*MockNode)(nil).WithMgmtNet), arg0)
}

// WithRuntime mocks base method.
func (m *MockNode) WithRuntime(arg0 runtime.ContainerRuntime) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithRuntime", arg0)
}

// WithRuntime indicates an expected call of WithRuntime.
func (mr *MockNodeMockRecorder) WithRuntime(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRuntime", reflect.TypeOf((*MockNode)(nil).WithRuntime), arg0)
}
