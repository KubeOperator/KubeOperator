#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function success(){
    echo -e "\033[32m kubeOperator 卸载完成. \033[0m"
}
function remove_dir() {
    echo -e "删除 KubeOperator 工作目录"
    rm -rf /opt/kubeoperator 2&> /dev/null
}
function remove_service() {
    echo -e "停止 KubeOperator 服务进程"
    systemctl stop kubeops
    systemctl disable kubeops
    rm -rf /etc/systemd/system/kubeops.service
}

function remove_images() {
    echo -e "删除 docker 镜像"
    docker ps -a | grep 'nexus-helm' | awk  '{print $1}'|xargs docker rm -f 2&> /dev/null
    docker images|grep -v IMAGE|awk '{print $3}'|xargs docker rmi 2&> /dev/null
}

function main() {
    remove_service
    remove_images
    remove_dir
    success
}

main
