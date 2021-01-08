package clab

import (
	"fmt"
	"net"
	"os"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type vEthEndpoint struct {
	Link     netlink.Link
	LinkName string
	NSName   string // netns name
	NSPath   string // netns path
	Bridge   string // bridge name a veth is destined to be connected to
}

// CreateVirtualWiring provides the virtual topology between the containers
func (c *CLab) CreateVirtualWiring(l *Link) (err error) {
	log.Infof("Creating virtual wire: %s:%s <--> %s:%s", l.A.Node.ShortName, l.A.EndpointName, l.B.Node.ShortName, l.B.EndpointName)

	// connect containers (or container and a bridge) using veth pair
	// based on the link configuration contained within *Link struct
	// veth side A
	vA := vEthEndpoint{
		LinkName: l.A.EndpointName,
		NSName:   l.A.Node.LongName,
		NSPath:   l.A.Node.NSPath,
	}
	// veth side B
	vB := vEthEndpoint{
		LinkName: l.B.EndpointName,
		NSName:   l.B.Node.LongName,
		NSPath:   l.B.Node.NSPath,
	}

	// get random names for veth sides as they will be created in root netns first
	ARndmName := fmt.Sprintf("clab-%s", genIfName())
	BRndmName := fmt.Sprintf("clab-%s", genIfName())

	// set bridge name for endpoint that should be connect to linux bridge
	switch {
	case l.A.Node.Kind == "bridge":
		vA.Bridge = l.A.Node.ShortName
		// veth endpoint destined to connect to the bridge in the host netns
		// will not have a random name
		ARndmName = l.A.EndpointName
	case l.B.Node.Kind == "bridge":
		vB.Bridge = l.B.Node.ShortName
		BRndmName = l.B.EndpointName
	}

	// create veth pair in the root netns
	vA.Link, vB.Link, err = createVethIface(ARndmName, BRndmName, l.MTU)
	if err != nil {
		return err
	}

	// once veth pair is created, disable tx offload for veth pair
	if err := EthtoolTXOff(ARndmName); err != nil {
		return err
	}
	if err := EthtoolTXOff(BRndmName); err != nil {
		return err
	}

	if err = vA.setVethLink(); err != nil {
		netlink.LinkDel(vA.Link)
		return err
	}
	if err = vB.setVethLink(); err != nil {
		netlink.LinkDel(vB.Link)
	}
	return err

}

// createVethIface takes two veth endpoint structs and create a veth pair and return
// veth interface links.
func createVethIface(ifName, peerName string, mtu int) (linkA netlink.Link, linkB netlink.Link, err error) {
	linkA = &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:  ifName,
			Flags: net.FlagUp,
			MTU:   mtu,
		},
		PeerName: peerName,
	}

	if err := netlink.LinkAdd(linkA); err != nil {
		return nil, nil, err
	}

	if linkB, err = netlink.LinkByName(peerName); err != nil {
		err = fmt.Errorf("failed to lookup %q: %v", peerName, err)
	}

	return
}

// setVethLink sets the veth link endpoints to the relevant namespaces and/or connects one end to the bridge
func (veth *vEthEndpoint) setVethLink() error {
	// if veth is destined to connect to a linux bridge in the host netns
	if veth.Bridge != "" {
		return veth.toBridge()
	}
	// otherwise it needs to be put into a netns
	return veth.toNS()
}

// vethToNS puts a veth endpoint to a given netns and renames its random name to a desired name
func (veth *vEthEndpoint) toNS() error {
	var vethNS ns.NetNS
	var err error
	if vethNS, err = ns.GetNS(veth.NSPath); err != nil {
		return err
	}
	// move veth endpoint to namespace
	if err = netlink.LinkSetNsFd(veth.Link, int(vethNS.Fd())); err != nil {
		return err
	}
	err = vethNS.Do(func(_ ns.NetNS) error {
		if err = netlink.LinkSetName(veth.Link, veth.LinkName); err != nil {
			return fmt.Errorf(
				"failed to rename link: %v", err)
		}

		if err = netlink.LinkSetUp(veth.Link); err != nil {
			return fmt.Errorf("failed to set %q up: %v",
				veth.LinkName, err)
		}
		return nil
	})
	return err
}

func (veth *vEthEndpoint) toBridge() error {
	var vethNS ns.NetNS
	var err error
	// bride is in the host netns, thus we need to get current netns
	if vethNS, err = ns.GetCurrentNS(); err != nil {
		return err
	}
	err = vethNS.Do(func(_ ns.NetNS) error {
		br, err := bridgeByName(veth.Bridge)
		if err != nil {
			return err
		}

		// connect host veth end to the bridge
		if err := netlink.LinkSetMaster(veth.Link, br); err != nil {
			return fmt.Errorf("failed to connect %q to bridge %v: %v", veth.LinkName, veth.Bridge, err)
		}

		if err = netlink.LinkSetUp(veth.Link); err != nil {
			return fmt.Errorf("failed to set %q up: %v", veth.LinkName, err)
		}
		return nil
	})
	return err
}

// DeleteNetnsSymlinks deletes the symlink file created for each container netns
func (c *CLab) DeleteNetnsSymlinks() (err error) {
	for _, node := range c.Nodes {
		if node.Kind != "bridge" {
			log.Debugf("Deleting %s network namespace", node.LongName)
			if err := deleteNetnsSymlink(node.LongName); err != nil {
				return err
			}
		}

	}

	return nil
}

func genIfName() string {
	s, _ := uuid.New().MarshalText() // .MarshalText() always return a nil error
	return string(s[:8])
}

// deleteNetnsSymlink deletes a network namespace and removes the symlink created by linkContainerNS func
func deleteNetnsSymlink(n string) error {
	log.Debug("Deleting netns symlink: ", n)
	sl := fmt.Sprintf("/run/netns/%s", n)
	err := os.Remove(sl)
	if err != nil {
		log.Debug("Failed to delete netns symlink by path:", sl)
	}
	return nil
}

func bridgeByName(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not lookup %q: %v", name, err)
	}
	br, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}
	return br, nil
}
