#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function success(){
    echo_green "kubeOperator 安装成功!\n 默认登陆信息:"  
    echo_green "username: admin"  
    echo_green "password: kubeoperator@admin123"  
    echo_green "[系统初始化中，请耐心等待5分钟再进行登录]"
}

function start_service(){
    echo_green "start service......"
    systemctl restart docker.service
    systemctl enable kubeops.service
    systemctl start kubeops.service
}

function main() {
    ${SCRIPTS_DIR}/8_check_install_env.sh
    if [[ $? != 0 ]]; then
        exit 1
    fi
    ${SCRIPTS_DIR}/1_set_iptables.sh
    ${SCRIPTS_DIR}/2_install_docker.sh
    ${SCRIPTS_DIR}/3_prepare_images.sh
    ${SCRIPTS_DIR}/4_copy_data.sh
    start_service
    success
}

main