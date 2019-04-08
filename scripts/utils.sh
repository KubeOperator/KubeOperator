#!/usr/bin/env bash
BASE_DIR=$(dirname $(dirname "$0"))
IMAGE_DIR="${BASE_DIR}/docker/images"


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