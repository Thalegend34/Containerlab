// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package clab

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/containers/podman/v4/pkg/util"
	"github.com/google/go-cmp/cmp"
	"github.com/srl-labs/containerlab/labels"
	"github.com/srl-labs/containerlab/utils"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	// setup function
	// create node' labdir with a file that is used in topo2 definition for binds resolver to check path existence
	f, _ := filepath.Abs("clab-topo2/node1/somefile")
	os.MkdirAll("clab-topo2/node1", 0777)

	if _, err := os.Create(f); err != nil {
		t.Error(err)
	}
	// teardown function
	return func(t *testing.T) {
		os.RemoveAll("clab-topo2")
	}
}

func TestLicenseInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		want string
	}{
		"node_license": {
			got:  "test_data/topo1.yml",
			want: "node1.lic",
		},
		"kind_license": {
			got:  "test_data/topo2.yml",
			want: "kind.lic",
		},
		"default_license": {
			got:  "test_data/topo3.yml",
			want: "default.lic",
		},
		"kind_overwrite": {
			got:  "test_data/topo4.yml",
			want: "node1.lic",
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			// fmt.Println(c.Config.Topology.Defaults)
			// fmt.Println(c.Config.Topology.Kinds)
			// fmt.Println(c.Config.Topology.Nodes)
			if filepath.Base(c.Nodes["node1"].Config().License) != tc.want {
				t.Fatalf("wanted '%s' got '%s'", tc.want, c.Nodes["node1"].Config().License)
			}
		})
	}
}

func TestBindsInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		want []string
	}{
		"node_single_bind": {
			got:  "test_data/topo1.yml",
			want: []string{"node1.lic:/dst"},
		},
		"node_many_binds": {
			got: "test_data/topo2.yml",
			want: []string{
				"node1.lic:/dst1",
				"kind.lic:/dst2",
				"${PWD}/clab-topo2/node1/somefile:/somefile",
			},
		},
		"kind_and_node_binds": {
			got:  "test_data/topo5.yml",
			want: []string{"kind.lic:/dst", "node1.lic:/dst2"},
		},
		"default_binds": {
			got:  "test_data/topo3.yml",
			want: []string{"default.lic:/dst"},
		},
		"default_and_kind_and_node_binds": {
			got:  "test_data/topo4.yml",
			want: []string{"node1.lic:/dst1", "kind.lic:/dst2", "default.lic:/dst3"},
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Error(err)
			}

			binds := c.Nodes["node1"].Config().Binds

			// expand env vars in bind paths, this is done during topology file load by clab
			utils.ExpandEnvVarsInStrSlice(tc.want)

			// resolve wanted paths as the binds paths are resolved as part of the c.ParseTopology
			err = c.resolveBindPaths(tc.want, c.Nodes["node1"].Config().LabDir)
			if err != nil {
				t.Error(err)
			}

			for _, b := range tc.want {
				if !util.StringInSlice(b, binds) {
					t.Errorf("bind %q is not found in resulting binds %q", b, binds)
				}
			}
		})
	}
}

func TestTypeInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		node string
		want string
	}{
		"undefined_type_returns_default": {
			got:  "test_data/topo1.yml",
			node: "node2",
			want: "ixrd2",
		},
		"node_type_override_kind_type": {
			got:  "test_data/topo2.yml",
			node: "node2",
			want: "ixr10",
		},
		"node_inherits_kind_type": {
			got:  "test_data/topo2.yml",
			node: "node1",
			want: "ixrd2",
		},
		"node_inherits_default_type": {
			got:  "test_data/topo3.yml",
			node: "node2",
			want: "ixrd2",
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(c.Nodes[tc.node].Config().NodeType, tc.want) {
				t.Fatalf("wanted %q got %q", tc.want, c.Nodes[tc.node].Config().NodeType)
			}
		})
	}
}

func TestEnvInit(t *testing.T) {
	tests := map[string]struct {
		got    string
		node   string
		envvar map[string]string
		want   map[string]string
	}{
		"env_defined_at_node_level": {
			got:  "test_data/topo1.yml",
			node: "node1",
			want: map[string]string{
				"env1": "val1",
				"env2": "val2",
			},
		},
		"env_defined_at_kind_level": {
			got:  "test_data/topo2.yml",
			node: "node2",
			want: map[string]string{
				"env1": "val1",
			},
		},
		"env_defined_at_defaults_level": {
			got:  "test_data/topo3.yml",
			node: "node1",
			want: map[string]string{
				"env1": "val1",
			},
		},
		"env_defined_at_node_and_kind_and_default_level": {
			got:  "test_data/topo4.yml",
			node: "node1",
			want: map[string]string{
				"env1": "node",
				"env2": "kind",
				"env3": "global",
				"env4": "kind",
				"env5": "node",
			},
		},
		"expand_env_variables": {
			got:  "test_data/topo9.yml",
			node: "node1",
			envvar: map[string]string{
				"CONTAINERLAB_TEST_ENV5": "node",
			},
			want: map[string]string{
				"env1": "node",
				"env2": "kind",
				"env3": "global",
				"env4": "kind",
				"env5": "node",
			},
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.envvar {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			// nodeCfg := c.Config.Topology.Nodes[tc.node]
			// kind := strings.ToLower(c.kindInitialization(nodeCfg))
			env := c.Config.Topology.GetNodeEnv(tc.node)
			// env := c.envInit(nodeCfg, kind)
			if !reflect.DeepEqual(env, tc.want) {
				t.Fatalf("wanted %q got %q", tc.want, env)
			}
		})
	}
}

func TestUserInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		node string
		want string
	}{
		"user_defined_at_node_level": {
			got:  "test_data/topo1.yml",
			node: "node2",
			want: "custom",
		},
		"user_defined_at_kind_level": {
			got:  "test_data/topo2.yml",
			node: "node2",
			want: "customkind",
		},
		"user_defined_at_defaults_level": {
			got:  "test_data/topo3.yml",
			node: "node1",
			want: "customglobal",
		},
		"user_defined_at_node_and_kind_and_default_level": {
			got:  "test_data/topo4.yml",
			node: "node1",
			want: "customnode",
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			// nodeCfg := c.Config.Topology.Nodes[tc.node]
			// kind := strings.ToLower(c.kindInitialization(nodeCfg))
			user := c.Config.Topology.GetNodeUser(tc.node)
			// user := c.userInit(nodeCfg, kind)
			if user != tc.want {
				t.Fatalf("wanted %q got %q", tc.want, user)
			}
		})
	}
}

func TestVerifyLinks(t *testing.T) {
	tests := map[string]struct {
		got  string
		want string
	}{
		"two_duplicated_links": {
			got:  "test_data/topo6.yml",
			want: "endpoints [\"lin1:eth1\" \"lin2:eth2\"] appeared more than once in the links section of the topology file",
		},
		"no_duplicated_links": {
			got:  "test_data/topo1.yml",
			want: "",
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			err = c.verifyLinks()
			if err != nil && err.Error() != tc.want {
				t.Fatalf("wanted %q got %q", tc.want, err.Error())
			}
			if err == nil && tc.want != "" {
				t.Fatalf("wanted %q got %q", tc.want, err.Error())
			}
		})
	}
}

func TestLabelsInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		node string
		want map[string]string
	}{
		"only_default_labels": {
			got:  "test_data/topo1.yml",
			node: "node1",
			want: map[string]string{
				labels.Containerlab: "topo1",
				labels.NodeName:     "node1",
				labels.NodeKind:     "srl",
				labels.NodeType:     "ixrd2",
				labels.NodeGroup:    "",
				labels.NodeLabDir:   "../clab-topo1/node1",
				labels.TopoFile:     "topo1.yml",
			},
		},
		"custom_node_label": {
			got:  "test_data/topo1.yml",
			node: "node2",
			want: map[string]string{
				labels.Containerlab: "topo1",
				labels.NodeName:     "node2",
				labels.NodeKind:     "srl",
				labels.NodeType:     "ixrd2",
				labels.NodeGroup:    "",
				labels.NodeLabDir:   "../clab-topo1/node2",
				labels.TopoFile:     "topo1.yml",
				"node-label":        "value",
			},
		},
		"custom_kind_label": {
			got:  "test_data/topo2.yml",
			node: "node1",
			want: map[string]string{
				labels.Containerlab: "topo2",
				labels.NodeName:     "node1",
				labels.NodeKind:     "srl",
				labels.NodeType:     "ixrd2",
				labels.NodeGroup:    "",
				labels.NodeLabDir:   "../clab-topo2/node1",
				labels.TopoFile:     "topo2.yml",
				"kind-label":        "value",
			},
		},
		"custom_default_label": {
			got:  "test_data/topo3.yml",
			node: "node2",
			want: map[string]string{
				labels.Containerlab: "topo3",
				labels.NodeName:     "node2",
				labels.NodeKind:     "srl",
				labels.NodeType:     "ixrd2",
				labels.NodeGroup:    "",
				labels.NodeLabDir:   "../clab-topo3/node2",
				labels.TopoFile:     "topo3.yml",
				"default-label":     "value",
			},
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			tc.want[labels.NodeLabDir] = utils.ResolvePath(tc.want[labels.NodeLabDir], c.TopoPaths.TopologyFileDir())
			tc.want[labels.TopoFile] = utils.ResolvePath(tc.want[labels.TopoFile], c.TopoPaths.TopologyFileDir())

			labels := c.Nodes[tc.node].Config().Labels

			if !cmp.Equal(labels, tc.want) {
				t.Errorf("failed at '%s', expected\n%v, got\n%+v", name, tc.want, labels)
			}

			// test that labels were propagated to env vars as CLAB_LABEL_<label-name>:<label-value>
			env := c.Nodes[tc.node].Config().Env
			fmt.Printf("%v\n", env)
			for k, v := range tc.want {
				// sanitize label key to be used as an env key
				sk := utils.ToEnvKey(k)
				// fail if env vars map doesn't have env var with key CLAB_LABEL_<label-name> and label value matches env value
				if val, exists := env["CLAB_LABEL_"+sk]; !exists || val != v {
					t.Errorf("env var %q promoted from a label %q was not found", "CLAB_LABEL_"+sk, k)
				}
			}
		})
	}
}

func TestVerifyRootNetnsInterfaceUniqueness(t *testing.T) {
	opts := []ClabOption{
		WithTopoFile("test_data/topo7-dup-rootnetns.yml", "", []string{}),
	}
	c, err := NewContainerLab(opts...)
	if err != nil {
		t.Fatal(err)
	}

	err = c.verifyRootNetnsInterfaceUniqueness()
	if err == nil {
		t.Fatalf("expected duplicate rootns links error")
	}
	t.Logf("error: %v", err)
}

func TestEnvFileInit(t *testing.T) {
	tests := map[string]struct {
		got  string
		node string
		want map[string]string
	}{
		"env-file_defined_at_node_and_default_1": {
			got:  "test_data/topo10.yml",
			node: "node1",
			want: map[string]string{
				"env1":     "val1",
				"env2":     "val2",
				"ENVFILE1": "SOMEOTHERDATA",
				"ENVFILE2": "THISANDTHAT",
			},
		},
		"env-file_defined_at_node_and_default_2": {
			got:  "test_data/topo10.yml",
			node: "node2",
			want: map[string]string{
				"ENVFILE1": "SOMEENVVARDATA",
				"ENVFILE2": "THISANDTHAT",
			},
		},
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := []ClabOption{
				WithTopoFile(tc.got, "", []string{}),
			}
			c, err := NewContainerLab(opts...)
			if err != nil {
				t.Fatal(err)
			}

			env := c.Nodes[tc.node].Config().Env
			// check all the want key/values are there
			for k, v := range tc.want {
				// check keys defined in tc.want exist and values are equal
				if val, exists := env[k]; !(exists && val == v) {
					t.Fatalf("wanted %q to be contained in env, but got %q", tc.want, env)
				}
			}
		})
	}
}
