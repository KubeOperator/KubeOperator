#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")

nexus_file_path="${BASE_DIR}/../docker/nexus/nexus-data.tar.gz"
nexus_data_dir="${BASE_DIR}/../docker/nexus/data/"

if [[ ! -f "${nexus_file_path}" ]];then
    wget "http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/tmp/nexus-data.tar.gz" -O ${nexus_file_path}
fi
