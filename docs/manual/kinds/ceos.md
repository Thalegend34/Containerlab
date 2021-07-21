# Arista cEOS

[Arista cEOS](https://www.arista.com/en/products/software-controlled-container-networking) is identified with `ceos` kind in the [topology file](../topo-def-file.md). The `ceos` kind defines a supported feature set and a startup procedure of a `ceos` node.

cEOS nodes launched with containerlab comes up with

* their management interface `eth0` configured with IPv4/6 addresses as assigned by docker
* hostname assigned to the node name
* gNMI, Netconf and eAPI services enabled
* `admin` user created with password `admin`

## Managing ceos nodes
Arista cEOS node launched with containerlab can be managed via the following interfaces:

=== "bash"
    to connect to a `bash` shell of a running ceos container:
    ```bash
    docker exec -it <container-name/id> bash
    ```
=== "CLI"
    to connect to the ceos CLI
    ```bash
    docker exec -it <container-name/id> Cli
    ```
=== "NETCONF"
    NETCONF server is running over port 830
    ```bash
    ssh root@<container-name> -p 830 -s netconf
    ```
=== "gNMI"
    gNMI server is running over port 6030 in non-secure mode
    using the best in class [gnmic](https://gnmic.kmrd.dev) gNMI client as an example:
    ```bash
    gnmic -a <container-name/node-mgmt-address>:6030 --insecure \
    -u admin -p admin \
    capabilities
    ```

!!!info
    Default user credentials: `admin:admin`

## Interfaces mapping
ceos container uses the following mapping for its linux interfaces:

* `eth0` - management interface connected to the containerlab management network
* `eth1` - first data interface

When containerlab launches ceos node, it will set IPv4/6 addresses as assigned by docker to the `eth0` interface and ceos node will boot with that addresses configured. Data interfaces `eth1+` need to be configured with IP addressing manually.

???note "ceos interfaces output"
    This output demonstrates the IP addressing of the linux interfaces of ceos node.
    ```
    bash-4.2# ip address
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/24 scope host lo
        valid_lft forever preferred_lft forever
        inet6 ::1/128 scope host
        valid_lft forever preferred_lft forever
    <SNIP>
    5877: eth0@if5878: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
        link/ether 02:42:ac:14:14:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
        inet 172.20.20.2/24 brd 172.20.20.255 scope global eth0
        valid_lft forever preferred_lft forever
        inet6 2001:172:20:20::2/80 scope global
        valid_lft forever preferred_lft forever
        inet6 fe80::42:acff:fe14:1402/64 scope link
        valid_lft forever preferred_lft forever
    ```
    This output shows how the linux interfaces are mapped into the ceos OS.
    ```
    ceos>sh ip int br
                                                                                Address
    Interface         IP Address           Status       Protocol           MTU    Owner
    ----------------- -------------------- ------------ -------------- ---------- -------
    Management0       172.20.20.2/24       up           up                1500

    ceos>sh ipv6 int br
    Interface       Status        MTU       IPv6 Address                     Addr State    Addr Source
    --------------- ------------ ---------- -------------------------------- ---------------- -----------
    Ma0             up           1500       fe80::42:acff:fe14:1402/64       up            link local
                                            2001:172:20:20::2/80             up            config
    ```
    As you see, the management interface `Ma0` inherits the IP address that docker assigned to ceos container management interface.


## Additional interface naming considerations

While many users will be fine with the default ceos naming of `eth`, some ceos users may find that they need to name their interfaces `et`. Interfaces named `et` provide consistency with the underlying interface mappings within ceos. This enables the correct operation of commands/features which depend on `et` format interface naming.

In order to align interfaces in this manner, the `INTFTYPE` environment variable must be set to `et` in the topology definition file and the links which are defined must be named `et`, as opposed to `eth`. This naming requirement does not apply to the `eth0` interface automatically created by containerlab. This is only required for links that are used for interconnection with other elements in a topology.

example:
```yml
topology:
  defaults:
    env:
      INTFTYPE: et
  nodes:
  { ... snipped misc. node definition for brevity ... }
  links:
    - endpoints: ["ceos_rtr1:et1", "ceos_rtr2:et1"]
    - endpoints: ["ceos_rtr1:et2", "ceos_rtr3:et1"]
```

## Features and options
### Node configuration
cEOS nodes have a dedicated [`config`](../conf-artifacts.md#identifying-a-lab-directory) directory that is used to persist the configuration of the node. It is possible to launch nodes of `ceos` kind with a basic config or to provide a custom config file that will be 
used as a startup config instead.

#### Default node configuration
When a node is defined without `config` statement present, containerlab will generate an empty config from [this template](https://github.com/srl-labs/containerlab/blob/master/nodes/ceos/ceos.cfg) and copy it to the config directory of the node.

```yaml
# example of a topo file that does not define a custom config
# as a result, the config will be generated from a template
# and used by this node
name: ceos
topology:
  nodes:
    ceos:
      kind: ceos
```

The generated config will be saved by the path `clab-<lab_name>/<node-name>/flash/startup-config`. Using the example topology presented above, the exact path to the config will be `clab-ceos/ceos/flash/startup-config`.

cEOS Ma0 interface will be configured with a random MAC address with `00:1c:73` OUI part. Containerlab will also create a `system_mac_address` file in the node's lab directory with the value of a System MAC address. The System MAC address value is calculated as `Ma0-MAC-addr + 1`.

#### User defined config
It is possible to make ceos nodes to boot up with a user-defined config instead of a built-in one. With a [`startup-config`](../nodes.md#startup-config) property a user sets the path to the config file that will be mounted to a container and used as a startup config:

```yaml
name: ceos_lab
topology:
  nodes:
    ceos:
      kind: ceos
      config: myconfig.conf
```

When a config file is passed via `config` parameter, it will override any configuration that may have left upon lab destroy.

With such topology file containerlab is instructed to take a file `myconfig.conf` from the current working directory, copy it to the lab directory for that specific node under the `/flash/startup-config` name and mount that dir to the container. This will result in this config to act as a startup config for the node.

It is possible to change the default config which every ceos node will start with with the following steps:

1. Save the [default configuration template](https://github.com/srl-labs/containerlab/blob/master/nodes/ceos/ceos.cfg) under some local file name[^2] and add the necessary changes to it
2. Use this file as a startup-config for ceos kind:
    ```
    name: ceos

    topology:
    kinds:
        ceos:
        startup-config: ceos-custom-startup.cfg
    nodes:
        # ceos1 will boot with ceos-custom-startup.cfg as set in the kind parameters
        ceos1:
        kind: ceos
        image: ceos:4.25.0F
        # ceos2 will boot with its own specific startup config, as it overrides the kind variables
        ceos2: 
        kind: ceos
        image: ceos:4.25.0F
        startup-config: node-specific-startup.cfg
    links:
        - endpoints: ["ceos1:eth1", "ceos2:eth1"]
    ```

#### Configuration persistency

It is important to understand how configuration persistency behaves when a single lab is going through rounds of `deploy->destroy` actions.

When the lab with cEOS nodes gets deployed for the first time the configuration file is generated with the IPv4/6 address assigned to `Ma0` management interface. These management interface addresses match the IP addresses that docker has assigned to cEOS containers. This makes it possible to have the cEOS nodes to start up with Management interface already correctly addressed.

When a user later configures the nodes during the lab exercise and saves it with `wr mem` or similar, the changes will be written to `startup-config` file of cEOS.

User then may destroy the lab and the config changes will persist on disk, this is done with `destroy` command. During this operation the containers will be destroyed, but their configuration files will still be kept in the lab directory by the path `clab-$labName`.

If a user then desires to start this lab once again it may lead to a problem. Since docker may assign new IP addresses to the cEOS nodes of the lab, the configuration saved on disk may not match those new docker-assigned addresses, and that will result in an incorrect management interface configuration.

To avoid this, and be able to start the nodes with the previously saved configuration, users may do the following:

1. Address the nodes explicitly via [user defined addresses](../network.md#user-defined-addresses). This will instruct docker to use the addresses as specified by a user in a clab file.
2. Leverage [user defined config](#user-defined-config), if all you need is to have a startup config.

#### Saving configuration
In addition to cli commands such as `write memory` user can take advantage of the [`containerlab save`](../../cmd/save.md) command. It saves running cEOS configuration into a file by `conf-saved.conf` path in the relevant node directory.

## Container configuration
To start an Arista cEOS node containerlab uses the configuration instructions described in Arista Forums[^1]. The exact parameters are outlined below.

=== "Startup command"
    `/sbin/init systemd.setenv=INTFTYPE=eth systemd.setenv=ETBA=4 systemd.setenv=SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT=1 systemd.setenv=CEOS=1 systemd.setenv=EOS_PLATFORM=ceoslab systemd.setenv=container=docker systemd.setenv=MAPETH0=1 systemd.setenv=MGMT_INTF=eth0`
=== "Environment variables"
    `CEOS:1`  
    `EOS_PLATFORM":ceoslab`  
    `container:docker`  
    `ETBA:4`  
    `SKIP_ZEROTOUCH_BARRIER_IN_SYSDBINIT:1`  
    `INTFTYPE:eth`  
    `MAPETH0:1`  
    `MGMT_INTF:eth0`

### File mounts
When a user starts a lab, containerlab creates a node directory for storing [configuration artifacts](../conf-artifacts.md). For `ceos` kind containerlab creates `flash` directory for each ceos node and mounts these folders by `/mnt/flash` paths.

```
❯ tree clab-srlceos01/ceos
clab-srlceos01/ceos
└── flash
    ├── AsuFastPktTransmit.log
    ├── debug
    │   └── proc
    │       └── modules
    ├── fastpkttx.backup
    ├── Fossil
    ├── kickstart-config
    ├── persist
    │   ├── local
    │   ├── messages
    │   ├── persistentRestartLog
    │   ├── secure
    │   └── sys
    ├── schedule
    │   └── tech-support
    │       └── ceos_tech-support_2021-01-14.0907.log.gz
    ├── SsuRestoreLegacy.log
    ├── SsuRestore.log
    ├── system_mac_address
    └── startup-config

9 directories, 11 files
```

## Lab examples
The following labs feature cEOS node:

- [SR Linux and cEOS](../../lab-examples/srl-ceos.md)

## Known issues or limitations
### cgroups v1

As of this writing (22-June, 2021), ceos-lab image requires a cgroups v1 environment.  For many users this should not require any changes to the runtime environment.  However, some linux distributions (ref: [#467](https://github.com/srl-labs/containerlab/issues/467)) may be configured to use cgroups v2 out-of-the-box, which will prevent ceos-lab image from booting. In such cases, the users will need to configure their system to utilize a cgroups v1 environment.  

Consult your distribution's documentation for details regarding configuring cgroups v1 in case you see similar startup issues as indicated in [#467](https://github.com/srl-labs/containerlab/issues/467). 

[^1]: https://eos.arista.com/ceos-lab-topo/
[^2]: do not remove the template variables from the `Management0` interface, otherwise the nodes will not apply the IP address from docker IPAM service.