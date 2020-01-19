#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh
OS=$(uname)

OFFLINE_DOCKER_DIR="${PROJECT_DIR}/docker/bin"

function prepare_install {
    yum install -y epel-release yum-utils device-mapper-persistent-data lvm2
}


function install_docker_online {
    prepare_install
    yum -y remove docker docker-common docker-engine
    wget -O /etc/yum.repos.d/docker-ce.repo https://download.docker.com/linux/centos/docker-ce.repo
    sed -i 's+download.docker.com+mirrors.tuna.tsinghua.edu.cn/docker-ce+' /etc/yum.repos.d/docker-ce.repo
    yum install -y docker-ce 
}

function install_docker-compose_online {
    yum install -y docker-compose
}

function install_docker_offline {
    cp ${OFFLINE_DOCKER_DIR}/docker/docker* /usr/bin/
    cp ${OFFLINE_DOCKER_DIR}/docker.service /etc/systemd/system/
    chmod +x /usr/bin/docker* && chmod 754 /etc/systemd/system/docker.service
}

function install_docker-compose_offline {
    cp ${OFFLINE_DOCKER_DIR}/docker-compose /usr/bin/
    chmod +x /usr/bin/docker-compose
}

function install_docker {
    echo ">> Install docker"
    if [[ "${OS}" == "Darwin" ]];then
        echo "Platform is MacOS, install manually"
        return
    fi
    if [[ -f "${OFFLINE_DOCKER_DIR}/docker/dockerd" ]];then
        install_docker_offline
    else
        install_docker_online
    fi
}

function install_docker-compose {
    echo ">> Install docker-compose"
    if [[ "${OS}" == "Darwin" ]];then
        echo "Platform is MacOS, install manually"
        return
    fi
    if [[ -f "${OFFLINE_DOCKER_DIR}/docker-compose" ]];then
        install_docker-compose_offline
    else
        install_docker-compose_online
    fi
}

function config_docker {
    set_docker_config registry-mirrors '["https://mirror.ccs.tencentyun.com"]'
}

function start_docker {
    systemctl start docker
    systemctl enable docker
}

function main {
    which docker >/dev/null 2>&1
    if [ $? -ne 0 ];then
       install_docker
       config_docker
       start_docker
    fi
    which docker-compose >/dev/null 2>&1
    if [ $? -ne 0 ];then
       install_docker-compose
    fi
}

main
