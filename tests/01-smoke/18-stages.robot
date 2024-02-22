*** Settings ***
Library             OperatingSystem
Library             String
Library             Process
Resource            ../common.robot

Suite Setup         Setup
Suite Teardown      Run Keyword    Teardown


*** Variables ***
${lab-file}     stages.clab.yml
${lab-name}     stages
${runtime}      docker


*** Test Cases ***
Pre-Pull Image
    ${output} =    Process.Run Process
    ...    ${runtime} pull debian:bookworm-slim
    ...    shell=True
    Log    ${output.stdout}
    Log    ${output.stderr}

    Should Be Equal As Integers    ${output.rc}    0

Deploy ${lab-name} lab
    ${output} =    Process.Run Process
    ...    sudo -E ${CLAB_BIN} --runtime ${runtime} deploy -d -t ${CURDIR}/${lab-file}
    ...    shell=True

    Log    ${output.stdout}
    Log    ${output.stderr}

    Should Be Equal As Integers    ${output.rc}    0

Ensure node3 started after node4
    [Documentation]    Ensure node3 is started after node4 since node3 waits for node4 to be healthy.
    ...    All containers write the unix timestamp whenever they are started to /tmp/time file and we compare the timestamps.
    ${node3} =    Process.Run Process
    ...    sudo -E docker exec clab-${lab-name}-node3 cat /tmp/time
    ...    shell=True

    Log    ${node3.stdout}
    Log    ${node3.stderr}

    ${node4} =    Process.Run Process
    ...    sudo -E docker exec clab-${lab-name}-node4 cat /tmp/time
    ...    shell=True

    Log    ${node4.stdout}
    Log    ${node4.stderr}

    ${n3} =    Convert To Integer    ${node3.stdout}
    ${n4} =    Convert To Integer    ${node4.stdout}

    Should Be True    ${n3} > ${n4}

Deploy ${lab-name} lab with a single worker
    Run Keyword    Teardown

    ${output} =    Process.Run Process
    ...    sudo -E ${CLAB_BIN} --runtime ${runtime} deploy --max-workers 1 -t ${CURDIR}/${lab-file}
    ...    shell=True

    Log    ${output.stdout}
    Log    ${output.stderr}

    Should Be Equal As Integers    ${output.rc}    0

Ensure node3 started after node4 after a single worker run
    [Documentation]    Ensure node3 is started after node4 since node3 waits for node4 to be healthy.
    ...    All containers write the unix timestamp whenever they are started to /tmp/time file and we compare the timestamps.
    ${node3} =    Process.Run Process
    ...    sudo -E docker exec clab-${lab-name}-node3 cat /tmp/time
    ...    shell=True

    Log    ${node3.stdout}
    Log    ${node3.stderr}

    ${node4} =    Process.Run Process
    ...    sudo -E docker exec clab-${lab-name}-node4 cat /tmp/time
    ...    shell=True

    Log    ${node4.stdout}
    Log    ${node4.stderr}

    ${n3} =    Convert To Integer    ${node3.stdout}
    ${n4} =    Convert To Integer    ${node4.stdout}

    Should Be True    ${n3} > ${n4}


*** Keywords ***
Teardown
    Run    sudo -E ${CLAB_BIN} --runtime ${runtime} destroy -c -t ${CURDIR}/${lab-file}

Setup
    ${output} =    Process.Run Process
    ...    sudo -E ${CLAB_BIN} --runtime ${runtime} destroy -c -t ${CURDIR}/${lab-file}
    ...    shell=True

    Log    ${output.stdout}
    Log    ${output.stderr}
    # skipping this test suite for podman for now
    Skip If    '${runtime}' == 'podman'
