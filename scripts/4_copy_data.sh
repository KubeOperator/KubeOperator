#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh



function main {
    mkdir -p /opt/kubeOperator >/dev/null 2>&1
    cp -rp ${SCRIPTS_DIR}/service/kubeops.service /etc/systemd/system/
    cp -rp ${PROJECT_DIR}/* /opt/kubeOperator/
    chmod -R 777 /opt/kubeOperator/docker/nexus
}
main