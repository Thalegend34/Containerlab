|                               |                                                                                    |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| **Description**               | A Nokia SR Linux connected back-to-back with Juniper cRPD                          |
| **Components**                | [Nokia SR Linux][srl], [Juniper cRPD][crpd]                                        |
| **Resource requirements**[^1] | :fontawesome-solid-microchip: 1 <br/>:fontawesome-solid-memory: 2 GB               |
| **Topology file**             | [srlcrpd01.yml][topofile]                                                          |
| **Name**                      | srlcrpd01                                                                          |
| **Version information**[^2]   | `containerlab:0.9.0`, `srlinux:20.6.3-145`, `crpd:20.2R1.10`, `docker-ce:19.03.13` |

## Description
A lab consists of an SR Linux node connected with Juniper cRPD via a point-to-point ethernet link. Both nodes are also connected with their management interfaces to the `clab` docker network.

<div class="mxgraph" style="max-width:100%;border:1px solid transparent;margin:0 auto; display:block;" data-mxgraph="{&quot;page&quot;:0,&quot;zoom&quot;:1.5,&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;nav&quot;:true,&quot;check-visible-state&quot;:true,&quot;resize&quot;:true,&quot;url&quot;:&quot;https://raw.githubusercontent.com/srl-wim/container-lab/diagrams/srlcrpd01&quot;}"></div>

## Use cases
This lab allows users to launch basic interoperability scenarios between Nokia SR Linux and Juniper cRPD network operating systems.

Below you will find configuration instructions to setup IS-IS routing protocol on both nodes and verify the successful control and data plane operation.
### OSPF
<div class="mxgraph" style="max-width:100%;border:1px solid transparent;margin:0 auto; display:block;" data-mxgraph="{&quot;page&quot;:2,&quot;zoom&quot;:1.5,&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;nav&quot;:true,&quot;check-visible-state&quot;:true,&quot;resize&quot;:true,&quot;url&quot;:&quot;https://raw.githubusercontent.com/srl-wim/container-lab/diagrams/srlcrpd01&quot;}"></div>

#### Configuration
Once the lab is deployed with containerlab, use the following configuration instructions to make interfaces configuration and enable OSPF on both nodes.

=== "srl"
    Get into SR Linux CLI with `docker exec -it clab-srlcrpd01-srl sr_cli` and start configuration
    ```bash
    # enter candidate datastore
    enter candidate
    
    # configure loopback and data interfaces
    set / interface ethernet-1/1 admin-state enable
    set / interface ethernet-1/1 subinterface 0 admin-state enable
    set / interface ethernet-1/1 subinterface 0 ipv4 address 192.168.1.1/24

    set / interface lo0 subinterface 0 admin-state enable
    set / interface lo0 subinterface 0 ipv4 address 10.10.10.1/32

    # configure OSPF
    set / network-instance default router-id 10.10.10.1
    set / network-instance default interface ethernet-1/1.0
    set / network-instance default interface lo0.0
    set / network-instance default protocols ospf instance main admin-state enable
    set / network-instance default protocols ospf instance main version ospf-v2
    set / network-instance default protocols ospf instance main area 0.0.0.0 interface ethernet-1/1.0 interface-type point-to-point
    set / network-instance default protocols ospf instance main area 0.0.0.0 interface ethernet-1/1.0

    # commit config
    commit now
    ```
=== "crpd"
    cRPD configuration needs to be done both from the container process, as well as within the CLI.  
    First attach to the container process `bash` shell and configure interfaces: `docker exec -it clab-srlcrpd01-crpd bash`
    ```bash
    # configure linux interfaces
    ip addr add 192.168.1.2/24 dev eth1
    ip addr add 10.10.10.2/32 dev lo
    ```
    Then launch the CLI and continue configuration `docker exec -it clab-srlcrpd01-crpd cli`:
    ```bash
    # enter configuration mode
    configure
    set routing-options router-id 10.10.10.2

    set protocols ospf area 0.0.0.0 interface eth1 interface-type p2p
    set protocols ospf area 0.0.0.0 interface lo.0 interface-type nbma
    
    # commit configuration
    commit
    ```

#### Verificaton
After the configuration is done on both nodes, verify the control plane by checking the route tables on both ends and ensuring dataplane was programmed as well by pinging the remote loopback

=== "srl"
    ```bash
    # control plane verification
    A:srl# / show network-instance default route-table ipv4-unicast summary | grep ospf
    | 10.10.10.2/32                 | 0     | true       | ospfv2          | 1       | 10    | 192.168.1.2 (direct)                     | ethernet-1/1.0    |
    ```
    ```
    # data plane verification
    A:srl# ping 10.10.10.2 network-instance default
    Using network instance default
    PING 10.10.10.2 (10.10.10.2) 56(84) bytes of data.
    64 bytes from 10.10.10.2: icmp_seq=1 ttl=64 time=1.15 ms
    ```
=== "crpd"
    ```bash
    # control plane verification
    root@crpd> show route | match OSPF
    10.10.10.1/32      *[OSPF/10] 00:01:24, metric 1
    224.0.0.5/32       *[OSPF/10] 00:05:49, metric 1
    ```

### IS-IS
<div class="mxgraph" style="max-width:100%;border:1px solid transparent;margin:0 auto; display:block;" data-mxgraph="{&quot;page&quot;:1,&quot;zoom&quot;:1.5,&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;nav&quot;:true,&quot;check-visible-state&quot;:true,&quot;resize&quot;:true,&quot;url&quot;:&quot;https://raw.githubusercontent.com/srl-wim/container-lab/diagrams/srlcrpd01&quot;}"></div>

Once the lab is deployed with containerlab, use the following configuration instructions to make interfaces configuration and enable IS-IS on both nodes.

=== "srl"
    Get into SR Linux CLI with `docker exec -it clab-srlcrpd01-srl sr_cli` and start configuration
    ```bash
    # enter candidate datastore
    enter candidate
    
    # configure loopback and data interfaces
    set / interface ethernet-1/1 admin-state enable
    set / interface ethernet-1/1 subinterface 0 admin-state enable
    set / interface ethernet-1/1 subinterface 0 ipv4 address 192.168.1.1/24

    set / interface lo0 subinterface 0 admin-state enable
    set / interface lo0 subinterface 0 ipv4 address 10.10.10.1/32

    # configure IS-IS
    set / network-instance default router-id 10.10.10.1
    set / network-instance default interface ethernet-1/1.0
    set / network-instance default interface lo0.0
    set / network-instance default protocols isis instance main admin-state enable
    set / network-instance default protocols isis instance main level-capability L1
    set / network-instance default protocols isis instance main net [ 49.0001.0100.1001.0001.00 ]
    set / network-instance default protocols isis instance main interface ethernet-1/1.0 admin-state enable
    set / network-instance default protocols isis instance main interface ethernet-1/1.0 circuit-type point-to-point
    set / network-instance default protocols isis instance main interface ethernet-1/1.0 level 1
    set / network-instance default protocols isis instance main interface ethernet-1/1.0 level 2 disable true
    set / network-instance default protocols isis instance main level 1 metric-style wide

    # commit config
    commit now
    ```
=== "crpd"
    cRPD configuration needs to be done both from the container process, as well as within the CLI.  
    First attach to the container process `bash` shell and configure interfaces: `docker exec -it clab-srlcrpd01-crpd bash`
    ```bash
    # configure linux interfaces
    ip addr add 192.168.1.2/24 dev eth1
    ip addr add 10.10.10.2/32 dev lo
    ```
    Then launch the CLI and continue configuration `docker exec -it clab-srlcrpd01-crpd cli`:
    ```bash
    # enter configuration mode
    configure
    set interfaces lo0 unit 0 family iso address 49.0001.0100.1001.0002.00
    set routing-options router-id 10.10.10.2

    set protocols isis interface eth1 point-to-point
    set protocols isis interface eth1 level 2 disable
    set protocols isis interface lo0.0
    set protocols isis level 1 wide-metrics-only
    set protocols isis level 2 wide-metrics-only
    set protocols isis reference-bandwidth 100g
    
    # commit configuration
    commit
    ```

[srl]: https://www.nokia.com/networks/products/service-router-linux-NOS/
[crpd]: https://www.juniper.net/documentation/us/en/software/crpd/crpd-deployment/topics/concept/understanding-crpd.html
[topofile]: https://github.com/srl-wim/container-lab/tree/master/lab-examples/srlcrpd01/srlcrpd01.yml

[^1]: Resource requirements are provisional. Consult with the installation guides for additional information.
[^2]: The lab has been validated using these versions of the required tools/components. Using versions other than stated might lead to a non-operational setup process.

<script type="text/javascript" src="https://cdn.jsdelivr.net/gh/hellt/drawio-js@main/embed2.js" async></script>