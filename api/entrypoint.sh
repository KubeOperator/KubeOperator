#!/bin/bash
function cleanup()
{
    local pids=`jobs -p`
    if [[ "${pids}" != ""  ]]; then
        kill ${pids} >/dev/null 2>/dev/null
    fi
}

trap cleanup EXIT

if [[ "$1" == "bash" ]];then
    bash
else
    echo -e "nameserver 8.8.8.8 \nnameserver 114.114.114.114" >> /etc/resolv.conf
    service dnsmasq start
    python kubeops.py start
fi

