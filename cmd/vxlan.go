package cmd

import (
	"fmt"
	"net"

	"github.com/jsimonetti/rtnetlink/rtnl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srl-wim/container-lab/clab"
	"github.com/vishvananda/netlink"
)

var vxlanRemote string
var cntLink string
var parentDev string
var vxlanMTU int
var vxlanID int
var delPrefix string

func init() {
	toolsCmd.AddCommand(vxlanCmd)
	vxlanCmd.AddCommand(vxlanCreateCmd)
	vxlanCmd.AddCommand(vxlanDeleteCmd)

	vxlanCreateCmd.Flags().IntVarP(&vxlanID, "id", "i", 10, "VxLAN ID (VNI)")
	vxlanCreateCmd.Flags().StringVarP(&vxlanRemote, "remote", "r", "", "address of the remote VTEP")
	vxlanCreateCmd.Flags().StringVarP(&parentDev, "dev", "", "", "parent (source) interface name for VxLAN")
	vxlanCreateCmd.Flags().StringVarP(&cntLink, "link", "l", "", "link to which 'attach' vxlan tunnel with tc redirect")
	vxlanCreateCmd.Flags().IntVarP(&vxlanMTU, "mtu", "m", 1554, "VxLAN MTU")

	vxlanCreateCmd.MarkFlagRequired("remote")
	vxlanCreateCmd.MarkFlagRequired("id")
	vxlanCreateCmd.MarkFlagRequired("link")

	vxlanDeleteCmd.Flags().StringVarP(&delPrefix, "prefix", "p", "vx-", "delete all containerlab created VxLAN interfaces which start with this prefix")
}

// vxlanCmd represents the vxlan command container
var vxlanCmd = &cobra.Command{
	Use:   "vxlan",
	Short: "VxLAN interface commands",
}

var vxlanCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create vxlan interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := netlink.LinkByName(cntLink); err != nil {
			return fmt.Errorf("failed to lookup link %q: %v",
				cntLink, err)
		}
		// if vxlan device was not set specifically, we will use
		// the device that is reported by `ip route get $remote`
		if parentDev == "" {
			conn, err := rtnl.Dial(nil)
			if err != nil {
				log.Fatal("can't establish netlink connection: ", err)
			}
			defer conn.Close()
			r, err := conn.RouteGet(net.ParseIP(vxlanRemote))
			if err != nil {
				return fmt.Errorf("failed to find a route to VxLAN remote address %s", vxlanRemote)
			}
			parentDev = r.Interface.Name
		}

		vxlanCfg := clab.VxLAN{
			Name:     "vx-" + cntLink,
			ID:       vxlanID,
			ParentIf: parentDev,
			Remote:   net.ParseIP(vxlanRemote),
			MTU:      vxlanMTU,
		}

		if err := clab.AddVxLanInterface(vxlanCfg); err != nil {
			return err
		}

		if err := clab.BindIfacesWithTC(vxlanCfg.Name, cntLink); err != nil {
			return err
		}

		return nil
	},
}

var vxlanDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete vxlan interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		var ls []netlink.Link
		var err error

		ls, err = clab.GetLiknsByNamePrefix(delPrefix)

		if err != nil {
			return err
		}
		for _, l := range ls {
			if l.Type() != "vxlan" {
				continue
			}
			log.Infof("Deleting VxLAN link %s", l.Attrs().Name)
			err := netlink.LinkDel(l)
			if err != nil {
				log.Warnf("error when deleting link %s: %v", l.Attrs().Name, err)
			}
		}
		return nil
	},
}
