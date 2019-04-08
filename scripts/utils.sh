#!/usr/bin/env bash
PROJECT_DIR=$(dirname $(dirname "$0"))
IMAGE_DIR="${PROJECT_DIR}/docker/images"
SCRIPTS_DIR="${PROJECT_DIR}/scripts"

function get_images(){
    images=(
       "redis:alpine"
       "mysql:5"
       "nginx:alpine"
       "sonatype/nexus3"
       "fit2openshift/api:latest"
       "fit2openshift/ui:latest"
       "fit2openshift/dns:latest"
    )
    for image in ${images[@]};do
        echo ${image}
    done
}