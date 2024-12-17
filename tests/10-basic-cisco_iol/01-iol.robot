*** Settings ***
Library             OperatingSystem
Resource            ../ssh.robot
Resource            ../common.robot

Suite Teardown      Run Keyword    Cleanup


*** Variables ***
${lab-name}         iol
${lab-file-name}    iol.clab.yml
${runtime}          docker


*** Test Cases ***
Deploy ${lab-name} lab
    Log    ${CURDIR}
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sudo -E ${CLAB_BIN} --runtime ${runtime} deploy -t ${CURDIR}/${lab-file-name}
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0

Wait 45s for nodes to boot
    Sleep    45s

Verify links in node router1
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-router1 sh ip int br | head -5
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    172.20.20.
    Should Contain    ${output}    up

Verify links in node switch
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-switch sh ip int br | head -5 | tail -1
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    172.20.20.
    Should Contain    ${output}    up

Verify parital startup configuration on router2
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-router2 show running-config interface Loopback0
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    PARTIAL_CFG

Verify full startup configuration on router3
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-router3 "sh run | inc hostname"
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    FULL_STARTUP_CFG-router3

Write configuration to NVRAM on router1
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-iol-router1 "write memory"
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    [OK]

Log IP addresses for router1
    ${rc}   ${ipv4_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq '.nodes.router1."mgmt-ipv4-address"'
    Log    ${ipv4_addr}
    Should Be Equal As Integers    ${rc}    0
    ${rc}   ${ipv6_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq '.nodes.router1."mgmt-ipv6-address"'
    Log    ${ipv6_addr}

Write configuration to NVRAM on switch
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-iol-switch "write memory"
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    [OK]

Log IP addresses for switch
    ${rc}   ${ipv4_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq '.nodes.switch."mgmt-ipv4-address"'
    Log    ${ipv4_addr}
    Should Be Equal As Integers    ${rc}    0
    ${rc}   ${ipv6_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq '.nodes.switch."mgmt-ipv6-address"'
    Log    ${ipv6_addr}

Destroy ${lab-name} lab
    Log     ${CURDIR}
    ${rc}   ${output} =    Run And Return Rc And Output
    ...     sudo -E ${CLAB_BIN} --runtime ${runtime} destroy -t ${CURDIR}/${lab-file-name}
    Log     ${output}
    Should Be Equal As Integers    ${rc}    0

Re-deploy ${lab-name} lab
    Log     ${CURDIR}
    ${rc}   ${output} =    Run And Return Rc And Output
    ...     sudo -E ${CLAB_BIN} --runtime ${runtime} deploy -t ${CURDIR}/${lab-file-name}
    Log     ${output}
    Should Be Equal As Integers    ${rc}    0

Wait 60s for nodes to boot
    Sleep    60s

Verify connectivity via new management addresses on router1
    ${rc}   ${ipv4_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq -r '.nodes.router1."mgmt-ipv4-address"'
    Should Be Equal As Integers    ${rc}    0
    ${rc}   ${ipv6_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq -r '.nodes.router1."mgmt-ipv6-address"'
    Should Be Equal As Integers    ${rc}    0
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-router1 "sh run interface Ethernet0/0"
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    ${ipv4_addr.upper()}
    Should Contain    ${output}    ${ipv6_addr.upper()}

Verify connectivity via new management addresses on switch
    ${rc}   ${ipv4_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq -r '.nodes.switch."mgmt-ipv4-address"'
    Should Be Equal As Integers    ${rc}    0
    ${rc}   ${ipv6_addr} =    Run And Return Rc And Output
    ...     cat ${CURDIR}/clab-${lab-name}/topology-data.json | jq -r '.nodes.switch."mgmt-ipv6-address"'
    Should Be Equal As Integers    ${rc}    0
    ${rc}    ${output} =    Run And Return Rc And Output
    ...    sshpass -p "admin" ssh -o "IdentitiesOnly=yes" admin@clab-${lab-name}-switch "sh run interface Ethernet0/0"
    Log    ${output}
    Should Be Equal As Integers    ${rc}    0
    Should Contain    ${output}    ${ipv4_addr.upper()}
    Should Contain    ${output}    ${ipv6_addr.upper()}

*** Keywords ***
Cleanup
    Run    sudo -E ${CLAB_BIN} --runtime ${runtime} destroy -t ${CURDIR}/${lab-file-name} --cleanup
