package clab

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// var debug bool

type CLab struct {
	Config       *Config
	TopoFile     *TopoFile
	m            *sync.RWMutex
	Nodes        map[string]*Node
	Links        map[int]*Link
	DockerClient *docker.Client
	Dir          *cLabDirectory

	debug            bool
	timeout          time.Duration
	gracefulShutdown bool
}

type cLabDirectory struct {
	Lab       string
	LabCA     string
	LabCARoot string
	LabGraph  string
}

type ClabOption func(c *CLab)

func WithDebug(d bool) ClabOption {
	return func(c *CLab) {
		c.debug = d
	}
}

func WithTimeout(dur time.Duration) ClabOption {
	return func(c *CLab) {
		c.timeout = dur
	}
}

func WithEnvDockerClient() ClabOption {
	return func(c *CLab) {
		var err error
		c.DockerClient, err = docker.NewEnvClient()
		if err != nil {
			log.Fatalf("failed to create docker client: %v", err)
		}
	}
}

func WithTopoFile(file string) ClabOption {
	return func(c *CLab) {
		if file == "" {
			return
		}
		if err := c.GetTopology(file); err != nil {
			log.Fatalf("failed to read topology file: %v", err)
		}
	}
}

func WithGracefulShutdown(gracefulShutdown bool) ClabOption {
	return func(c *CLab) {
		c.gracefulShutdown = gracefulShutdown
	}
}

// NewContainerLab function defines a new container lab
func NewContainerLab(opts ...ClabOption) *CLab {
	c := &CLab{
		Config:   new(Config),
		TopoFile: new(TopoFile),
		m:        new(sync.RWMutex),
		Nodes:    make(map[string]*Node),
		Links:    make(map[int]*Link),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *CLab) CreateNode(ctx context.Context, node *Node, certs *Certificates) error {
	if certs != nil {
		c.m.Lock()
		node.TLSCert = string(certs.Cert)
		node.TLSKey = string(certs.Key)
		c.m.Unlock()
	}
	err := c.CreateNodeDirStructure(node)
	if err != nil {
		return err
	}
	return c.CreateContainer(ctx, node)
}

// ExecPostDeployTasks executes tasks that some nodes might require to boot properly after start
func (c *CLab) ExecPostDeployTasks(ctx context.Context, node *Node, lworkers uint) error {
	switch node.Kind {
	case "ceos":
		log.Debugf("Running postdeploy actions for Arista cEOS '%s' node", node.ShortName)
		return ceosPostDeploy(ctx, c, node, lworkers)
	case "crpd":
		// exec `service ssh restart` to start ssh service and take into account mounted sshd_config
		execConfig := types.ExecConfig{Tty: false, AttachStdout: false, AttachStderr: false, Cmd: strings.Fields("service ssh restart")}
		respID, err := c.DockerClient.ContainerExecCreate(context.Background(), node.ContainerID, execConfig)
		if err != nil {
			return err
		}
		_, err = c.DockerClient.ContainerExecAttach(context.Background(), respID.ID, execConfig)
		if err != nil {
			return err
		}

	case "linux":
		log.Debugf("Running postdeploy actions for Linux '%s' node", node.ShortName)
		return disableTxOffload(node)

	case "sonic-vs":
		log.Debugf("Running postdeploy actions for sonic-vs '%s' node", node.ShortName)
		// TODO: change this calls to c.ExecNotWait
		// exec `supervisord` to start sonic services
		execConfig := types.ExecConfig{Tty: false, AttachStdout: false, AttachStderr: false, Cmd: strings.Fields("supervisord")}
		respID, err := c.DockerClient.ContainerExecCreate(context.Background(), node.ContainerID, execConfig)
		if err != nil {
			return err
		}
		_, err = c.DockerClient.ContainerExecAttach(context.Background(), respID.ID, execConfig)
		if err != nil {
			return err
		}
		// exec `/usr/lib/frr/bgpd` to start BGP service
		execConfig = types.ExecConfig{Tty: false, AttachStdout: false, AttachStderr: false, Cmd: strings.Fields("/usr/lib/frr/bgpd")}
		respID, err = c.DockerClient.ContainerExecCreate(context.Background(), node.ContainerID, execConfig)
		if err != nil {
			return err
		}
		_, err = c.DockerClient.ContainerExecAttach(context.Background(), respID.ID, execConfig)
		if err != nil {
			return err
		}
	case "mysocketio":
		log.Debugf("Running postdeploy actions for mysocketio '%s' node", node.ShortName)
		err := disableTxOffload(node)
		if err != nil {
			return fmt.Errorf("failed to disable tx checksum offload for mysocketio kind: %v", err)
		}

		log.Infof("Creating mysocketio tunnels...")
		err = createMysocketTunnels(ctx, c, node)
		return err
	}
	return nil
}

// CreateLinks creates links using the specified number of workers
// `postdeploy` indicates the stage of links creation.
// `postdeploy=true` means the links routine is called after nodes postdeploy tasks
func (c *CLab) CreateLinks(ctx context.Context, workers uint, postdeploy bool) {
	wg := new(sync.WaitGroup)
	wg.Add(int(workers))
	linksChan := make(chan *Link)

	log.Debug("creating links...")
	// wire the links between the nodes based on cabling plan
	for i := uint(0); i < workers; i++ {
		go func(i uint) {
			defer wg.Done()
			for {
				select {
				case link := <-linksChan:
					if link == nil {
						log.Debugf("Link worker %d terminating...", i)
						return
					}
					log.Debugf("Link worker %d received link: %+v", i, link)
					if err := c.CreateVirtualWiring(link); err != nil {
						log.Error(err)
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	for _, link := range c.Links {
		// skip the links of ceos kind
		// ceos containers need to be restarted in the postdeploy stage, thus their data links
		// will get recreated after post-deploy stage
		if !postdeploy {
			if link.A.Node.Kind == "ceos" || link.B.Node.Kind == "ceos" {
				continue
			}
			linksChan <- link
		} else {
			// postdeploy stage
			// create ceos links that were skipped during original links creation
			if link.A.Node.Kind == "ceos" || link.B.Node.Kind == "ceos" {
				linksChan <- link
			}
		}
	}
	// close channel to terminate the workers
	close(linksChan)
	// wait for all workers to finish
	wg.Wait()
}

func disableTxOffload(n *Node) error {
	// disable tx checksum offload for linux containers on eth0 interfaces
	nodeNS, err := ns.GetNS(n.NSPath)
	if err != nil {
		return err
	}
	err = nodeNS.Do(func(_ ns.NetNS) error {
		// disabling offload on lo0 interface
		err := EthtoolTXOff("eth0")
		if err != nil {
			log.Infof("Failed to disable TX checksum offload for 'eth0' interface for Linux '%s' node: %v", n.ShortName, err)
		}
		return err
	})
	return err
}

func StringInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
