#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh

function download_nexus_resource {
    NEXUS_TAR_PATH="${PROJECT_DIR}/docker/nexus/nexus-data.tar.gz"
    NEXUS_DATA_PATH="${PROJECT_DIR}/docker/nexus/data/"

    if [[ ! -f "${NEXUS_TAR_PATH}" && "${DOWNLOAD_NEXUS_DATA}" != "0" ]];then
        wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd/v3/nexus-data.tar.gz" -O ${NEXUS_TAR_PATH}
    fi

    if [[ ! -f "${NEXUS_DATA_PATH}/.done" ]];then
        tar xvf ${NEXUS_TAR_PATH} -C ${NEXUS_DATA_PATH} && echo > ${NEXUS_DATA_PATH}/.done
        chown -R 200 ${NEXUS_DATA_PATH}
    fi
}

function main {
    download_nexus_resource
}

main
