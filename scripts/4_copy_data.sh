#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh



function main {
    mkdir -p ${INSTALL_DIR}/kubeoperator >/dev/null 2>&1
    cp -rp ${SCRIPTS_DIR}/service/kubeops.service /etc/systemd/system/
    sed -i -r "s|INSTALL_DIR|${INSTALL_DIR}|g" /etc/systemd/system/kubeops.service
    cp -rp ${PROJECT_DIR}/* ${INSTALL_DIR}/kubeoperator/
    chmod -R 777 ${INSTALL_DIR}/kubeoperator/data
}
main