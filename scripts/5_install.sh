#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function success(){
    echo -e "\033[31m kubeOperator 安装成功!\n 默认登陆信息: \033[0m"  
    echo -e "\033[32m username: admin \033[0m"  
    echo -e "\033[32m password: kubeoperator@admin123 \033[0m"  
}

function start_service(){
    systemctl restart docker.service
    systemctl enable kubeops.service
    systemctl start kubeops.service
}

function main() {
    ${SCRIPTS_DIR}/1_set_iptables.sh
    ${SCRIPTS_DIR}/2_install_docker.sh
    ${SCRIPTS_DIR}/3_prepare_images.sh
    start_service
    success
}

main