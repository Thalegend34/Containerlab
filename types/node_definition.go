package types

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	importEnvsKey = "__IMPORT_ENVS"
)

// NodeDefinition represents a configuration a given node can have in the lab definition file.
type NodeDefinition struct {
	Kind                  string            `yaml:"kind,omitempty"`
	Group                 string            `yaml:"group,omitempty"`
	Type                  string            `yaml:"type,omitempty"`
	StartupConfig         string            `yaml:"startup-config,omitempty"`
	StartupDelay          uint              `yaml:"startup-delay,omitempty"`
	EnforceStartupConfig  *bool             `yaml:"enforce-startup-config,omitempty"`
	SuppressStartupConfig *bool             `yaml:"suppress-startup-config,omitempty"`
	AutoRemove            *bool             `yaml:"auto-remove,omitempty"`
	Config                *ConfigDispatcher `yaml:"config,omitempty"`
	Image                 string            `yaml:"image,omitempty"`
	ImagePullPolicy       string            `yaml:"image-pull-policy,omitempty"`
	License               string            `yaml:"license,omitempty"`
	Position              string            `yaml:"position,omitempty"`
	Entrypoint            string            `yaml:"entrypoint,omitempty"`
	Cmd                   string            `yaml:"cmd,omitempty"`
	// list of commands to run in container
	Exec []string `yaml:"exec,omitempty"`
	// list of bind mount compatible strings
	Binds []string `yaml:"binds,omitempty"`
	// list of port bindings
	Ports []string `yaml:"ports,omitempty"`
	// user-defined IPv4 address in the management network
	MgmtIPv4 string `yaml:"mgmt-ipv4,omitempty"`
	// user-defined IPv6 address in the management network
	MgmtIPv6 string `yaml:"mgmt-ipv6,omitempty"`
	// list of ports to publish with mysocketctl
	Publish []string `yaml:"publish,omitempty"`
	// environment variables
	Env map[string]string `yaml:"env,omitempty"`
	// external file containing environment variables
	EnvFiles []string `yaml:"env-files,omitempty"`
	// linux user used in a container
	User string `yaml:"user,omitempty"`
	// container labels
	Labels map[string]string `yaml:"labels,omitempty"`
	// container networking mode. if set to `host` the host networking will be used for this node, else bridged network
	NetworkMode string `yaml:"network-mode,omitempty"`
	// Ignite sandbox and kernel imageNames
	Sandbox string `yaml:"sandbox,omitempty"`
	Kernel  string `yaml:"kernel,omitempty"`
	// Override container runtime
	Runtime string `yaml:"runtime,omitempty"`
	// Set node CPU (cgroup or hypervisor)
	CPU float64 `yaml:"cpu,omitempty"`
	// Set node CPUs to use
	CPUSet string `yaml:"cpu-set,omitempty"`
	// Set node Memory (cgroup or hypervisor)
	Memory string `yaml:"memory,omitempty"`
	// Set the nodes Sysctl
	Sysctls map[string]string `yaml:"sysctls,omitempty"`
	// Extra options, may be kind specific
	Extras *Extras `yaml:"extras,omitempty"`
	// Deployment stages
	Stages *Stages `yaml:"stages,omitempty"`
	// DNS configuration
	DNS *DNSConfig `yaml:"dns,omitempty"`
	// Certificate configuration
	Certificate *CertificateConfig `yaml:"certificate,omitempty"`
	// Healthcheck configuration
	HealthCheck *HealthcheckConfig `yaml:"healthcheck,omitempty"`
}

// Stages represents a configuration of a given node deployment stage.
type Stages struct {
	Create      *StageCreate      `yaml:"create"`
	CreateLinks *StageCreateLinks `yaml:"create-links"`
	Configure   *StageConfigure   `yaml:"configure"`
	Healthy     *StageHealthy     `yaml:"healthy"`
	Exit        *StageExit        `yaml:"exit"`
}

func NewStages() *Stages {
	return &Stages{
		Create: &StageCreate{
			StageBase: StageBase{},
		},
		CreateLinks: &StageCreateLinks{
			StageBase: StageBase{},
		},
		Configure: &StageConfigure{
			StageBase: StageBase{},
		},
		Healthy: &StageHealthy{
			StageBase: StageBase{},
		},
		Exit: &StageExit{
			StageBase: StageBase{},
		},
	}
}

// GetWaitFor returns lists of nodes that need to be waited for in a map
// that is indexed by the state for which this dependency is to be evaluated
func (s *Stages) GetWaitFor() map[WaitForPhase]WaitForList {
	result := map[WaitForPhase]WaitForList{}

	result[WaitForConfigure] = s.Configure.WaitFor
	result[WaitForCreate] = s.Create.WaitFor
	result[WaitForCreateLinks] = s.CreateLinks.WaitFor
	result[WaitForHealthy] = s.Healthy.WaitFor
	result[WaitForExit] = s.Exit.WaitFor

	return result
}

// Merge merges the Stages configuration of other into s
func (s *Stages) Merge(other *Stages) error {
	var err error
	if other.Configure != nil {
		err = s.Configure.Merge(&other.Configure.StageBase)
		if err != nil {
			return err
		}
	}
	if other.Create != nil {
		err = s.Create.Merge(&other.Create.StageBase)
		if err != nil {
			return err
		}
	}
	if other.CreateLinks != nil {
		err = s.CreateLinks.Merge(&other.CreateLinks.StageBase)
		if err != nil {
			return err
		}
	}
	if other.Healthy != nil {
		err = s.Healthy.Merge(&other.Healthy.StageBase)
		if err != nil {
			return err
		}
	}
	if other.Exit != nil {
		err = s.Exit.Merge(&other.Exit.StageBase)
		if err != nil {
			return err
		}
	}
	return err
}

// StageCreate represents a creation stage of a given node.
type StageCreate struct {
	StageBase `yaml:",inline"`
}

// StageCreateLinks represents a stage of a given node when links are getting added to it.
type StageCreateLinks struct {
	StageBase `yaml:",inline"`
}

// StageConfigure represents a stage of a given node when it enters configuration workflow.
type StageConfigure struct {
	StageBase `yaml:",inline"`
}

// StageHealthy represents a stage of a given node when it reaches healthy status.
type StageHealthy struct {
	StageBase `yaml:",inline"`
}

// StageExit represents a stage of a given node when the node reaches exit state.
type StageExit struct {
	StageBase `yaml:",inline"`
}

// StageBase represents a configuration of a given stage.
type StageBase struct {
	WaitFor WaitForList `yaml:"wait-for,omitempty"`
}

type WaitForList []*WaitFor

func (w WaitForList) contains(newWf *WaitFor) bool {
	for _, entry := range w {
		if entry.Equals(newWf) {
			return true
		}
	}
	return false
}

func (s *StageBase) Merge(sc *StageBase) error {
	if sc == nil {
		return nil
	}
	for _, wf := range sc.WaitFor {
		// prevent adding the same dependency twice
		if s.WaitFor.contains(wf) {
			continue
		}
		s.WaitFor = append(s.WaitFor, wf)
	}
	return nil
}

// Interface compliance.
var _ yaml.Unmarshaler = &NodeDefinition{}

// UnmarshalYAML is a custom unmarshaller for NodeDefinition type that allows to map old attributes to new ones.
func (n *NodeDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// define an alias type to avoid recursion during unmarshalling
	type NodeDefinitionAlias NodeDefinition

	type NodeDefinitionWithDeprecatedFields struct {
		NodeDefinitionAlias `yaml:",inline"`
		DeprecatedMgmtIPv4  string `yaml:"mgmt_ipv4,omitempty"`
		DeprecatedMgmtIPv6  string `yaml:"mgmt_ipv6,omitempty"`
	}

	nd := &NodeDefinitionWithDeprecatedFields{}

	nd.NodeDefinitionAlias = (NodeDefinitionAlias)(*n)
	if err := unmarshal(nd); err != nil {
		return err
	}

	// process deprecated fields and use their values for new fields if new fields are not set
	if len(nd.DeprecatedMgmtIPv4) > 0 && len(nd.MgmtIPv4) == 0 {
		log.Warnf("Attribute \"mgmt_ipv4\" is deprecated and will be removed in future. Change it to \"mgmt-ipv4\"")
		nd.MgmtIPv4 = nd.DeprecatedMgmtIPv4
	}

	if len(nd.DeprecatedMgmtIPv6) > 0 && len(nd.MgmtIPv6) == 0 {
		log.Warnf("Attribute \"mgmt_ipv6\" is deprecated and will be removed in future. Change it to \"mgmt-ipv6\"")
		nd.MgmtIPv6 = nd.DeprecatedMgmtIPv6
	}

	*n = (NodeDefinition)(nd.NodeDefinitionAlias)

	return nil
}

func (n *NodeDefinition) GetKind() string {
	if n == nil {
		return ""
	}
	return n.Kind
}

func (n *NodeDefinition) GetGroup() string {
	if n == nil {
		return ""
	}
	return n.Group
}

func (n *NodeDefinition) GetType() string {
	if n == nil {
		return ""
	}
	return n.Type
}

func (n *NodeDefinition) GetStartupConfig() string {
	if n == nil {
		return ""
	}
	return n.StartupConfig
}

func (n *NodeDefinition) GetStartupDelay() uint {
	if n == nil {
		return 0
	}
	return n.StartupDelay
}

func (n *NodeDefinition) GetEnforceStartupConfig() *bool {
	if n == nil {
		return nil
	}
	return n.EnforceStartupConfig
}

func (n *NodeDefinition) GetSuppressStartupConfig() *bool {
	if n == nil {
		return nil
	}
	return n.SuppressStartupConfig
}

func (n *NodeDefinition) GetAutoRemove() *bool {
	if n == nil {
		return nil
	}
	return n.AutoRemove
}

func (n *NodeDefinition) GetConfigDispatcher() *ConfigDispatcher {
	if n == nil {
		return nil
	}
	return n.Config
}

func (n *NodeDefinition) GetImage() string {
	if n == nil {
		return ""
	}
	return n.Image
}

func (n *NodeDefinition) GetImagePullPolicy() string {
	if n == nil {
		return ""
	}
	return n.ImagePullPolicy
}

func (n *NodeDefinition) GetLicense() string {
	if n == nil {
		return ""
	}
	return n.License
}

func (n *NodeDefinition) GetPostion() string {
	if n == nil {
		return ""
	}
	return n.Position
}

func (n *NodeDefinition) GetEntrypoint() string {
	if n == nil {
		return ""
	}
	return n.Entrypoint
}

func (n *NodeDefinition) GetCmd() string {
	if n == nil {
		return ""
	}
	return n.Cmd
}

func (n *NodeDefinition) GetBinds() []string {
	if n == nil {
		return nil
	}
	return n.Binds
}

func (n *NodeDefinition) GetPorts() []string {
	if n == nil {
		return nil
	}
	return n.Ports
}

func (n *NodeDefinition) GetMgmtIPv4() string {
	if n == nil {
		return ""
	}
	return n.MgmtIPv4
}

func (n *NodeDefinition) GetMgmtIPv6() string {
	if n == nil {
		return ""
	}
	return n.MgmtIPv6
}

func (n *NodeDefinition) GetPublish() []string {
	if n == nil {
		return nil
	}
	return n.Publish
}

func (n *NodeDefinition) GetEnv() map[string]string {
	if n == nil {
		return nil
	}
	return n.Env
}

func (n *NodeDefinition) GetEnvFiles() []string {
	if n == nil {
		return nil
	}
	return n.EnvFiles
}

func (n *NodeDefinition) GetUser() string {
	if n == nil {
		return ""
	}
	return n.User
}

func (n *NodeDefinition) GetLabels() map[string]string {
	if n == nil {
		return nil
	}
	return n.Labels
}

func (n *NodeDefinition) GetNetworkMode() string {
	if n == nil {
		return ""
	}
	return n.NetworkMode
}

func (n *NodeDefinition) GetNodeSandbox() string {
	if n == nil {
		return ""
	}
	return n.Sandbox
}

func (n *NodeDefinition) GetNodeKernel() string {
	if n == nil {
		return ""
	}
	return n.Kernel
}

func (n *NodeDefinition) GetNodeRuntime() string {
	if n == nil {
		return ""
	}
	return n.Runtime
}

func (n *NodeDefinition) GetNodeCPU() float64 {
	if n == nil {
		return 0
	}
	return n.CPU
}

func (n *NodeDefinition) GetNodeCPUSet() string {
	if n == nil {
		return ""
	}
	return n.CPUSet
}

func (n *NodeDefinition) GetNodeMemory() string {
	if n == nil {
		return ""
	}
	return n.Memory
}

func (n *NodeDefinition) GetExec() []string {
	if n == nil {
		return nil
	}
	return n.Exec
}

func (n *NodeDefinition) GetSysctls() map[string]string {
	if n == nil || n.Sysctls == nil {
		return map[string]string{}
	}

	return n.Sysctls
}

func (n *NodeDefinition) GetExtras() *Extras {
	if n == nil {
		return nil
	}
	return n.Extras
}

func (n *NodeDefinition) GetStages() *Stages {
	if n == nil {
		return nil
	}
	return n.Stages
}

func (n *NodeDefinition) GetDns() *DNSConfig {
	if n == nil {
		return nil
	}
	return n.DNS
}

func (n *NodeDefinition) GetCertificateConfig() *CertificateConfig {
	if n == nil {
		return nil
	}
	return n.Certificate
}

func (n *NodeDefinition) GetHealthcheckConfig() *HealthcheckConfig {
	if n == nil {
		return nil
	}
	return n.HealthCheck
}

// ImportEnvs imports all environment variales defined in the shell
// if __IMPORT_ENVS is set to true.
func (n *NodeDefinition) ImportEnvs() {
	if n == nil || n.Env == nil {
		return
	}

	var importEnvs bool

	for k, v := range n.Env {
		if k == importEnvsKey && v == "true" {
			importEnvs = true
			break
		}
	}

	if !importEnvs {
		return
	}

	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		if _, exists := n.Env[kv[0]]; exists {
			continue
		}

		n.Env[kv[0]] = kv[1]
	}
}
