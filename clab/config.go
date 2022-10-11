// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package clab

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/klauspost/cpuid"
	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/nodes"
	clabRuntimes "github.com/srl-labs/containerlab/runtime"
	"github.com/srl-labs/containerlab/types"
	"github.com/srl-labs/containerlab/utils"
	"github.com/vishvananda/netlink"
)

const (
	// prefix is used to distinct containerlab created files/dirs/containers.
	defaultPrefix = "clab"
	// a name of a docker network that nodes management interfaces connect to.
	dockerNetName     = "clab"
	dockerNetIPv4Addr = "172.20.20.0/24"
	dockerNetIPv6Addr = "2001:172:20:20::/64"
	// NSPath value assigned to host interfaces.
	hostNSPath = "__host"
	// veth link mtu.
	DefaultVethLinkMTU = 9500
	// containerlab's reserved OUI.
	ClabOUI = "aa:c1:ab"

	// label names.
	ContainerlabLabel = "containerlab"
	NodeNameLabel     = "clab-node-name"
	NodeKindLabel     = "clab-node-kind"
	NodeTypeLabel     = "clab-node-type"
	NodeGroupLabel    = "clab-node-group"
	NodeLabDirLabel   = "clab-node-lab-dir"
	TopoFileLabel     = "clab-topo-file"
	NodeMgmtNetBr     = "clab-mgmt-net-bridge"

	// clab specific topology variables.
	clabDirVar = "__clabDir__"
	nodeDirVar = "__clabNodeDir__"
)

// Config defines lab configuration as it is provided in the YAML file.
type Config struct {
	Name     string          `json:"name,omitempty"`
	Prefix   *string         `json:"prefix,omitempty"`
	Mgmt     *types.MgmtNet  `json:"mgmt,omitempty"`
	Topology *types.Topology `json:"topology,omitempty"`
}

// ParseTopology parses the lab topology.
func (c *CLab) parseTopology() error {
	log.Infof("Parsing & checking topology file: %s", c.TopoFile.fullName)

	cwd, _ := filepath.Abs(os.Getenv("PWD"))

	if c.Config.Prefix == nil {
		c.Config.Prefix = new(string)
		*c.Config.Prefix = defaultPrefix
	}

	c.Dir = &Directory{}
	// labDir is always named clab-$labName, regardless of the prefix
	labDir := strings.Join([]string{"clab", c.Config.Name}, "-")
	c.Dir.Lab = filepath.Join(cwd, labDir)

	c.Dir.LabCA = filepath.Join(c.Dir.Lab, "ca")
	c.Dir.LabCARoot = filepath.Join(c.Dir.LabCA, "root")
	c.Dir.LabGraph = filepath.Join(c.Dir.Lab, "graph")

	// initialize Nodes and Links variable
	c.Nodes = make(map[string]nodes.Node)
	c.Links = make(map[int]*types.Link)

	// initialize the Node information from the topology map
	nodeNames := make([]string, 0, len(c.Config.Topology.Nodes))
	for nodeName := range c.Config.Topology.Nodes {
		nodeNames = append(nodeNames, nodeName)
	}
	sort.Strings(nodeNames)

	// collect node runtimes in a map[NodeName] -> RuntimeName
	nodeRuntimes := make(map[string]string)

	for nodeName, topologyNode := range c.Config.Topology.Nodes {
		// this case is when runtime was overridden at the node level
		if r := c.Config.Topology.GetNodeRuntime(nodeName); r != "" {
			nodeRuntimes[nodeName] = r
			continue
		}

		// this case if for non-default runtimes overriding the global default
		if r, ok := nodes.NonDefaultRuntimes[topologyNode.GetKind()]; ok {
			nodeRuntimes[nodeName] = r
			continue
		}

		// saving the global default runtime
		nodeRuntimes[nodeName] = c.globalRuntime
	}

	// initialize any extra runtimes
	for _, r := range nodeRuntimes {
		// this is the case for already init'ed runtimes
		if _, ok := c.Runtimes[r]; ok {
			continue
		}

		if rInit, ok := clabRuntimes.ContainerRuntimes[r]; ok {

			newRuntime := rInit()
			defaultConfig := c.Runtimes[c.globalRuntime].Config()
			err := newRuntime.Init(
				clabRuntimes.WithConfig(&defaultConfig),
			)
			if err != nil {
				return fmt.Errorf("failed to init the container runtime: %s", err)
			}

			c.Runtimes[r] = newRuntime
		}
	}

	var err error
	for idx, nodeName := range nodeNames {
		err = c.NewNode(nodeName, nodeRuntimes[nodeName], c.Config.Topology.Nodes[nodeName], idx)
		if err != nil {
			return err
		}
	}
	for i, l := range c.Config.Topology.Links {
		// i represents the endpoint integer and l provide the link struct
		c.Links[i] = c.NewLink(l)
	}

	// set any containerlab defaults after we've parsed the input
	c.setDefaults()

	return nil
}

// NewNode initializes a new node object.
func (c *CLab) NewNode(nodeName, nodeRuntime string, nodeDef *types.NodeDefinition, idx int) error {
	nodeCfg, err := c.createNodeCfg(nodeName, nodeDef, idx)
	if err != nil {
		return err
	}

	// Init
	nodeInitializer, ok := nodes.Nodes[nodeCfg.Kind]
	if !ok {
		// collect all kind names registered by the nodes in a slice
		kinds := make([]string, 0, len(nodes.Nodes))
		for k := range nodes.Nodes {
			kinds = append(kinds, k)
		}
		return fmt.Errorf("node %q refers to a kind %q which is not supported. Supported kinds are %q", nodeCfg.ShortName, nodeCfg.Kind, kinds)
	}
	n := nodeInitializer()
	// Init

	err = n.Init(nodeCfg, nodes.WithRuntime(c.Runtimes[nodeRuntime]), nodes.WithMgmtNet(c.Config.Mgmt))
	if err != nil {
		log.Errorf("failed to initialize node %q: %v", nodeCfg.ShortName, err)
		return fmt.Errorf("failed to initialize node %q: %v", nodeCfg.ShortName, err)
	}

	c.Nodes[nodeName] = n

	c.addDefaultLabels(n)

	labelsToEnvVars(n.Config())

	return nil
}

func (c *CLab) createNodeCfg(nodeName string, nodeDef *types.NodeDefinition, idx int) (*types.NodeConfig, error) {
	// default longName follows $prefix-$lab-$nodeName pattern
	longName := fmt.Sprintf("%s-%s-%s", *c.Config.Prefix, c.Config.Name, nodeName)

	switch {
	// when prefix is an empty string longName will match shortName/nodeName
	case *c.Config.Prefix == "":
		longName = nodeName
	case *c.Config.Prefix == "__lab-name":
		longName = fmt.Sprintf("%s-%s", c.Config.Name, nodeName)
	}

	nodeCfg := &types.NodeConfig{
		ShortName:            nodeName, // just the node name as seen in the topo file
		LongName:             longName, // by default clab-$labName-$nodeName
		Fqdn:                 strings.Join([]string{nodeName, c.Config.Name, "io"}, "."),
		LabDir:               filepath.Join(c.Dir.Lab, nodeName),
		Index:                idx,
		Group:                c.Config.Topology.GetNodeGroup(nodeName),
		Kind:                 strings.ToLower(c.Config.Topology.GetNodeKind(nodeName)),
		NodeType:             c.Config.Topology.GetNodeType(nodeName),
		Position:             c.Config.Topology.GetNodePosition(nodeName),
		Image:                c.Config.Topology.GetNodeImage(nodeName),
		User:                 c.Config.Topology.GetNodeUser(nodeName),
		Entrypoint:           c.Config.Topology.GetNodeEntrypoint(nodeName),
		Cmd:                  c.Config.Topology.GetNodeCmd(nodeName),
		Exec:                 c.Config.Topology.GetNodeExec(nodeName),
		Env:                  c.Config.Topology.GetNodeEnv(nodeName),
		NetworkMode:          strings.ToLower(c.Config.Topology.GetNodeNetworkMode(nodeName)),
		MgmtIPv4Address:      nodeDef.GetMgmtIPv4(),
		MgmtIPv4PrefixLength: nodeDef.GetMgmtIPv4Prefix(),
		MgmtIPv6Address:      nodeDef.GetMgmtIPv6(),
		Publish:              c.Config.Topology.GetNodePublish(nodeName),
		Sysctls:              c.Config.Topology.GetSysCtl(nodeName),
		Endpoints:            make([]types.Endpoint, 0),
		Sandbox:              c.Config.Topology.GetNodeSandbox(nodeName),
		Kernel:               c.Config.Topology.GetNodeKernel(nodeName),
		Runtime:              c.Config.Topology.GetNodeRuntime(nodeName),
		CPU:                  c.Config.Topology.GetNodeCPU(nodeName),
		CPUSet:               c.Config.Topology.GetNodeCPUSet(nodeName),
		Memory:               c.Config.Topology.GetNodeMemory(nodeName),
		StartupDelay:         c.Config.Topology.GetNodeStartupDelay(nodeName),
		SANs:                 c.Config.Topology.GetSANs(nodeName),

		// Extras
		Extras:  c.Config.Topology.GetNodeExtras(nodeName),
		WaitFor: c.Config.Topology.GetWaitFor(nodeName),
	}

	var err error

	// Load content of the EnvVarFiles
	envFileContent, err := utils.LoadEnvVarFiles(c.TopoFile.dir,
		c.Config.Topology.GetNodeEnvFiles(nodeName))
	if err != nil {
		return nil, err
	}
	// Merge EnvVarFiles content and the existing env variable
	nodeCfg.Env = utils.MergeStringMaps(envFileContent, nodeCfg.Env)

	log.Debugf("node config: %+v", nodeCfg)
	// initialize config
	p, err := c.Config.Topology.GetNodeStartupConfig(nodeCfg.ShortName)
	if err != nil {
		return nil, err
	}
	// resolve the startup config path to an abs path
	nodeCfg.StartupConfig = utils.ResolvePath(p, c.TopoFile.dir)

	nodeCfg.EnforceStartupConfig = c.Config.Topology.GetNodeEnforceStartupConfig(nodeCfg.ShortName)

	// initialize license field
	p, err = c.Config.Topology.GetNodeLicense(nodeCfg.ShortName)
	if err != nil {
		return nil, err
	}
	// resolve the lic path to an abs path
	nodeCfg.License = utils.ResolvePath(p, c.TopoFile.dir)

	// initialize bind mounts
	binds := c.Config.Topology.GetNodeBinds(nodeName)
	err = c.resolveBindPaths(binds, nodeCfg.LabDir)
	if err != nil {
		return nil, err
	}
	nodeCfg.Binds = binds
	nodeCfg.PortSet, nodeCfg.PortBindings, err = c.Config.Topology.GetNodePorts(nodeName)
	if err != nil {
		return nil, err
	}

	nodeCfg.Labels = c.Config.Topology.GetNodeLabels(nodeCfg.ShortName)

	nodeCfg.Config = c.Config.Topology.GetNodeConfigDispatcher(nodeCfg.ShortName)

	return nodeCfg, nil
}

// NewLink initializes a new link object.
func (c *CLab) NewLink(l *types.LinkConfig) *types.Link {
	if len(l.Endpoints) != 2 {
		log.Fatalf("endpoint %q has wrong syntax, unexpected number of items", l.Endpoints) // skipcq: RVV-A0003
	}

	return &types.Link{
		A:      c.NewEndpoint(l.Endpoints[0]),
		B:      c.NewEndpoint(l.Endpoints[1]),
		MTU:    DefaultVethLinkMTU,
		Labels: l.Labels,
		Vars:   l.Vars,
	}
}

// NewEndpoint initializes a new endpoint object.
func (c *CLab) NewEndpoint(e string) *types.Endpoint {
	// initialize a new endpoint
	endpoint := new(types.Endpoint)

	// split the string to get node name and endpoint name
	split := strings.Split(e, ":")
	if len(split) != 2 {
		log.Fatalf("endpoint %s has wrong syntax", e) // skipcq: GO-S0904, RVV-A0003
	}
	nName := split[0] // node name

	// initialize the endpoint name based on the split function
	endpoint.EndpointName = split[1] // endpoint name
	if len(endpoint.EndpointName) > 15 {
		log.Fatalf("interface '%s' name exceeds maximum length of 15 characters",
			endpoint.EndpointName) // skipcq: RVV-A0003
	}
	// generate unique MAC
	endpoint.MAC = utils.GenMac(ClabOUI)

	// search the node pointer for a node name referenced in endpoint section
	switch nName {
	// "host" is a special reference to host namespace
	// for which we create an special Node with kind "host"
	case "host":
		endpoint.Node = &types.NodeConfig{
			Kind:             "host",
			ShortName:        "host",
			NSPath:           hostNSPath,
			DeploymentStatus: "created",
		}
	// mgmt-net is a special reference to a bridge of the docker network
	// that is used as the management network
	case "mgmt-net":
		endpoint.Node = &types.NodeConfig{
			Kind:             "bridge",
			ShortName:        "mgmt-net",
			DeploymentStatus: "created",
		}
	default:
		c.m.Lock()
		if n, ok := c.Nodes[nName]; ok {
			endpoint.Node = n.Config()
			n.Config().Endpoints = append(n.Config().Endpoints, *endpoint)
		}
		c.m.Unlock()
	}

	// stop the deployment if the matching node element was not found
	// "host" node name is an exception, it may exist without a matching node
	if endpoint.Node == nil {
		log.Fatalf("not all nodes are specified in the 'topology.nodes' section or the names don't match in the 'links.endpoints' section: %s", nName) // skipcq: GO-S0904, RVV-A0003
	}

	return endpoint
}

// CheckTopologyDefinition runs topology checks and returns any errors found.
func (c *CLab) CheckTopologyDefinition(ctx context.Context) error {
	var err error
	if err = c.verifyBridgesExist(); err != nil {
		return err
	}
	if err = c.verifyLinks(); err != nil {
		return err
	}
	if err = c.verifyDuplicateAddresses(); err != nil {
		return err
	}
	if err = c.verifyRootNetnsInterfaceUniqueness(); err != nil {
		return err
	}
	if err = c.VerifyContainersUniqueness(ctx); err != nil {
		return err
	}
	if err = c.verifyVirtSupport(); err != nil {
		return err
	}
	if err = c.verifyHostIfaces(); err != nil {
		return err
	}
	if err = c.verifyLicFilesExist(); err != nil {
		return err
	}
	if err = c.verifyStartupConfigFilesExist(); err != nil {
		return err
	}
	if err = c.checkIfSignatures(); err != nil {
		return err
	}
	if err = c.checkHostRequirements(); err != nil {
		return err
	}

	return c.VerifyImages(ctx)
}

// VerifyBridgeExists verifies if every node of kind=bridge/ovs-bridge exists on the lab host.
func (c *CLab) verifyBridgesExist() error {
	for name, node := range c.Nodes {
		if node.Config().Kind == nodes.NodeKindBridge ||
			node.Config().Kind == nodes.NodeKindOVS {
			if _, err := netlink.LinkByName(name); err != nil {
				return fmt.Errorf("bridge %s is referenced in the endpoints section but was not found in the default network namespace", name)
			}
		}
	}
	return nil
}

func (c *CLab) checkHostRequirements() error {
	for _, n := range c.Nodes {
		// sse3 instruction set check
		if n.Config().HostRequirements.SSE3 && !cpuid.CPU.SSE3() {
			return fmt.Errorf("SSE3 CPU instruction set is required by kind '%s' but not available. If containerlab runs in a VM, check if that VM has been launched with full host-cpu support", n.Config().Kind)
		}
	}

	return nil
}

func (c *CLab) verifyLinks() error {
	endpoints := map[string]struct{}{}
	// dups accumulates duplicate links
	dups := []string{}
	for _, lc := range c.Config.Topology.Links {
		for _, e := range lc.Endpoints {
			if err := checkEndpoint(e); err != nil {
				return err
			}
			if _, ok := endpoints[e]; ok {
				dups = append(dups, e)
			}
			endpoints[e] = struct{}{}
		}
	}
	if len(dups) != 0 {
		return fmt.Errorf("endpoints %q appeared more than once in the links section of the topology file", dups)
	}
	return nil
}

// verifyDuplicateAddresses checks that every static IP address in the topology is unique.
func (c *CLab) verifyDuplicateAddresses() error {
	dupIps := map[string]struct{}{}
	for _, node := range c.Nodes {
		ips := []string{node.Config().MgmtIPv4Address, node.Config().MgmtIPv6Address}
		if ips[0] == "" && ips[1] == "" {
			continue
		}

		for _, ip := range ips {
			if _, exists := dupIps[ip]; !exists {
				if ip == "" {
					continue
				}
				dupIps[ip] = struct{}{}
			} else {
				return fmt.Errorf("management IP address %s appeared more than once in the topology file", ip)
			}
		}
	}

	return nil
}

// verifyLicFilesExist checks if referenced license files exist.
func (c *CLab) verifyLicFilesExist() error {
	for _, node := range c.Nodes {
		lic := node.Config().License
		if lic == "" {
			continue
		}

		rlic := utils.ResolvePath(lic, c.TopoFile.dir)
		if !utils.FileExists(rlic) {
			return fmt.Errorf("node's %q license file not found by the path %s", node.Config().ShortName, rlic)
		}
	}

	return nil
}

// verifyStartupConfigFilesExist checks if referenced startup config files exist.
func (c *CLab) verifyStartupConfigFilesExist() error {
	for _, node := range c.Nodes {
		cfg := node.Config().StartupConfig
		if cfg == "" {
			continue
		}

		rcfg := utils.ResolvePath(cfg, c.TopoFile.dir)
		if !utils.FileExists(rcfg) {
			return fmt.Errorf("node's %q startup-config file not found by the path %s", node.Config().ShortName, rcfg)
		}
	}

	return nil
}

// VerifyImages will check if image referred in the node config
// either pullable or is available in the local image store.
func (c *CLab) VerifyImages(ctx context.Context) error {
	images := make(map[string]string)

	for _, node := range c.Nodes {
		for imageKey, imageName := range node.GetImages() {
			if imageName == "" {
				return fmt.Errorf("missing required %q image for node %q", imageKey, node.Config().ShortName)
			}
			images[imageName] = node.GetRuntime().GetName()
		}
	}

	for image, runtimeName := range images {
		err := c.Runtimes[runtimeName].PullImageIfRequired(ctx, image)
		if err != nil {
			return err
		}
	}
	return nil
}

// VerifyContainersUniqueness ensures that nodes defined in the topology do not have names of the existing containers
// additionally it checks that the lab name is unique and no containers are currently running with the same lab name label.
func (c *CLab) VerifyContainersUniqueness(ctx context.Context) error {
	nctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	containers, err := c.ListContainers(nctx, nil)
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return nil
	}

	dups := []string{}
	for _, n := range c.Nodes {
		for _, cnt := range containers {
			if "/"+n.Config().LongName == cnt.Names[0] {
				dups = append(dups, n.Config().LongName)
			}
		}
	}
	if len(dups) != 0 {
		return fmt.Errorf("containers %q already exist. Add '--reconfigure' flag to the deploy command to first remove the containers and then deploy the lab", dups)
	}

	// check that none of the existing containers has a label that matches
	// the lab name of a currently deploying lab
	// this ensures lab uniqueness
	for _, cnt := range containers {
		if cnt.Labels[ContainerlabLabel] == c.Config.Name {
			return fmt.Errorf("the '%s' lab has already been deployed. Destroy the lab before deploying a lab with the same name", c.Config.Name)
		}
	}

	return nil
}

// verifyHostIfaces ensures that host interfaces referenced in the topology
// do not exist already in the root namespace
// and ensure that nodes that are configured with host networking mode do not have any interfaces defined.
func (c *CLab) verifyHostIfaces() error {
	for _, l := range c.Links {
		if l.A.Node.ShortName == "host" {
			if nl, _ := netlink.LinkByName(l.A.EndpointName); nl != nil {
				return fmt.Errorf("host interface %s referenced in topology already exists", l.A.EndpointName)
			}
		}
		if l.A.Node.NetworkMode == "host" {
			return fmt.Errorf("node '%s' is defined with host network mode, it can't have any links. Remove '%s' node links from the topology definition",
				l.A.Node.ShortName, l.A.Node.ShortName)
		}
		if l.B.Node.ShortName == "host" {
			if nl, _ := netlink.LinkByName(l.B.EndpointName); nl != nil {
				return fmt.Errorf("host interface %s referenced in topology already exists", l.B.EndpointName)
			}
		}
		if l.B.Node.NetworkMode == "host" {
			return fmt.Errorf("node '%s' is defined with host network mode, it can't have any links. Remove '%s' node links from the topology definition",
				l.B.Node.ShortName, l.B.Node.ShortName)
		}
	}
	return nil
}

// verifyRootNetnsInterfaceUniqueness ensures that interafaces that appear in the root ns (bridge, ovs-bridge and host)
// are uniquely defined in the topology file.
func (c *CLab) verifyRootNetnsInterfaceUniqueness() error {
	rootNsIfaces := map[string]struct{}{}
	for _, l := range c.Links {
		endpoints := [2]*types.Endpoint{l.A, l.B}
		for _, e := range endpoints {
			if e.Node.Kind == nodes.NodeKindBridge || e.Node.Kind == nodes.NodeKindOVS || e.Node.Kind == nodes.NodeKindHOST {
				if _, ok := rootNsIfaces[e.EndpointName]; ok {
					return fmt.Errorf(`interface %s defined for node %s has already been used in other bridges, ovs-bridges or host interfaces.
					Make sure that nodes of these kinds use unique interface names`, e.EndpointName, e.Node.ShortName)
				}
				rootNsIfaces[e.EndpointName] = struct{}{}
			}
		}
	}
	return nil
}

// verifyVirtSupport checks if virtualization is supported by a cpu in case topology contains VM-based nodes
// when clab itself is being invoked as a container, this check is bypassed.
func (c *CLab) verifyVirtSupport() error {
	// check if we are being executed in a container environment
	// in that case we skip this check as there are no convenient ways to interrogate hosts capabilities
	// check if /proc/2 exists, and if it does, check if the name of the proc is kthreadd
	// otherwise it is a container env

	f, err := os.Open("/proc/2/status")
	if err != nil {
		log.Debug("proc/2/status file was not found. This means we run in a container and no virt checks are possible")
		return nil
	}

	defer f.Close()

	// read first line of a /proc/2/status file to check if it contains kthreadd
	// if it doesn't, we are in a container
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	if !strings.Contains(scanner.Text(), "kthreadd") {
		log.Debug("proc/2/status first line doesn't contain kthreadd. This means we run in a container and no virt checks are possible")
		return nil
	}

	virtRequired := false
	for _, n := range c.Nodes {
		if n.Config().HostRequirements.VirtRequired {
			virtRequired = true

			log.Debug("detected virtualization required")

			break
		}
	}

	if !virtRequired {
		return nil
	}

	f, err = os.Open("/proc/cpuinfo")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner = bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "vmx") ||
			strings.Contains(scanner.Text(), "svm") {

			log.Debug("virtualization support found")

			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return fmt.Errorf("virtualization seems to be not supported and it is required for VM-based nodes.\nEnable virtualization or, in case you're using a VM, make sure virtualization flags are propagated to a guest")
}

// checkIfSignatures ensures that users provide valid endpoint names
// mostly this is useful for srlinux nodes which require special naming convention to be followed.
func (c *CLab) checkIfSignatures() error {
	for _, node := range c.Nodes {
		switch {
		case node.Config().Kind == nodes.NodeKindSRL:
			// regexp to match srlinux endpoint names
			srlIfRe := regexp.MustCompile(`e\d+-\d+(-\d+)?`)
			for _, e := range node.Config().Endpoints {
				if !srlIfRe.MatchString(e.EndpointName) {
					return fmt.Errorf("nokia sr linux endpoint %q doesn't match required pattern. SR Linux endpoints should be named as e1-1 or e1-1-1", e.EndpointName)
				}
			}
		}
	}

	return nil
}

// checkEndpoint runs checks on the endpoint syntax.
func checkEndpoint(e string) error {
	split := strings.Split(e, ":")
	if len(split) != 2 {
		return fmt.Errorf("malformed endpoint definition: %s", e)
	}
	if split[1] == "eth0" {
		return fmt.Errorf("eth0 interface can't be used in the endpoint definition as it is added by docker automatically: '%s'", e)
	}
	return nil
}

// resolveBindPaths resolves the host paths in a bind string, such as /hostpath:/remotepath(:options) string
// it allows host path to have `~` and relative path to an absolute path
// the list of binds will be changed in place.
// if the host path doesn't exist, the error will be returned.
func (c *CLab) resolveBindPaths(binds []string, nodedir string) error {
	for i := range binds {
		// host path is a first element in a /hostpath:/remotepath(:options) string
		elems := strings.Split(binds[i], ":")

		// replace special variable
		r := strings.NewReplacer(clabDirVar, c.Dir.Lab, nodeDirVar, nodedir)
		hp := r.Replace(elems[0])
		hp = utils.ResolvePath(hp, c.TopoFile.dir)

		_, err := os.Stat(hp)
		if err != nil {
			// check if the hostpath mount has a reference to ansible-inventory.yml or topology-data.json
			// if that is the case, we do not emit an error on missing file, since these files
			// will be created by containerlab upon lab deployment
			labdir := filepath.Base(filepath.Dir(nodedir))
			s := strings.Split(hp, string(os.PathSeparator))
			// creating a path from last two elements of a resolved host path
			h := filepath.Join(s[len(s)-2], s[len(s)-1])
			if h != filepath.Join(labdir, "ansible-inventory.yml") &&
				h != filepath.Join(labdir, "topology-data.json") {
				return fmt.Errorf("failed to verify bind path: %v", err)
			}
		}
		elems[0] = hp
		binds[i] = strings.Join(elems, ":")
	}

	return nil
}

// CheckResources runs container host resources check.
func (c *CLab) CheckResources() error {
	vcpu := runtime.NumCPU()
	log.Debugf("Number of vcpu: %d", vcpu)
	if vcpu < 2 {
		log.Warn("Only 1 vcpu detected on this container host. Most containerlab nodes require at least 2 vcpu")
		if c.HasKind(nodes.NodeKindSRL) {
			return errors.New("not enough vcpus. Nokia SR Linux nodes require at least 2 vcpus")
		}
	}

	freeMemG := sysMemory("free") / 1024 / 1024 / 1024
	if freeMemG < 1 {
		log.Warnf("it appears that container host has low memory available: ~%dGi. This might lead to runtime errors. Consider freeing up more memory.", freeMemG)
	}

	return nil
}

// sets defaults after the topology has been parsed.
func (c *CLab) setDefaults() {
	for _, n := range c.Nodes {
		// Injecting the env var with expected number of links
		numLinks := map[string]string{
			types.CLAB_ENV_INTFS: fmt.Sprintf("%d", len(n.Config().Endpoints)),
		}
		n.Config().Env = utils.MergeStringMaps(n.Config().Env, numLinks)

	}
}

// sysMemory reports on total installed or free memory (in bytes)
// used from https://github.com/pbnjay/memory
func sysMemory(v string) uint64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	var m uint64
	// If this is a 32-bit system, then these fields are
	// uint32 instead of uint64.
	// So we always convert to uint64 to match signature.
	switch v {
	case "total":
		m = uint64(in.Totalram) * uint64(in.Unit)
	case "free":
		m = uint64(in.Freeram) * uint64(in.Unit)
	}
	return m
}

// returns nodeCfg.ShortName based on the provided containerName and labName.
func getShortName(labName string, labPrefix *string, containerName string) (string, error) {
	if *labPrefix == "" {
		return containerName, nil
	}

	result := strings.Split(containerName, "-"+labName+"-")
	if len(result) != 2 {
		return "", fmt.Errorf("failed to parse container name %q", containerName)
	}
	return result[1], nil
}

// HasKind returns true if kind k is found in the list of nodes.
func (c *CLab) HasKind(k string) bool {
	for _, n := range c.Nodes {
		if n.Config().Kind == k {
			log.Warn("found")
			return true
		}
	}

	return false
}

// addDefaultLabels adds default labels to node's config struct.
func (c *CLab) addDefaultLabels(n nodes.Node) {
	cfg := n.Config()
	if cfg.Labels == nil {
		cfg.Labels = map[string]string{}
	}

	cfg.Labels[ContainerlabLabel] = c.Config.Name
	cfg.Labels[NodeNameLabel] = cfg.ShortName
	cfg.Labels[NodeKindLabel] = cfg.Kind
	cfg.Labels[NodeTypeLabel] = cfg.NodeType
	cfg.Labels[NodeGroupLabel] = cfg.Group
	cfg.Labels[NodeLabDirLabel] = cfg.LabDir
	cfg.Labels[TopoFileLabel] = c.TopoFile.path
}

// labelsToEnvVars adds labels to env vars with CLAB_LABEL_ prefix added
// and labels value sanitized.
func labelsToEnvVars(n *types.NodeConfig) {
	if n.Env == nil {
		n.Env = map[string]string{}
	}

	for k, v := range n.Labels {
		// add the value to the node env with a prefixed and special chars cleaned up key
		n.Env["CLAB_LABEL_"+utils.ToEnvKey(k)] = v
	}
}
