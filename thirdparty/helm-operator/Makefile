.DEFAULT: all
.PHONY: all clean realclean test integration-test check-generated lint-e2e

SUDO := $(shell docker info > /dev/null 2> /dev/null || echo "sudo")

TEST_FLAGS?=

BATS_COMMIT := 3a1c2f28be260f8687ff83183cef4963faabedd6
SHELLCHECK_VERSION := 0.7.0
SHFMT_VERSION := 2.6.4

include docker/kubectl.version
include docker/helm2.version
include docker/helm3.version

# NB default target architecture is amd64. If you would like to try the
# other one -- pass an ARCH variable, e.g.,
#  `make ARCH=arm64`
ifeq ($(ARCH),)
	ARCH=amd64
endif
CURRENT_OS_ARCH=$(shell echo `go env GOOS`-`go env GOARCH`)
GOBIN?=$(shell echo `go env GOPATH`/bin)

MAIN_GO_MODULE:=$(shell go list -m -f '{{ .Path }}')
LOCAL_GO_MODULES:=$(shell go list -m -f '{{ .Path }}' all | grep $(MAIN_GO_MODULE))
godeps=$(shell go list -deps -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}' $(1) | sed "s%${PWD}/%%g")

HELM_OPERATOR_DEPS:=$(call godeps,./cmd/helm-operator/...)

IMAGE_TAG:=$(shell ./docker/image-tag)
VCS_REF:=$(shell git rev-parse HEAD)
BUILD_DATE:=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

DOCS_PORT:=8000

all: $(GOBIN)/bin/helm-operator build/.helm-operator.done

clean:
	go clean ./cmd/helm-operator
	rm -rf ./build
	rm -f test/bin/kubectl test/bin/helm2 test/bin/helm3 test/bin/kind

realclean: clean
	rm -rf ./cache

test: test/bin/helm2 test/bin/helm3
	PATH="${PWD}/bin:${PWD}/test/bin:${PATH}" go test ${TEST_FLAGS} $(shell go list $(patsubst %, %/..., $(LOCAL_GO_MODULES)) | sort -u)

e2e: test/bin/helm2 test/bin/helm3 test/bin/kubectl test/e2e/bats build/.helm-operator.done
	PATH="${PWD}/test/bin:${PATH}" CURRENT_OS_ARCH=$(CURRENT_OS_ARCH) test/e2e/run.bash

E2E_BATS_FILES := test/e2e/*.bats
E2E_BASH_FILES := test/e2e/run.bash test/e2e/lib/*
SHFMT_DIFF_CMD := test/bin/shfmt -i 2 -sr -d
SHFMT_WRITE_CMD := test/bin/shfmt -i 2 -sr -w
lint-e2e: test/bin/shfmt test/bin/shellcheck
	@# shfmt is not compatible with .bats files, so we preprocess them to turn '@test's into functions
	for I in $(E2E_BATS_FILES); do \
	  ( cat "$$I" | sed 's%@test.*%test() {%' | $(SHFMT_DIFF_CMD) ) || ( echo "Please correct the diff for file $$I"; exit 1 ); \
	done
	$(SHFMT_DIFF_CMD) $(E2E_BASH_FILES) || ( echo "Please run '$(SHFMT_WRITE_CMD) $(E2E_BASH_FILES)'"; exit 1 )
	test/bin/shellcheck $(E2E_BASH_FILES) $(E2E_BATS_FILES)

build/.%.done: docker/Dockerfile.%
	mkdir -p ./build/docker/$*
	cp $^ ./build/docker/$*/
	$(SUDO) docker build -t docker.io/fluxcd/$* -t docker.io/fluxcd/$*:$(IMAGE_TAG) \
		--build-arg VCS_REF="$(VCS_REF)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		-f build/docker/$*/Dockerfile.$* ./build/docker/$*
	touch $@

build/.helm-operator.done: build/helm-operator build/kubectl build/helm2 build/helm3 docker/ssh_config docker/known_hosts.sh docker/helm-repositories.yaml

build/helm-operator: $(HELM_OPERATOR_DEPS)
build/helm-operator: cmd/helm-operator/*.go
	CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -o $@ $(LDFLAGS) -ldflags "-X main.version=$(shell ./docker/image-tag)" ./cmd/helm-operator

build/kubectl: cache/linux-$(ARCH)/kubectl-$(KUBECTL_VERSION)
test/bin/kubectl: cache/$(CURRENT_OS_ARCH)/kubectl-$(KUBECTL_VERSION)
build/helm2: cache/linux-$(ARCH)/helm-$(HELM2_VERSION)
build/helm3: cache/linux-$(ARCH)/helm-$(HELM3_VERSION)
test/bin/helm2: cache/$(CURRENT_OS_ARCH)/helm-$(HELM2_VERSION)
test/bin/helm3: cache/$(CURRENT_OS_ARCH)/helm-$(HELM3_VERSION)
test/bin/shellcheck: cache/$(CURRENT_OS_ARCH)/shellcheck-$(SHELLCHECK_VERSION)
test/bin/shfmt: cache/$(CURRENT_OS_ARCH)/shfmt-$(SHFMT_VERSION)
build/kubectl test/bin/kubectl build/helm2 build/helm3 test/bin/helm2 test/bin/helm3 test/bin/shellcheck test/bin/shfmt:
	mkdir -p build
	cp $< $@
	if [ `basename $@` = "build" -a $(CURRENT_OS_ARCH) = "linux-$(ARCH)" ]; then strip $@; fi
	chmod a+x $@

cache/%/kubectl-$(KUBECTL_VERSION): docker/kubectl.version
	mkdir -p cache/$*
	curl --fail -L -o cache/$*/kubectl-$(KUBECTL_VERSION).tar.gz "https://dl.k8s.io/$(KUBECTL_VERSION)/kubernetes-client-$*.tar.gz"
	[ $* != "linux-$(ARCH)" ] || echo "$(KUBECTL_CHECKSUM_$(ARCH))  cache/$*/kubectl-$(KUBECTL_VERSION).tar.gz" | shasum -a 512 -c
	tar -m --strip-components 3 -C ./cache/$* -xzf cache/$*/kubectl-$(KUBECTL_VERSION).tar.gz kubernetes/client/bin/kubectl
	mv ./cache/$*/kubectl $@

cache/%/helm-$(HELM2_VERSION): docker/helm2.version
	mkdir -p cache/$*
	curl --fail -L -o cache/$*/helm-$(HELM2_VERSION).tar.gz "https://storage.googleapis.com/kubernetes-helm/helm-v$(HELM2_VERSION)-$*.tar.gz"
	[ $* != "linux-$(ARCH)" ] || echo "$(HELM2_CHECKSUM_$(ARCH))  cache/$*/helm-$(HELM2_VERSION).tar.gz" | shasum -a 256 -c
	tar -m -C ./cache -xzf cache/$*/helm-$(HELM2_VERSION).tar.gz $*/helm
	mv cache/$*/helm $@

cache/%/helm-$(HELM3_VERSION): docker/helm3.version
	mkdir -p cache/$*
	curl --fail -L -o cache/$*/helm-$(HELM3_VERSION).tar.gz "https://get.helm.sh/helm-v$(HELM3_VERSION)-$*.tar.gz"
	[ $* != "linux-$(ARCH)" ] || echo "$(HELM3_CHECKSUM_$(ARCH))  cache/$*/helm-$(HELM3_VERSION).tar.gz" | shasum -a 256 -c
	tar -m -C ./cache -xzf cache/$*/helm-$(HELM3_VERSION).tar.gz $*/helm
	mv cache/$*/helm $@

$(GOBIN)/bin/helm-operator: $(HELM_OPERATOR_DEPS)
	go install ./cmd/helm-operator

pkg/install/generated_templates.gogen.go: pkg/install/templates/*
	cd pkg/install && go run generate.go embedded-templates

cache/%/shellcheck-$(SHELLCHECK_VERSION):
	mkdir -p cache/$*
	curl --fail -L -o cache/$*/shellcheck-$(SHELLCHECK_VERSION).tar.xz "https://storage.googleapis.com/shellcheck/shellcheck-v$(SHELLCHECK_VERSION).$(CURRENT_OS).x86_64.tar.xz"
	tar -C cache/$* --strip-components 1 -xvJf cache/$*/shellcheck-$(SHELLCHECK_VERSION).tar.xz shellcheck-v$(SHELLCHECK_VERSION)/shellcheck
	mv cache/$*/shellcheck $@

cache/%/shfmt-$(SHFMT_VERSION):
	mkdir -p cache/$*
	curl --fail -L -o $@ "https://github.com/mvdan/sh/releases/download/v$(SHFMT_VERSION)/shfmt_v$(SHFMT_VERSION)_`echo $* | tr - _`"

test/e2e/bats: cache/bats-core-$(BATS_COMMIT).tar.gz
	mkdir -p $@
	tar -C $@ --strip-components 1 -xzf $<

cache/bats-core-$(BATS_COMMIT).tar.gz:
	# Use 2opremio's fork until https://github.com/bats-core/bats-core/pull/255 is merged
	curl --fail -L -o $@ https://github.com/2opremio/bats-core/archive/$(BATS_COMMIT).tar.gz

generate: generate-codegen generate-deploy

generate-codegen:
	./hack/update/generated.sh

generate-deploy: pkg/install/generated_templates.gogen.go
	cd deploy && go run ../pkg/install/generate.go deploy
	cp ./deploy/flux-helm-release-crd.yaml ./chart/helm-operator/crds/helmrelease.yaml

check-generated: generate-deploy pkg/install/generated_templates.gogen.go
	git diff --exit-code -- pkg/install/generated_templates.gogen.go
	./hack/update/verify.sh

build-docs:
	@cd docs && docker build -t flux-docs .

test-docs: build-docs
	@docker run -it flux-docs /usr/bin/linkchecker _build/html/index.html

serve-docs: build-docs
	@echo Stating docs website on http://localhost:${DOCS_PORT}/_build/html/index.html
	@docker run -i -p ${DOCS_PORT}:8000 -e USER_ID=$$UID flux-docs
