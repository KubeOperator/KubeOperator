#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

STARTTIME=$(date +%s)
source_root=$(dirname "${0}")/..

prefix="openshift/origin-ansible"
version="latest"
verbose=false
options="-f images/installer/Dockerfile"
help=false

for args in "$@"
do
  case $args in
      --prefix=*)
        prefix="${args#*=}"
        ;;
      --version=*)
        version="${args#*=}"
        ;;
      --no-cache)
        options="${options} --no-cache"
        ;;
      --verbose)
        verbose=true
        ;;
     --help)
        help=true
        ;;
  esac
done

# allow ENV to take precedent over switches
prefix="${PREFIX:-$prefix}"
version="${OS_TAG:-$version}"

if [ "$help" = true ]; then
  echo "Builds the docker images for openshift-ansible"
  echo
  echo "Options: "
  echo "  --prefix=PREFIX"
  echo "  The prefix to use for the image names."
  echo "  default: docker.io/openshift/origin-ansible"
  echo
  echo "  --version=VERSION"
  echo "  The version used to tag the image (can be a comma-separated list)"
  echo "  default: latest"
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

BUILD_STARTTIME=$(date +%s)
comp_path=$source_root/

# turn comma-separated versions into -t args for docker build
IFS=',' read -r -a version_arr <<< "$version"
docker_tags=()
for tag in "${version_arr[@]}"; do
  docker_tags+=("-t" "${prefix}:${tag}")
done

echo
echo
echo "--- Building component '$comp_path' with docker tag(s) '$version' ---"
docker build ${options} "${docker_tags[@]}" $comp_path
BUILD_ENDTIME=$(date +%s); echo "--- ${version} took $(($BUILD_ENDTIME - $BUILD_STARTTIME)) seconds ---"
echo
echo

echo
echo
echo "++ Active images"
docker images | grep ${prefix} | sort
echo


ret=$?; ENDTIME=$(date +%s); echo "$0 took $(($ENDTIME - $STARTTIME)) seconds"; exit "$ret"
