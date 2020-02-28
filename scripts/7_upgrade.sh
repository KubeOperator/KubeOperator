#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh
KUBEOPS_DIR="${INSTALL_DIR}/kubeoperator/"



function success(){
    echo_green "kubeOperator 升级完成."
}

function stop_service() {
    echo -ne "停止 KubeOperator 服务进程 ... "
    systemctl stop kubeops
    if [ "$?" -eq "0" ];then
        echo_green "[OK]"
    else
        echo_red "[ERROR]"
        exit 1
    fi
    echo_green "停止 nexus 容器 ..."
    container_id=(`docker ps -a | grep nexus-helm | awk '{print $1}'`)
    for cid in ${container_id[@]}; do
        docker rm -f $cid
    done
}

function start_service() {
    echo -ne "启动 KubeOperator 服务进程 ... "
    systemctl start kubeops
    if [ "$?" -eq "0" ];then
        echo_green "[OK]"
    else
        echo_red "[ERROR]"
        exit 1
    fi
}

function check_docker_service() {
    echo -ne "检测 docker 服务状态 ... "
    result=`docker ps 2>&1`
    if [ "$?" -eq "0" ]; then
        echo_green "[OK]"
    else
        if [[ $(echo $result | grep 'not found' | wc -l) == 1 ]]; then
            echo_red "[ERROR] 没有找到 docker 服务"
            exit 1
        elif [[ $(echo $result | grep 'running?' | wc -l) == 1 ]]; then
            echo_red "[ERROR] docker 服务未运行，请启动 docker 服务"
            exit 1
        else
            echo_red "[ERROR] 请检查 docker 服务"
        fi
    fi
}

function upgrade_service() {
    echo_green "开始升级 KubeOperator"
    echo -ne "清理旧镜像文件 ... "
    rm -rf ${KUBEOPS_DIR}/docker/images
    echo_green "[OK]"

    compose_file="${KUBEOPS_DIR}/docker-compose.yml"
    compose_bak="${KUBEOPS_DIR}/docker-compose.yml.bak"
    \mv -f ${compose_file} ${compose_bak}

    echo -ne "更新升级文件 ... "
    \cp -rf ${PROJECT_DIR}/* ${KUBEOPS_DIR}/
    chmod -R 777 ${KUBEOPS_DIR}/data
    echo_green "[OK]"

    kubeops_service_old="/etc/systemd/system/kubeops.service"
    kubeops_service_new="${PROJECT_DIR}/scripts/service/kubeops.service"
    diffLine=`diff ${kubeops_service_old} ${kubeops_service_new} | wc -l`
    if [ ! "$diffLine" -eq "0" ];then
        echo -ne "kubeops 服务有更新，升级 kubeops 服务 ... "
        \cp -f ${PROJECT_DIR}/scripts/service/kubeops.service /etc/systemd/system/kubeops.service
        systemctl daemon-reload
        echo_green "[OK]"
    fi

    echo -ne "移除旧版本镜像 ... "
    for image in $(grep "\simage: " ${compose_bak}  | awk -F' ' '{print $2}'); do
        docker rmi -f $image > /dev/null 2>&1
    done
    echo_green "[OK]"

    echo -ne "加载最新镜像 ... "
    docker_images_folder="${KUBEOPS_DIR}/docker/images"
    for docker_image in ${docker_images_folder}/*; do
        temp_file=`basename $docker_image`
        docker load -q -i ${docker_images_folder}/$temp_file > /dev/null 2>&1
    done
    echo_green "[OK]"
}

function main() {
    check_docker_service && stop_service && upgrade_service && start_service && success
}

main
