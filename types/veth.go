package types

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type RawVEthLink struct {
	RawLinkTypeAlias `yaml:",inline"`
	Mtu              int            `yaml:"mtu,omitempty"`
	Endpoints        []*EndpointRaw `yaml:"endpoints"`
}

func (r *RawVEthLink) Resolve(res NodeResolver) (Link, error) {
	result := &VEthLink{
		Endpoints: make([]*Endpoint, len(r.Endpoints)),
		LinkGenericAttrs: LinkGenericAttrs{
			Labels: r.Labels,
			Vars:   r.Vars,
		},
		Mtu: r.Mtu,
	}

	var err error
	for idx, e := range r.Endpoints {
		result.Endpoints[idx], err = e.Resolve(res)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func vEthFromLinkConfig(lc LinkConfig) (*RawVEthLink, error) {
	nodeA, nodeAIf, nodeB, nodeBIf := extractHostNodeInterfaceData(lc, 0)

	result := &RawVEthLink{
		RawLinkTypeAlias: RawLinkTypeAlias{
			Type:     string(LinkTypeVEth),
			Labels:   lc.Labels,
			Vars:     lc.Vars,
			Instance: nil,
		},
		Mtu: lc.MTU,
		Endpoints: []*EndpointRaw{
			{
				Node:  nodeA,
				Iface: nodeAIf,
			},
			{
				Node:  nodeB,
				Iface: nodeBIf,
			},
		},
	}
	return result, nil
}

type VEthLink struct {
	LinkGenericAttrs
	Mtu       int
	Endpoints []*Endpoint
}

func (m *VEthLink) GetType() (LinkType, error) {
	return LinkTypeVEth, nil
}

func (m *VEthLink) Deploy(ctx context.Context) error {
	linkA := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:         m.Endpoints[0].GetRandName(),
			HardwareAddr: m.Endpoints[0].MacAddress,
			Flags:        net.FlagUp,
			MTU:          m.Mtu,
		},
		PeerName:         m.Endpoints[1].GetRandName(),
		PeerHardwareAddr: m.Endpoints[1].MacAddress,
	}

	// Add the veth Pair
	if err := netlink.LinkAdd(linkA); err != nil {
		return err
	}
	// acquire netlink.Link via peer name
	linkB, err := netlink.LinkByName(m.Endpoints[1].GetRandName())
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", m.Endpoints[1].GetRandName(), err)
	}

	// diable TXOffloading for the endpoints
	for _, e := range m.Endpoints {
		err := e.DisableTxOffload(TxOffloadLinkNameRandom)
		if err != nil {
			return err
		}
	}

	// push interfaces to namespaces and rename to final interface names
	links := []netlink.Link{linkA, linkB}
	for idx, endpoint := range m.Endpoints {
		err := toNS(links[idx], endpoint.Node.GetNamespacePath(), endpoint.Iface)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *VEthLink) Remove(_ context.Context) error {
	// TODO
	log.Warn("not implemented yet")
	return nil
}

func (m *VEthLink) GetEndpoints() []*Endpoint {
	return m.Endpoints
}
