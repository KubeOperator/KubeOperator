#!/usr/bin/env bash

function main() {
    ./1_set_iptables.sh
    ./2_install_docker.sh
    ./3_download_resource.sh
    ./4_prepare_images.sh
}

main