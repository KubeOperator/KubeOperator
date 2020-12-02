#!/usr/bin/env bash
# Kindly borrowed from Kind
# https://github.com/kubernetes-sigs/kind/blob/c5298b2/hack/update/generated.sh
#
# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# 'go generate's kind, using tools from vendor (go-bindata)
set -o errexit -o nounset -o pipefail

# cd to the repo root
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
cd "${REPO_ROOT}"

# build the generators using the tools module
cd "hack/tools"
"${REPO_ROOT}/hack/go_container.sh" go build -o /out/deepcopy-gen k8s.io/code-generator/cmd/deepcopy-gen
"${REPO_ROOT}/hack/go_container.sh" go build -o /out/client-gen k8s.io/code-generator/cmd/client-gen
"${REPO_ROOT}/hack/go_container.sh" go build -o /out/lister-gen k8s.io/code-generator/cmd/lister-gen
"${REPO_ROOT}/hack/go_container.sh" go build -o /out/informer-gen k8s.io/code-generator/cmd/informer-gen
# go back to the root
cd "${REPO_ROOT}"

# turn off module mode before running the generators
# https://github.com/kubernetes/code-generator/issues/69
# we also need to populate vendor
hack/go_container.sh go mod tidy
hack/go_container.sh go mod vendor
export GO111MODULE="off"

# fake being in a gopath
FAKE_GOPATH="$(mktemp -d)"
trap 'rm -rf ${FAKE_GOPATH}' EXIT

FAKE_REPOPATH="${FAKE_GOPATH}/src/github.com/fluxcd/helm-operator"
mkdir -p "$(dirname "${FAKE_REPOPATH}")" && ln -s "${REPO_ROOT}" "${FAKE_REPOPATH}"

export GOPATH="${FAKE_GOPATH}"
cd "${FAKE_REPOPATH}"

# run the generators
echo "Generating deepcopy funcs"
bin/deepcopy-gen \
  --input-dirs "./pkg/apis/helm.fluxcd.io/v1" \
  -O zz_generated.deepcopy \
  --bounding-dirs "./pkg/apis" \
  --go-header-file hack/tools/boilerplate.go.txt

echo "Generating clientset"
bin/client-gen \
  --clientset-name "versioned" \
  --input-base "" \
  --input-dirs "./pkg/apis/helm.fluxcd.io/v1" \
  --output-package "./pkg/client/clientset" \
  --go-header-file hack/tools/boilerplate.go.txt

echo "Generating listers"
bin/lister-gen \
  --input-dirs "./pkg/apis/helm.fluxcd.io/v1" \
  --output-package "./pkg/client/listers" \
  --go-header-file hack/tools/boilerplate.go.txt

echo "Generating informers"
bin/informer-gen \
  --input-dirs "./pkg/apis/helm.fluxcd.io/v1" \
  --versioned-clientset-package "./pkg/client/clientset/versioned" \
  --listers-package "./pkg/client/listers" \
  --output-package "./pkg/client/informers" \
  --go-header-file hack/tools/boilerplate.go.txt

# set module mode back, return to repo root
export GO111MODULE="on"
cd "${REPO_ROOT}"
