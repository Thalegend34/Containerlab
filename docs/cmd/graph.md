# graph command

### Description

The `graph` command generates graphical representations of a topology.

Two graphing options are available:

* an HTML page served by `containerlab` web-server based on a user-provided HTML template and static files.
* a [graph description file in dot format](https://en.wikipedia.org/wiki/DOT_(graph_description_language)) that can be rendered using [Graphviz](https://graphviz.org/) or viewed [online](https://dreampuf.github.io/GraphvizOnline/).[^1]

#### HTML

The HTML-based graph representation is the default graphing option. The topology will be graphed and served online using the embedded web server.

The default graph template is based on the [NeXt UI](https://developer.cisco.com/site/neXt/) framework[^2].

![animation](https://user-images.githubusercontent.com/11521160/155654224-d46b346d-7051-49f8-ba93-6dee6d22a39f.gif)

To render a topology using this default graph engine:

```
containerlab graph -t <path/to/topo.clab.yml>
```

##### NeXt UI
Topology graph created with NeXt UI has some control elements that allow you to choose the color theme of the web view, scaling and panning. Besides these generic controls it is possible to enable auto-layout of the components using buttons at the top of the screen.

###### Layout and sorting
The graph engine can automatically pan and sort elements in your topology based on their _role_. We encode the role via `group` property of a node.

Today we have the following sort orders available to users:

```yaml
sortOrder: ['10', '9', 'superspine', '8', 'dc-gw', '7', '6', 'spine', '5', '4', 'leaf', 'border-leaf', '3', 'server', '2', '1'],
```
The values are sorted so that `10` is placed higher in the hierarchy than `9` and so on.

Consider the following snippet:

```yaml
topology:
  nodes:
    ### SPINES ###
    spine1:
      group: spine
    
    ### LEAFS ###
    leaf1:
      group: leaf

    ### CLIENTS ###
    client1:
      kind: linux
      group: server
```

The `group` property set to the predefined value will automatically auto-align the elements based on their role.

#### Graphviz

When `graph` command is called without the `--srv` flag, containerlab will generate a [graph description file in dot format](https://en.wikipedia.org/wiki/DOT_(graph_description_language)).

The dot file can be used to view the graphical representation of the topology either by rendering the dot file into a PNG file or using [online dot viewer](https://dreampuf.github.io/GraphvizOnline/).

### Online vs offline graphing
When HTML graph option is used, containerlab will try to build the topology graph by inspecting the running containers which are part of the lab. This essentially means, that the lab must be running. Although this method provides some additional details (like IP addresses), it is not always convenient to run a lab to see its graph.

The other option is to use the topology file solely to build the graph. This is done by adding `--offline` flag.

If `--offline` flag was not provided and no containers were found matching the lab name, containerlab will use the topo file only (as if offline mode was set).
### Usage

`containerlab [global-flags] graph [local-flags]`

### Flags

#### topology

With the global `--topo | -t` flag a user sets the path to the topology file that will be used to get the nodes of a lab.

#### srv

The `--srv` flag allows a user to customize the HTTP address and port for the web server. Default value is `:50080`.

A single path `/` is served, where the graph is generated based on either a default template or on the template supplied using `--template`.

#### template

The `--template` flag allows to customize the HTML based graph by supplying a user defined template that will be rendered and exposed on the address specified by `--srv`.

#### static-dir

The `--static-dir` flag enables the embedded HTML web-server to serve static files from the specified directory. Must be used together with the `--template` flag.

With this flag, it is possible to link to local files (JS, CSS, fonts, etc.) from the custom HTML template.

#### dot
With `--dot` flag provided containerlab will generate the `dot` file instead of serving the topology with embedded HTTP server.

### Examples

```bash

# render a graph from running lab or topo file if lab is not running#
# using HTML graph option with default server address :50080
containerlab graph --topo /path/to/topo1.clab.yml

# start an http server on :3002 where topo1 graph will be rendered using a custom template my_template.html
containerlab graph --topo /path/to/topo1.clab.yml --srv ":3002" --template my_template.html

# start an http server on default port :50080 
# using a custom template that links to local files located at /path/to/static_files directory
containerlab graph --topo /path/to/topo1.clab.yml --template my_template.html --static-dir /path/to/static_files
```

[^1]: This method is prone to errors when node names contain dashes and special symbols. Use with caution, and prefer the HTML server alternative.
[^2]: NeXt UI css/js files can be found at `/etc/containerlab/templates/graph/nextui` directory