#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function remove_dir {
    rm -rf /opt/kubeOperator 2&> /dev/null
}

function remove_images {
    echo -e "停止KubeOperator服务进程"
    cd ${PROJECT_DIR} && ${PROJECT_DIR}/kubeopsctl.sh down
    docker images|grep -v IMAGE|awk '{print $3}'|xargs docker rmi
}

function main {
    remove_images
    remove_dir
}

main
