#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function download_docker() {
    DOCKER_VERSION=18.06.2-ce
    DOCKER_MD5=8c4a1d65ddcecf91ae357b434dffe039
    DOCKER_COMPOSE_VERSION=1.23.2
    DOCKER_COMPOSE_MD5=7f508b543123e8c81ca138d5b36001a2
    DOCKER_BIN_DIR="${PROJECT_DIR}/docker/bin"
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
    echo ">>> 开始build镜像"
    images=$(get_images)
    cd ${PROJECT_DIR}
    docker-compose build

    echo ">>> 开始保存镜像"
    for image in ${images};do
        if [[ ! ${image} =~ 'kubeops' ]];then
            docker pull ${image}
        fi
        filename=$(basename ${image}).tar
        docker save -o ${IMAGE_DIR}/${filename} | gzip -c > ${image}
    done
}


function main() {
    download_docker && build_and_save_images  || exit 10
}

main