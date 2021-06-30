GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist
GOGINDATA=go-bindata
GOBUILDMODE=pie

KO_SERVER_NAME=ko-server
KO_CONFIG_DIR=etc/ko
KO_BIN_DIR=usr/local/bin
KO_DATA_DIR=usr/local/lib/ko

GOPROXY="https://goproxy.cn,direct"

build_server_linux:
	GOOS=linux GOARCH=$(GOARCH) --buildmode=$(GOBUILDMODE) -trimpath $(GOGINDATA) -o ./pkg/i18n/locales.go -pkg i18n ./locales/...
	GOOS=linux GOARCH=$(GOARCH)  $(GOBUILD) -o $(BUILDDIR)/$(KO_BIN_DIR)/$(KO_SERVER_NAME) main.go
	mkdir -p $(BUILDDIR)/$(KO_CONFIG_DIR) && cp -r  $(BASEPATH)/conf/app.yaml $(BUILDDIR)/$(KO_CONFIG_DIR)
	mkdir -p $(BUILDDIR)/$(KO_DATA_DIR)
	cp -r  $(BASEPATH)/migration $(BUILDDIR)/$(KO_DATA_DIR)

docker_ui:
	docker build -t kubeoperator/ui:master --build-arg GOARCH=$(GOARCH) ./ui --no-cache

docker_server:
	docker build -t kubeoperator/server:master --build-arg GOPROXY=$(GOPROXY) --build-arg GOARCH=$(GOARCH) --build-arg XPACK="no" .

docker_server_xpack:
	docker build -t kubeoperator/server:master --build-arg GOPROXY=$(GOPROXY) --build-arg GOARCH=$(GOARCH) --build-arg XPACK="yes" .

clean:
	rm -fr ./dist
	go clean
