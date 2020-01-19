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
    if [ -a /etc/systemd/system/kubeops.serviced ] ;then
        systemctl stop kubeops 2&> /dev/null
        systemctl disable kubeops 2&> /dev/null
        rm -rf /etc/systemd/system/kubeops.serviced
    fi
    if [ -a /opt/kubeoperator/kubeopsctl.sh ]; then
        cd /opt/kubeoperator && docker-compose down -v 
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
    echo -e "清理镜像中..."
    docker images -q|xargs docker rmi -f 2&> /dev/null
}

function main() {
    remove_service
    remove_images
    remove_dir
    success
}

main
