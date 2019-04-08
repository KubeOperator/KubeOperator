#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function main() {
    ${SCRIPTS_DIR}/1_set_iptables.sh
    ${SCRIPTS_DIR}/2_install_docker.sh
    ${SCRIPTS_DIR}/3_download_resource.sh
    ${SCRIPTS_DIR}/4_prepare_images.sh
}

main