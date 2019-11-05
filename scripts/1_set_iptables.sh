#!/usr/bin/env bash
success=0

function set_firewall {
    systemctl stop firewalld >/dev/null 2>&1
    systemctl disable  firewalld >/dev/null 2>&1

}


function main {
    which firewall-cmd &> /dev/null
    if [[ "$?" == "0" ]];then
        set_firewall
        exit 0
    fi
}

main

