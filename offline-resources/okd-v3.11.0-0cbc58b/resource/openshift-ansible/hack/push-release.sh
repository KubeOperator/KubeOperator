#!/bin/bash

# This script pushes a built image to a registry.
#
# Set OS_PUSH_BASE_REGISTRY to prefix the destination images e.g.
# OS_PUSH_BASE_REGISTRY="docker.io/"
#
# Set OS_PUSH_TAG with a comma-separated list for pushing same image
# to multiple tags e.g.
# OS_PUSH_TAG="latest,v3.6"

set -o errexit
set -o nounset
set -o pipefail

starttime=$(date +%s)

# image name without repo or tag.
image="${PREFIX:-docker.io/openshift/origin-ansible}"

# existing local tag on the image we want to push
source_tag="${OS_TAG:-latest}"

# Enable retagging a build with one or more tags for push
IFS=',' read -r -a push_tags <<< "${OS_PUSH_TAG:-latest}"
registry="${OS_PUSH_BASE_REGISTRY:-}"

# force push if available
PUSH_OPTS=""
if docker push --help | grep -q force; then
  PUSH_OPTS="--force"
fi

set -x
for tag in "${push_tags[@]}"; do
  docker tag "${image}:${source_tag}" "${registry}${image}:${tag}"
  docker push ${PUSH_OPTS} "${registry}${image}:${tag}"
done
set +x

endtime=$(date +%s); echo "$0 took $(($endtime - $starttime)) seconds"; exit 0
