#!/usr/bin/env bash
PROJECT_DIR=$(dirname $(cd $(dirname "$0");pwd))
IMAGE_DIR="${PROJECT_DIR}/docker/images"
SCRIPTS_DIR="${PROJECT_DIR}/scripts"

function get_images(){
    images=(
       "docker.io/redis:alpine"
       "docker.io/mysql:5"
       "docker.io/nginx:alpine"
       "kube-operator/core:2.5.0"
       "kube-operator/ui:2.5.0"
       "registry.fit2cloud.com/public/nexus-helm:3.15.2-01"
       "kubeoperator/webkubectl:v2.0"
       "elasticsearch:7.4.1"
    )
    for image in ${images[@]};do
        echo ${image}
    done
}


function set_docker_config() {
   key=$1
   value=$2
   DOCKER_CONFIG="/etc/docker/daemon.json"

   if [[ ! -f "${DOCKER_CONFIG}" ]];then
       config_dir=$(dirname ${DOCKER_CONFIG})
       if [[ ! -d ${config_dir} ]];then
           mkdir -p ${config_dir}
       fi
        echo -e "{\n}" >> ${DOCKER_CONFIG}
   fi
   $(python -c "import json
key = '${key}'
value = '${value}'
try:
    value = json.loads(value)
except:
    pass
filepath = \"${DOCKER_CONFIG}\"
f = open(filepath);
config = json.load(f);
config[key] = value
f.close();
f = open(filepath, 'w');
json.dump(config, f, indent=True, sort_keys=True);
f.close()
")
}
