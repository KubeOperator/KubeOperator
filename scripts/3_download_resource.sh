#!/usr/bin/env bash
source ./utils.sh

function download_nexus_resource {
    NEXUS_TAR_PATH="${BASE_DIR}/docker/nexus/nexus-data.tar.gz"
    NEXUS_DATA_PATH="${BASE_DIR}/docker/nexus/data/"

    if [[ ! -f "${NEXUS_TAR_PATH}" ]];then
        wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/tmp/nexus-data.tar.gz" -O ${NEXUS_DATA_PATH}
    fi

    if [[ ! -f "${NEXUS_DATA_PATH}/.done" ]];then
        tar xvf ${NEXUS_TAR_PATH} -C ${NEXUS_DATA_PATH} && echo > ${NEXUS_DATA_PATH}/.done
    fi
}

function main {
    download_nexus_resource
}

main
