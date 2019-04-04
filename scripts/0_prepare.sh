#!/usr/bin/env bash
BASE_DIR=$(dirname $(dirname "$0"))


function download_docker() {
    DOCKER_VERSION=18.06.2-ce
    DOCKER_MD5=8c4a1d65ddcecf91ae357b434dffe039
    DOCKER_COMPOSE_VERSION=1.23.2
    DOCKER_COMPOSE_MD5=7f508b543123e8c81ca138d5b36001a2
    DOCKER_BIN_DIR="${BASE_DIR}/docker/bin"
    DOCKER_COMPOSE_BIN="${DOCKER_BIN_DIR}/docker-compose"

    echo ">>> 开始下载 docker程序"
    if [[ ! -f /tmp/docker.tar.gz || ! `md5sum /tmp/docker.tar.gz | awk '$1~/'"$1"'/ {print $1}'` == ${DOCKER_MD5} ]]; then
       wget http://jms-pkg.oss-cn-beijing.aliyuncs.com/docker/docker-${DOCKER_VERSION}.tgz -O /tmp/docker.tar.gz
    fi
    if [[ ! -f "${DOCKER_BIN_DIR}/dockerd" ]];then
        tar xzf /tmp/docker.tar.gz -C ${DOCKER_BIN_DIR}
    fi
    if [[ ! -f "${DOCKER_COMPOSE_BIN}" || ! `md5sum ${DOCKER_COMPOSE_BIN} | awk '$1~/'"$1"'/ {print $1}'` == ${DOCKER_COMPOSE_MD5}  ]]; then
        wget http://jms-pkg.oss-cn-beijing.aliyuncs.com/docker/docker-compose-${DOCKER_COMPOSE_VERSION} -O ${DOCKER_COMPOSE_BIN}
    fi
}


function build_and_save_images() {
    IMAGE_DIR="${BASE_DIR}/docker/images"

    echo ">>> 开始build镜像"
    images=(
        "redis:alpine"
        "mysql:5"
        "nginx:alpine"
        "sonatype/nexus3"
        "fit2openshift/api:latest"
        "fit2openshift/ui:latest"
        "fit2openshift/dns:latest"
    )
    cd ${BASE_DIR}
    docker-compose pull
    docker-compose build

    echo ">>> 开始保存镜像"
    for image in ${images[@]};do
        filename=$(basename ${image}).tar
        docker save -o ${IMAGE_DIR}/${filename} ${image}
    done
}

function download_resources() {
    echo ">>> 开始下载resource"
    NEXUS_TAR_PATH="${BASE_DIR}/docker/nexus/nexus-data.tar.gz"
    NEXUS_DATA_PATH="${BASE_DIR}/docker/nexus/data/"

    if [[ ! -f "${NEXUS_TAR_PATH}" ]];then
        wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11//tmp/nexus-data.tar.gz" -O ${NEXUS_TAR_PATH}
    elif [[ $(du -sh nexus-data.tar.gz | grep 'G' | awk -F. '{ print $1 }') -gt 5 ]];then
        wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11//tmp/nexus-data.tar.gz" -O ${NEXUS_TAR_PATH}
    fi
}

function main() {
    download_docker
    build_and_save_images
    download_resources
}

main