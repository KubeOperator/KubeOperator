#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh
images=$(get_images)

function load_images() {
    echo ">>> 开始加载镜像"
    for image in ${images};do
        filename=$(basename ${image}).tar
        docker load < ${IMAGE_DIR}/${filename}
    done
}

function build_image() {
    echo ">>> 开始build镜像"
    cd ${PROJECT_DIR}
    for image in ${images};do
        if [[ ! ${image} =~ 'kube-operator' ]];then
            docker pull ${image}
        fi
    done
    docker-compose build
    cd -
}

function find_offline_images() {
    ok=1
    for image in ${images};do
        filename=$(basename ${image}).tar
        if [[ ! -f ${IMAGE_DIR}/${filename} ]];then
            ok=0
        fi
    done
    echo ${ok}
    return ${ok}
}

function main() {
    chmod -R 777 ${PROJECT_DIR}/docker/nexus/data
    find_offline_images
    ok=$?
    if [[ ${ok} == "1" ]];then
        load_images
    else
        build_image
    fi
}

main