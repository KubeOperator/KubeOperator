#!/bin/bash

# This is intended to run either locally (in which case a push is not
# necessary) or in a CI job (where the results should be pushed to a
# registry for use in later CI test jobs). Images are tagged locally with
# both the base name (e.g. "test-target-base") and with the prefix given;
# then only the prefixed name is pushed if --push is specified, assuming
# any necessary credentials are available for the push. The same prefix
# can then be used for the testing script. By default a local (non-registry)
# prefix is used and no push can occur. To push to e.g. dockerhub:
#
# ./build-images.sh --push --prefix=docker.io/openshift/ansible-integration-

set -o errexit
set -o nounset
set -o pipefail

STARTTIME=$(date +%s)
source_root=$(dirname "${0}")

prefix="${PREFIX:-openshift-ansible-integration-}"
push=false
verbose=false
build_options="${DOCKER_BUILD_OPTIONS:-}"
help=false

for args in "$@"
do
  case $args in
      --prefix=*)
        prefix="${args#*=}"
        ;;
      --push)
        push=true
        ;;
      --no-cache)
        build_options="${build_options} --no-cache"
        ;;
      --verbose)
        verbose=true
        ;;
      --help)
        help=true
        ;;
  esac
done

if [ "$help" = true ]; then
  echo "Builds the docker images for openshift-ansible integration tests"
  echo "and pushes them to a central registry."
  echo
  echo "Options: "
  echo "  --prefix=PREFIX"
  echo "  The prefix to use for the image names."
  echo "  default: openshift-ansible-integration-"
  echo
  echo "  --push"
  echo "  If set will push the tagged image"
  echo 
  echo "  --no-cache"
  echo "  If set will perform the build without a cache."
  echo
  echo "  --verbose"
  echo "  Enables printing of the commands as they run."
  echo
  echo "  --help"
  echo "  Prints this help message"
  echo
  exit 0
fi

if [ "$verbose" = true ]; then
  set -x
fi


declare -a build_order                       ; declare -A images
build_order+=( test-target-base )            ; images[test-target-base]=openshift_health_checker/builds/test-target-base
build_order+=( preflight-aos-package-checks ); images[preflight-aos-package-checks]=openshift_health_checker/builds/aos-package-checks
for image in "${build_order[@]}"; do
  BUILD_STARTTIME=$(date +%s)
  docker_tag=${prefix}${image}
  echo
  echo "--- Building component '$image' with docker tag '$docker_tag' ---"
  docker build ${build_options} -t $image -t $docker_tag "$source_root/${images[$image]}"
  echo
  BUILD_ENDTIME=$(date +%s); echo "--- build $docker_tag took $(($BUILD_ENDTIME - $BUILD_STARTTIME)) seconds ---"
  if [ "$push" = true ]; then
    docker push $docker_tag
    PUSH_ENDTIME=$(date +%s); echo "--- push $docker_tag took $(($PUSH_ENDTIME - $BUILD_ENDTIME)) seconds ---"
  fi
done

echo
echo
echo "++ Active images"
docker images | grep ${prefix} | sort
echo


ret=$?; ENDTIME=$(date +%s); echo "$0 took $(($ENDTIME - $STARTTIME)) seconds"; exit "$ret"
