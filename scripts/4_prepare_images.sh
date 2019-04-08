#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function load_images() {
    images=$(get_images)
    echo ">>> 加载镜像"
    for image in ${images};do
        filename=$(basename ${image}).tar
        docker load < ${IMAGE_DIR}/${filename}
    done
}

function build_image() {
    cd ${PROJECT_DIR}
    docker-compose pull &> /dev/null
    docker-compose build
    cd -
}

function find_offline_images() {
    ok=1
    for image in ${images};do
        filename=$(basename ${image}).tar
        if [[ ! -f ${filename} ]];then
            ok=0
        fi
    done
    echo ${ok}
}

function main() {
    ok=$(find_offline_images)
    if [[ ${ok} == "1" ]];then
        load_images
    else
        build_image
    fi
}

main