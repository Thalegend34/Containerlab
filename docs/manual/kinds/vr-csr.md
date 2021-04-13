# Cisco CSR1000v

Cisco CSR1000v virtualized router is identified with `vr-csr` kind in the [topology file](../topo-def-file.md). It is built using [vrnetlab](../vrnetlab.md) project and essentially is a Qemu VM packaged in a docker container format.

vr-csr nodes launched with containerlab comes up pre-provisioned with SSH, SNMP, NETCONF and gNMI services enabled.

## Managing vr-csr nodes

!!!note
    Containers with CSR1000v inside will take ~6min to fully boot.  
    You can monitor the progress with `docker logs -f <container-name>`.

Cisco CSR1000v node launched with containerlab can be managed via the following interfaces:

=== "bash"
    to connect to a `bash` shell of a running vr-csr container:
    ```bash
    docker exec -it <container-name/id> bash
    ```
=== "CLI"
    to connect to the CSR1000v CLI
    ```bash
    ssh admin@<container-name/id>
    ```
=== "NETCONF"
    NETCONF server is running over port 830
    ```bash
    ssh admin@<container-name> -p 830 -s netconf
    ```

!!!info
    Default user credentials: `admin:admin`

## Interfaces mapping
vr-csr container can have up to 144 interfaces and uses the following mapping rules:

* `eth0` - management interface connected to the containerlab management network
* `eth1` - first data interface, mapped to first data port of CSR1000v line card
* `eth2+` - second and subsequent data interface

When containerlab launches vr-csr node, it will assign IPv4/6 address to the `eth0` interface. These addresses can be used to reach management plane of the router.

Data interfaces `eth1+` needs to be configured with IP addressing manually using CLI/management protocols.


## Features and options
### Node configuration
vr-csr nodes come up with a basic configuration where only `admin` user and management interfaces such as NETCONF provisioned.
