// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/srl-labs/containerlab/clab"
	"github.com/srl-labs/containerlab/runtime"
	"github.com/srl-labs/containerlab/types"
)

var format string
var details bool
var all bool

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:     "inspect",
	Short:   "inspect lab details",
	Long:    "show details about a particular lab or all running labs\nreference: https://containerlab.dev/cmd/inspect/",
	Aliases: []string{"ins", "i"},
	PreRunE: sudoCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		if name == "" && topo == "" && !all {
			fmt.Println("provide either a lab name (--name) or a topology file path (--topo) or the flag --all")
			return nil
		}
		opts := []clab.ClabOption{
			clab.WithTimeout(timeout),
			clab.WithRuntime(rt,
				&runtime.RuntimeConfig{
					Debug:            debug,
					Timeout:          timeout,
					GracefulShutdown: graceful,
				},
			),
		}
		if topo != "" {
			opts = append(opts, clab.WithTopoFile(topo, varsFile))
		}
		c, err := clab.NewContainerLab(opts...)
		if err != nil {
			return fmt.Errorf("could not parse the topology file: %v", err)
		}

		if name == "" {
			name = c.Config.Name
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var glabels []*types.GenericFilter
		if all {
			glabels = []*types.GenericFilter{{FilterType: "label", Field: "containerlab", Operator: "exists"}}
		} else {
			if name != "" {
				glabels = []*types.GenericFilter{{FilterType: "label", Match: name, Field: "containerlab", Operator: "="}}
			} else if topo != "" {
				glabels = []*types.GenericFilter{{FilterType: "label", Match: c.Config.Name, Field: "containerlab", Operator: "="}}
			}
		}

		containers, err := c.ListContainers(ctx, glabels)
		if err != nil {
			return fmt.Errorf("failed to list containers: %s", err)
		}

		if len(containers) == 0 {
			log.Println("no containers found")
			return nil
		}
		if details {
			b, err := json.MarshalIndent(containers, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal containers struct: %v", err)
			}
			fmt.Println(string(b))
			return nil
		}

		err = printContainerInspect(c, containers, format)
		return err
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.Flags().BoolVarP(&details, "details", "", false, "print all details of lab containers")
	inspectCmd.Flags().StringVarP(&format, "format", "f", "table", "output format. One of [table, json]")
	inspectCmd.Flags().BoolVarP(&all, "all", "a", false, "show all deployed containerlab labs")
}

func toTableData(det []types.ContainerDetails) [][]string {
	tabData := make([][]string, 0, len(det))
	for i := range det {
		d := &det[i]

		if all {
			tabData = append(tabData, []string{fmt.Sprintf("%d", i+1), d.LabPath, d.LabName, d.Name, d.ContainerID, d.Image, d.Kind, d.State, d.IPv4Address, d.IPv6Address})
			continue
		}
		tabData = append(tabData, []string{fmt.Sprintf("%d", i+1), d.Name, d.ContainerID, d.Image, d.Kind, d.State, d.IPv4Address, d.IPv6Address})
	}
	return tabData
}

func printContainerInspect(c *clab.CLab, containers []types.GenericContainer, format string) error {
	contDetails := make([]types.ContainerDetails, 0, len(containers))
	// do not print published ports unless mysocketio kind is found
	printMysocket := false
	var mysocketCID string

	for i := range containers {
		cont := &containers[i]
		// get topo file path relative of the cwd
		cwd, _ := os.Getwd()
		path, _ := filepath.Rel(cwd, cont.Labels["clab-topo-file"])

		cdet := &types.ContainerDetails{
			LabName:     cont.Labels["containerlab"],
			LabPath:     path,
			Image:       cont.Image,
			State:       cont.State,
			IPv4Address: cont.GetContainerIPv4(),
			IPv6Address: cont.GetContainerIPv6(),
		}
		cdet.ContainerID = cont.ShortID

		if len(cont.Names) > 0 {
			cdet.Name = strings.TrimLeft(cont.Names[0], "/")
		}
		if kind, ok := cont.Labels["clab-node-kind"]; ok {
			cdet.Kind = kind
			if kind == "mysocketio" {
				printMysocket = true
				mysocketCID = cont.ID
			}
		}
		if group, ok := cont.Labels["clab-node-group"]; ok {
			cdet.Group = group
		}
		contDetails = append(contDetails, *cdet)
	}

	sort.Slice(contDetails, func(i, j int) bool {
		if contDetails[i].LabName == contDetails[j].LabName {
			return contDetails[i].Name < contDetails[j].Name
		}
		return contDetails[i].LabName < contDetails[j].LabName
	})

	if format == "json" {
		b, err := json.MarshalIndent(contDetails, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal container details: %v", err)
		}
		fmt.Println(string(b))
		return nil
	}
	tabData := toTableData(contDetails)
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{
		"Lab Name",
		"Name",
		"Container ID",
		"Image",
		"Kind",
		"State",
		"IPv4 Address",
		"IPv6 Address",
	}
	if all {
		table.SetHeader(append([]string{"#", "Topo Path"}, header...))
	} else {
		table.SetHeader(append([]string{"#"}, header[1:]...))
	}
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	// merge cells with lab name and topo file path
	table.SetAutoMergeCellsByColumnIndex([]int{1, 2})
	table.AppendBulk(tabData)
	table.Render()

	if !printMysocket {
		return nil
	}

	runtime := c.GlobalRuntime()

	stdout, stderr, err := runtime.Exec(context.Background(), mysocketCID, []string{"mysocketctl", "socket", "ls"})
	if err != nil {
		return fmt.Errorf("failed to execute cmd: %v", err)

	}
	if len(stderr) > 0 {
		log.Infof("errors during listing mysocketio sockets: %s", string(stderr))
	}

	fmt.Println("Published ports:")
	fmt.Println(string(stdout))

	return nil
}
