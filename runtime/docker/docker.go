// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	dockerC "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/dustin/go-humanize"
	"github.com/google/shlex"
	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/runtime"
	"github.com/srl-labs/containerlab/types"
	"github.com/srl-labs/containerlab/utils"
	"github.com/vishvananda/netlink"
)

const (
	runtimeName    = "docker"
	sysctlBase     = "/proc/sys"
	defaultTimeout = 30 * time.Second
)

func init() {
	runtime.Register(runtimeName, func() runtime.ContainerRuntime {
		return &DockerRuntime{
			mgmt: new(types.MgmtNet),
		}
	})
}

type DockerRuntime struct {
	config runtime.RuntimeConfig
	Client *dockerC.Client
	mgmt   *types.MgmtNet
}

func (d *DockerRuntime) Init(opts ...runtime.RuntimeOption) error {
	var err error
	log.Debug("Runtime: Docker")
	d.Client, err = dockerC.NewClientWithOpts(dockerC.FromEnv, dockerC.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	for _, o := range opts {
		o(d)
	}
	return nil
}

func (d *DockerRuntime) WithKeepMgmtNet() {
	d.config.KeepMgmtNet = true
}
func (*DockerRuntime) GetName() string                 { return runtimeName }
func (d *DockerRuntime) Config() runtime.RuntimeConfig { return d.config }

// Mgmt return management network struct of a runtime
func (d *DockerRuntime) Mgmt() *types.MgmtNet { return d.mgmt }

func (d *DockerRuntime) WithConfig(cfg *runtime.RuntimeConfig) {
	d.config.Timeout = cfg.Timeout
	d.config.Debug = cfg.Debug
	d.config.GracefulShutdown = cfg.GracefulShutdown
	if d.config.Timeout <= 0 {
		d.config.Timeout = defaultTimeout
	}
}

func (d *DockerRuntime) WithMgmtNet(n *types.MgmtNet) {
	d.mgmt = n
}

// CreateDockerNet creates a docker network or reusing if it exists
func (d *DockerRuntime) CreateNet(ctx context.Context) (err error) {
	nctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	// linux bridge name that is used by docker network
	bridgeName := d.mgmt.Bridge

	log.Debugf("Checking if docker network %q exists", d.mgmt.Network)
	netResource, err := d.Client.NetworkInspect(nctx, d.mgmt.Network, dockerTypes.NetworkInspectOptions{})
	switch {
	case dockerC.IsErrNotFound(err):
		log.Debugf("Network %q does not exist", d.mgmt.Network)
		log.Infof("Creating docker network: Name=%q, IPv4Subnet=%q, IPv6Subnet=%q, MTU=%q",
			d.mgmt.Network, d.mgmt.IPv4Subnet, d.mgmt.IPv6Subnet, d.mgmt.MTU)

		enableIPv6 := false
		var ipamConfig []network.IPAMConfig

		// check if IPv4/6 addr are assigned to a mgmt bridge
		var v4gw, v6gw string
		if d.mgmt.Bridge != "" {
			v4gw, v6gw, err = utils.FirstLinkIPs(d.mgmt.Bridge)
			if err != nil {
				// only return error if the error is not about link not found
				// we will create the bridge if it doesn't exist
				if !errors.As(err, &netlink.LinkNotFoundError{}) {
					return err
				}
			}
			log.Debugf("bridge %q has ipv4 adrr of %q and ipv6 addr of %q", d.mgmt.Bridge, v4gw, v6gw)
		}

		if d.mgmt.IPv4Subnet != "" {
			if d.mgmt.IPv4Gw != "" {
				v4gw = d.mgmt.IPv4Gw
			}
			ipamConfig = append(ipamConfig, network.IPAMConfig{
				Subnet:  d.mgmt.IPv4Subnet,
				Gateway: v4gw,
			})
		}

		if d.mgmt.IPv6Subnet != "" {
			if d.mgmt.IPv6Gw != "" {
				v6gw = d.mgmt.IPv6Gw
			}
			ipamConfig = append(ipamConfig, network.IPAMConfig{
				Subnet:  d.mgmt.IPv6Subnet,
				Gateway: v6gw,
			})
			enableIPv6 = true
		}

		ipam := &network.IPAM{
			Driver: "default",
			Config: ipamConfig,
		}

		netwOpts := map[string]string{
			"com.docker.network.driver.mtu": d.mgmt.MTU,
		}

		if bridgeName != "" {
			netwOpts["com.docker.network.bridge.name"] = bridgeName
		}

		opts := dockerTypes.NetworkCreate{
			CheckDuplicate: true,
			Driver:         "bridge",
			EnableIPv6:     enableIPv6,
			IPAM:           ipam,
			Internal:       false,
			Attachable:     false,
			Labels: map[string]string{
				"containerlab": "",
			},
			Options: netwOpts,
		}

		netCreateResponse, err := d.Client.NetworkCreate(nctx, d.mgmt.Network, opts)
		if err != nil {
			return err
		}

		if len(netCreateResponse.ID) < 12 {
			return fmt.Errorf("could not get bridge ID")
		}
		// when bridge is not set by a user explicitly
		// we use the 12 chars of docker net as its name
		if bridgeName == "" {
			bridgeName = "br-" + netCreateResponse.ID[:12]
		}

	case err == nil:
		log.Debugf("network %q was found. Reusing it...", d.mgmt.Network)
		if len(netResource.ID) < 12 {
			return fmt.Errorf("could not get bridge ID")
		}
		switch d.mgmt.Network {
		case "bridge":
			bridgeName = "docker0"
		default:
			if netResource.Options["com.docker.network.bridge.name"] != "" {
				bridgeName = netResource.Options["com.docker.network.bridge.name"]
			} else {
				bridgeName = "br-" + netResource.ID[:12]
			}
		}

	default:
		return err
	}

	if d.mgmt.Bridge == "" {
		d.mgmt.Bridge = bridgeName
	}

	log.Debugf("Docker network %q, bridge name %q", d.mgmt.Network, bridgeName)

	return d.postCreateNetActions()
}

// postCreateNetActions performs additional actions after the network has been created
func (d *DockerRuntime) postCreateNetActions() (err error) {
	log.Debug("Disable RPF check on the docker host")
	err = setSysctl("net/ipv4/conf/all/rp_filter", 0)
	if err != nil {
		return fmt.Errorf("failed to disable RP filter on docker host for the 'all' scope: %v", err)
	}
	err = setSysctl("net/ipv4/conf/default/rp_filter", 0)
	if err != nil {
		return fmt.Errorf("failed to disable RP filter on docker host for the 'default' scope: %v", err)
	}

	log.Debugf("Enable LLDP on the linux bridge %s", d.mgmt.Bridge)
	file := "/sys/class/net/" + d.mgmt.Bridge + "/bridge/group_fwd_mask"

	err = ioutil.WriteFile(file, []byte(strconv.Itoa(16384)), 0640)
	if err != nil {
		log.Warnf("failed to enable LLDP on docker bridge: %v", err)
	}

	log.Debugf("Disabling TX checksum offloading for the %s bridge interface...", d.mgmt.Bridge)
	err = utils.EthtoolTXOff(d.mgmt.Bridge)
	if err != nil {
		log.Warnf("failed to disable TX checksum offloading for the %s bridge interface: %v", d.mgmt.Bridge, err)
	}
	err = d.installIPTablesFwdRule()
	if err != nil {
		log.Warnf("%v", err)
	}

	return nil
}

// DeleteNet deletes a docker bridge
func (d *DockerRuntime) DeleteNet(ctx context.Context) (err error) {
	network := d.mgmt.Network
	if network == "bridge" || d.config.KeepMgmtNet {
		log.Debugf("Skipping deletion of %q network", network)
		return nil
	}
	nctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	nres, err := d.Client.NetworkInspect(ctx, network, dockerTypes.NetworkInspectOptions{})
	if err != nil {
		return err
	}
	numEndpoints := len(nres.Containers)
	if numEndpoints > 0 {
		if d.config.Debug {
			log.Debugf("network %q has %d active endpoints, deletion skipped", d.mgmt.Network, numEndpoints)
			for _, endp := range nres.Containers {
				log.Debugf("%q is connected to %s", endp.Name, network)
			}
		}
		return nil
	}
	err = d.Client.NetworkRemove(nctx, network)
	if err != nil {
		return err
	}

	// bridge name associated with the network
	br := "br-" + nres.ID[:12]
	err = d.deleteIPTablesFwdRule(br)
	if err != nil {
		log.Warnf("%v", err)
	}

	return nil
}

// CreateContainer creates a docker container (but does not start it)
func (d *DockerRuntime) CreateContainer(ctx context.Context, node *types.NodeConfig) (string, error) {
	log.Infof("Creating container: %q", node.ShortName)
	nctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	cmd, err := shlex.Split(node.Cmd)
	if err != nil {
		return "", err
	}

	var entrypoint []string
	if node.Entrypoint != "" {
		entrypoint, err = shlex.Split(node.Entrypoint)
		if err != nil {
			return "", err
		}
	}

	containerConfig := &container.Config{
		Image:        node.Image,
		Entrypoint:   entrypoint,
		Cmd:          cmd,
		Env:          utils.ConvertEnvs(node.Env),
		AttachStdout: true,
		AttachStderr: true,
		Hostname:     node.ShortName,
		Tty:          true,
		User:         node.User,
		Labels:       node.Labels,
		ExposedPorts: node.PortSet,
		MacAddress:   node.MacAddress,
	}
	containerHostConfig := &container.HostConfig{
		Binds:        node.Binds,
		PortBindings: node.PortBindings,
		Sysctls:      node.Sysctls,
		Privileged:   true,
		// Network mode will be defined below via switch
		NetworkMode: "",
		ExtraHosts:  node.ExtraHosts, // add static /etc/hosts entries
	}
	var resources container.Resources
	if node.Memory != "" {
		mem, err := humanize.ParseBytes(node.Memory)
		if err != nil {
			return "", err
		}
		resources.Memory = int64(mem)
	}
	if node.CPU != 0 {
		resources.CPUQuota = int64(node.CPU * 100000)
		resources.CPUPeriod = 100000
	}
	if node.CPUSet != "" {
		resources.CpusetCpus = node.CPUSet
	}
	containerHostConfig.Resources = resources
	containerNetworkingConfig := &network.NetworkingConfig{}

	netMode := strings.SplitN(node.NetworkMode, ":", 2)
	switch netMode[0] {
	case "container":
		// We expect exactly two arguments in this case ("container" keyword & cont. name/ID)
		if len(netMode) != 2 {
			return "", fmt.Errorf("container network mode was specified for container %q, but no container name was found: %q", node.ShortName, netMode)
		}
		// also cont. ID shouldn't be empty
		if netMode[1] == "" {
			return "", fmt.Errorf("container network mode was specified for container %q, but no container name was found: %q", node.ShortName, netMode)
		}
		// Extract lab/topo prefix to provide a full (long) container name. Hackish way.
		prefix := strings.SplitN(node.LongName, node.ShortName, 2)[0]
		// Compile the net spec
		containerHostConfig.NetworkMode = container.NetworkMode("container:" + prefix + netMode[1])
		// unset the hostname as it is not supported in this case
		containerConfig.Hostname = ""
	case "host":
		containerHostConfig.NetworkMode = "host"
	default:
		containerHostConfig.NetworkMode = container.NetworkMode(d.mgmt.Network)

		containerNetworkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			d.mgmt.Network: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: node.MgmtIPv4Address,
					IPv6Address: node.MgmtIPv6Address,
				},
			},
		}
	}

	// regular linux containers may benefit from automatic restart on failure
	// note, that veth pairs added to this container (outside of eth0) will be lost on restart
	if node.Kind == "linux" {
		containerHostConfig.RestartPolicy.Name = "on-failure"
	}

	cont, err := d.Client.ContainerCreate(
		nctx,
		containerConfig,
		containerHostConfig,
		containerNetworkingConfig,
		nil,
		node.LongName,
	)
	log.Debugf("Container %q create response: %+v", node.ShortName, cont)
	if err != nil {
		return "", err
	}
	return cont.ID, nil
}

// GetNSPath inspects a container by its name/id and returns an netns path using the pid of a container
func (d *DockerRuntime) GetNSPath(ctx context.Context, cID string) (string, error) {
	nctx, cancelFn := context.WithTimeout(ctx, d.config.Timeout)
	defer cancelFn()
	cJSON, err := d.Client.ContainerInspect(nctx, cID)
	if err != nil {
		return "", err
	}
	return "/proc/" + strconv.Itoa(cJSON.State.Pid) + "/ns/net", nil
}

func (d *DockerRuntime) PullImageIfRequired(ctx context.Context, imageName string) error {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)

	ilo := dockerTypes.ImageListOptions{
		All:     false,
		Filters: filter,
	}

	log.Debugf("Looking up %s Docker image", imageName)

	images, err := d.Client.ImageList(ctx, ilo)
	if err != nil {
		return err
	}

	// If Image doesn't exist, we need to pull it
	if len(images) > 0 {
		log.Debugf("Image %s present, skip pulling", imageName)
		return nil
	}

	canonicalImageName := utils.GetCanonicalImageName(imageName)
	authString := ""

	// get docker config based on an empty path (default docker config path will be assumed)
	dockerConfig, err := GetDockerConfig("")
	if err != nil {
		log.Debug("docker config file not found")
	} else {
		authString, err = GetDockerAuth(dockerConfig, canonicalImageName)
		if err != nil {
			return err
		}
	}

	log.Infof("Pulling %s Docker image", canonicalImageName)
	reader, err := d.Client.ImagePull(ctx, canonicalImageName, dockerTypes.ImagePullOptions{
		RegistryAuth: authString,
	})
	if err != nil {
		return err
	}
	defer reader.Close()
	// must read from reader, otherwise image is not properly pulled
	_, _ = io.Copy(ioutil.Discard, reader)
	log.Infof("Done pulling %s", canonicalImageName)

	return nil
}

// StartContainer starts a docker container
func (d *DockerRuntime) StartContainer(ctx context.Context, cID string, node *types.NodeConfig) (interface{}, error) {
	nctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()
	log.Debugf("Start container: %q", node.LongName)
	err := d.Client.ContainerStart(nctx,
		cID,
		dockerTypes.ContainerStartOptions{
			CheckpointID:  "",
			CheckpointDir: "",
		},
	)
	if err != nil {
		return nil, err
	}
	log.Debugf("Container started: %q", node.LongName)
	err = d.postStartActions(ctx, cID, node)
	return nil, err
}

// postStartActions performs misc. tasks that are needed after the container starts
func (d *DockerRuntime) postStartActions(ctx context.Context, cID string, node *types.NodeConfig) error {
	var err error
	node.NSPath, err = d.GetNSPath(ctx, cID)
	if err != nil {
		return err
	}
	err = utils.LinkContainerNS(node.NSPath, node.LongName)
	return err
}

// ListContainers lists all containers with labels []string
func (d *DockerRuntime) ListContainers(ctx context.Context, gfilters []*types.GenericFilter) ([]types.GenericContainer, error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
	defer cancel()

	filter := d.buildFilterString(gfilters)
	ctrs, err := d.Client.ContainerList(ctx, dockerTypes.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return nil, err
	}
	var nr []dockerTypes.NetworkResource
	if d.mgmt.Network == "" {
		nctx, cancel := context.WithTimeout(ctx, d.config.Timeout)
		defer cancel()
		// fetch containerlab created networks
		f := filters.NewArgs()
		f.Add("label", "containerlab")
		nr, err = d.Client.NetworkList(nctx, dockerTypes.NetworkListOptions{
			Filters: f,
		})

		if err != nil {
			return nil, err
		}

		// fetch default bridge network
		f = filters.NewArgs()
		f.Add("name", "bridge")
		bridgenet, err := d.Client.NetworkList(nctx, dockerTypes.NetworkListOptions{
			Filters: f,
		})

		if err != nil {
			return nil, err
		}

		nr = append(nr, bridgenet...)
	}
	return d.produceGenericContainerList(ctrs, nr)
}

func (d *DockerRuntime) GetContainer(ctx context.Context, cID string) (*types.GenericContainer, error) {
	var ctr *types.GenericContainer
	gFilter := types.GenericFilter{
		FilterType: "name",
		Field:      "",
		Operator:   "",
		Match:      cID,
	}
	ctrs, err := d.ListContainers(ctx, []*types.GenericFilter{&gFilter})
	if err != nil {
		return ctr, err
	}
	if len(ctrs) != 1 {
		return ctr, fmt.Errorf("found unexpected number of containers: %d", len(ctrs))
	}
	return &ctrs[0], nil
}

func (*DockerRuntime) buildFilterString(gfilters []*types.GenericFilter) filters.Args {
	filter := filters.NewArgs()
	for _, filterentry := range gfilters {
		filterstring := filterentry.Field
		if filterentry.Operator != "exists" {
			filterstring = filterstring + filterentry.Operator + filterentry.Match
		}
		log.Debug("Filterstring: " + filterstring)
		filter.Add(filterentry.FilterType, filterstring)
	}
	return filter
}

// Transform docker-specific to generic container format
func (d *DockerRuntime) produceGenericContainerList(inputContainers []dockerTypes.Container, inputNetworkRessources []dockerTypes.NetworkResource) ([]types.GenericContainer, error) {
	var result []types.GenericContainer

	for _, i := range inputContainers {
		ctr := types.GenericContainer{
			Names:           i.Names,
			ID:              i.ID,
			ShortID:         i.ID[:12],
			Image:           i.Image,
			State:           i.State,
			Status:          i.Status,
			Labels:          i.Labels,
			NetworkSettings: types.GenericMgmtIPs{},
		}
		bridgeName := d.mgmt.Network
		// if bridgeName is "", try to find a network created by clab that the container is connected to
		if bridgeName == "" && inputNetworkRessources != nil {
			for _, nr := range inputNetworkRessources {
				if _, ok := i.NetworkSettings.Networks[nr.Name]; ok {
					bridgeName = nr.Name
					break
				}
			}
		}
		if ifcfg, ok := i.NetworkSettings.Networks[bridgeName]; ok {
			ctr.NetworkSettings.IPv4addr = ifcfg.IPAddress
			ctr.NetworkSettings.IPv4pLen = ifcfg.IPPrefixLen
			ctr.NetworkSettings.IPv6addr = ifcfg.GlobalIPv6Address
			ctr.NetworkSettings.IPv6pLen = ifcfg.GlobalIPv6PrefixLen
		}
		result = append(result, ctr)
	}

	return result, nil
}

// Exec executes cmd on container identified with id and returns stdout, stderr bytes and an error
func (d *DockerRuntime) Exec(ctx context.Context, cID string, cmd []string) ([]byte, []byte, error) {
	cont, err := d.Client.ContainerInspect(ctx, cID)
	if err != nil {
		return nil, nil, err
	}
	execID, err := d.Client.ContainerExecCreate(ctx, cID, dockerTypes.ExecConfig{
		User:         "root",
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	if err != nil {
		log.Errorf("failed to create exec in container %s: %v", cont.Name, err)
		return nil, nil, err
	}
	log.Debugf("%s exec created %v", cont.Name, cID)

	rsp, err := d.Client.ContainerExecAttach(ctx, execID.ID, dockerTypes.ExecStartCheck{})
	if err != nil {
		log.Errorf("failed exec in container %s: %v", cont.Name, err)
		return nil, nil, err
	}
	defer rsp.Close()
	log.Debugf("%s exec attached %v", cont.Name, cID)

	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, rsp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return outBuf.Bytes(), errBuf.Bytes(), err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	return outBuf.Bytes(), errBuf.Bytes(), nil
}

// ExecNotWait executes cmd on container identified with id but doesn't wait for output nor attaches stdout/err
func (d *DockerRuntime) ExecNotWait(_ context.Context, cID string, cmd []string) error {
	execConfig := dockerTypes.ExecConfig{Tty: false, AttachStdout: false, AttachStderr: false, Cmd: cmd}
	respID, err := d.Client.ContainerExecCreate(context.Background(), cID, execConfig)
	if err != nil {
		return err
	}

	execStartCheck := dockerTypes.ExecStartCheck{}
	_, err = d.Client.ContainerExecAttach(context.Background(), respID.ID, execStartCheck)
	if err != nil {
		return err
	}
	return nil
}

// DeleteContainer tries to stop a container then remove it
func (d *DockerRuntime) DeleteContainer(ctx context.Context, cID string) error {
	var err error
	force := !d.config.GracefulShutdown
	if d.config.GracefulShutdown {
		log.Infof("Stopping container: %s", cID)
		err = d.Client.ContainerStop(ctx, cID, &d.config.Timeout)
		if err != nil {
			log.Errorf("could not stop container %q: %v", cID, err)
			force = true
		}
	}
	log.Debugf("Removing container: %s", strings.TrimLeft(cID, "/"))
	err = d.Client.ContainerRemove(ctx, cID, dockerTypes.ContainerRemoveOptions{Force: force})
	if err != nil {
		return err
	}
	log.Infof("Removed container: %s", strings.TrimLeft(cID, "/"))
	return nil
}

// setSysctl writes sysctl data by writing to a specific file
func setSysctl(sysctl string, newVal int) error {
	return ioutil.WriteFile(path.Join(sysctlBase, sysctl), []byte(strconv.Itoa(newVal)), 0640)
}

func (d *DockerRuntime) StopContainer(ctx context.Context, name string) error {
	return d.Client.ContainerKill(ctx, name, "kill")
}
