package srl

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/srl-labs/containerlab/clab/exec"
)

const banner = `................................................................
:                  Welcome to Nokia SR Linux!                  :
:              Open Network OS for the NetOps era.             :
:                                                              :
:    This is a freely distributed official container image.    :
:                      Use it - Share it                       :
:                                                              :
: Get started: https://learn.srlinux.dev                       :
: Container:   https://go.srlinux.dev/container-image          :
: Docs:        https://doc.srlinux.dev/%s-%-2s                   :
: Rel. notes:  https://doc.srlinux.dev/rn%s-%s-%s               :
: YANG:        https://yang.srlinux.dev/v%s.%s.%s               :
: Discord:     https://go.srlinux.dev/discord                  :
: Contact:     https://go.srlinux.dev/contact-sales            :
................................................................
`

// banner returns a banner string with a docs version filled in based on the version information queried from the node.
func (n *srl) banner(ctx context.Context) (string, error) {
	cmd, _ := exec.NewExecCmdFromString(`sr_cli -d "info from state /system information version | grep version"`)

	execResult, err := n.RunExec(ctx, cmd)
	if err != nil {
		return "", err
	}

	log.Debugf("node %s. stdout: %s, stderr: %s", n.Cfg.ShortName, execResult.GetStdOutString(), execResult.GetStdErrString())

	n.swVersion = n.parseVersionString(execResult.GetStdOutString())

	// if minor is a single digit value, we need to add extra space to patch version
	// to have banner table aligned nicely
	if len(n.swVersion.minor) == 1 {
		n.swVersion.patch = n.swVersion.patch + " "
	}

	b := fmt.Sprintf(banner,
		n.swVersion.major, n.swVersion.minor,
		n.swVersion.major, n.swVersion.minor, n.swVersion.patch,
		n.swVersion.major, n.swVersion.minor, n.swVersion.patch)

	return b, nil
}
