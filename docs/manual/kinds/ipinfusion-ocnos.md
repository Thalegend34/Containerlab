# IPInfusion OcNOS

IPInfusion OcNOS virtualized router is identified with `ipinfusion-ocnos` kind in the [topology file](../topo-def-file.md). It is built using [boxen](https://github.com/carlmontanari/boxen) project and essentially is a Qemu VM packaged in a docker container format.

ipinfusion-ocnos nodes launched with containerlab come up pre-provisioned with SSH, and NETCONF services enabled.

## Managing ipinfusion-ocnos nodes

!!!note
    Containers with OcNOS inside will take ~3min to fully boot.  
    You can monitor the progress with `docker logs -f <container-name>` and `docker exec -it <container-name> tail -f /console.log`.

IPInfusion OcNOS node launched with containerlab can be managed via the following interfaces:

=== "bash"
    to connect to a `bash` shell of a running ipinfusion-ocnos container:
    ```bash
    docker exec -it <container-name/id> bash
    ```
=== "CLI"
    to connect to the OcNOS CLI
    ```bash
    ssh ocnos@<container-name/id>
    ```
=== "NETCONF"
    NETCONF server is running over port 830
    ```bash
    ssh ocnos@<container-name> -p 830 -s netconf
    ```

!!!info
    Default user credentials: `ocnos:ocnos`

## Interfaces mapping
ipinfusion-ocnos container can have up to 144 interfaces and uses the following mapping rules:

* `eth0` - management interface connected to the containerlab management network
* `eth1` - first data interface, mapped to first data port of OcNOS line card
* `eth2+` - second and subsequent data interface

When containerlab launches ipinfusion-ocnos node, it will assign IPv4 address to the `eth0` interface. This address can be used to reach management plane of the router.

Data interfaces `eth1+` need to be configured with IP addressing manually using CLI/management protocols.
