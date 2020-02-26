#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function success(){
    echo_green "kubeOperator 卸载完成."
}
function remove_dir() {
    echo_green "删除 KubeOperator 工作目录"
    rm -rf ${INSTALL_DIR}/kubeoperator 2&> /dev/null
}
function remove_service() {
    echo_green "停止 KubeOperator 服务进程"
    if [ -a /etc/systemd/system/kubeops.service ] ;then
        systemctl stop kubeops 2&> /dev/null
        systemctl disable kubeops 2&> /dev/null
        rm -rf /etc/systemd/system/kubeops.service
    fi
    if [ -a ${INSTALL_DIR}/kubeoperator/kubeopsctl.sh ]; then
        cd ${INSTALL_DIR}/kubeoperator && docker-compose down -v 
        docker ps |grep -i nexus|awk '{print $1}'|xargs docker rm -f 2&> /dev/null
    else
        read -p "强力卸载将会完全清除主机上的所有容器，是否继续： y/n : " yn
        if [ "$yn" == "Y" ] || [ "$yn" == "y" ]; then
            docker stop  $(docker ps -q -a)
            docker rm -f -v $(docker ps -q -a) 2&> /dev/null
        else
            exit 0
        fi
    fi
}

function remove_images() {
    echo_green "清理镜像中..."
    docker images -q|xargs docker rmi -f 2&> /dev/null
}

function main() {
    remove_service
    remove_images
    remove_dir
    success
}

main
