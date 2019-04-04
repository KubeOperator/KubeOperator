#!/usr/bin/env bash
BASE_DIR=$(dirname $(dirname "$0"))

NEXUS_TAR_PATH="${BASE_DIR}/docker/nexus/nexus-data.tar.gz"
NEXUS_DATA_PATH="${BASE_DIR}/docker/nexus/data/"

if [[ ! -f "${NEXUS_TAR_PATH}" ]];then
    wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/tmp/nexus-data.tar.gz" -O ${NEXUS_DATA_PATH}
fi

tar xvf ${NEXUS_TAR_PATH} -C ${NEXUS_DATA_PATH}
