#!/bin/bash

# This script runs the golang integration tests in the directories underneath.
# It should be run from the same directory it is in, or in a directory above.
# Specify the same image prefix used (if any) with build-images.sh
#
# Example:
# ./run-tests.sh --prefix=docker.io/openshift/ansible-integration- --parallel=16

set -o errexit
set -o nounset
set -o pipefail

source_root=$(dirname "${0}")

prefix="${PREFIX:-openshift-ansible-integration-}"
gotest_options="${GOTEST_OPTIONS:--v}"
push=false
verbose=false
help=false

for args in "$@"
do
  case $args in
      --prefix=*)
        prefix="${args#*=}"
        ;;
      --parallel=*)
        gotest_options="${gotest_options} -parallel ${args#*=}"
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
  echo "Runs the openshift-ansible integration tests."
  echo
  echo "Options: "
  echo "  --prefix=PREFIX"
  echo "  The prefix to use for the image names."
  echo "  default: openshift-ansible-integration-"
  echo
  echo "  --parallel=NUMBER"
  echo "  Number of tests to run in parallel."
  echo "  default: GOMAXPROCS (typically, number of processors)"
  echo
  echo "  --verbose"
  echo "  Enables printing of the commands as they run."
  echo
  echo "  --help"
  echo "  Prints this help message"
  echo
  exit 0
fi



if ! [ -d $source_root/../../.tox/integration ]; then
  # have tox create a consistent virtualenv
  pushd $source_root/../..; tox -e integration; popd
fi
# use the virtualenv from tox
set +o nounset; source $source_root/../../.tox/integration/bin/activate; set -o nounset

if [ "$verbose" = true ]; then
  set -x
fi

# Run the tests. NOTE: "go test" requires a relative path for this purpose.
# The PWD trick below will only work if cwd is in/above where this script lives.
retval=0
IMAGE_PREFIX="${prefix}" env -u GOPATH \
  go test ./${source_root#$PWD}/... ${gotest_options}


