FROM golang:1.14 as stage-build
LABEL stage=stage-build
WORKDIR /build/ko
ARG GOPROXY
ARG GOARCH
ARG XPACK

ENV GOARCH=$GOARCH
ENV GO111MODULE=on
ENV GOOS=linux
ENV CGO_ENABLED=1

RUN apt-get update && apt-get install unzip

COPY go.mod go.sum ./
RUN go mod download

RUN wget https://github.com/go-bindata/go-bindata/archive/v3.1.3.zip -O /tmp/go-bindata.zip  \
    && cd /tmp \
    && unzip  /tmp/go-bindata.zip  \
    && cd /tmp/go-bindata-3.1.3 \
    && go build \
    && cd go-bindata \
    && go build \
    && cp go-bindata /go/bin

RUN export PATH=$PATH:$GOPATH/bin

COPY . .
RUN make build_server_linux GOARCH=$GOARCH

RUN if [ "$XPACK" = "yes" ] ; then  cd xpack && sed -i 's/ ..\/KubeOperator/ \..\/..\/ko/g' go.mod && make build_linux GOARCH=$GOARCH && cp -r dist/* ../dist/  ; fi

FROM kubeoperator/euleros:2.1
ARG GOARCH

RUN if [ "$GOARCH" = "amd64" ] ; then \
        echo > /etc/yum.repos.d/Euler-Base.repo; \
        echo -e "[base]\nname=EulerOS-2.0SP5 base\nbaseurl=http://mirrors.huaweicloud.com/euler/2.5/os/x86_64/\nenabled=1\ngpgcheck=1\ngpgkey=http://mirrors.huaweicloud.com/euler/2.5/os/RPM-GPG-KEY-EulerOS" >> /etc/yum.repos.d/Euler-Base.repo; \
    fi

RUN if [ "$GOARCH" = "arm64" ] ; then \
        echo > /etc/yum.repos.d/Euler-Base.repo; \
        echo -e "[base]\nname=EulerOS-2.0SP8 base\nbaseurl=http://repo.huaweicloud.com/euler/2.8/os/aarch64/\nenabled=1\ngpgcheck=1\ngpgkey=http://repo.huaweicloud.com/euler/2.8/os/RPM-GPG-KEY-EulerOS" >> /etc/yum.repos.d/Euler-Base.repo; \
    fi

RUN cd /usr/local/bin && \
    yum install -y wget && \
    wget https://fit2cloud-support.oss-cn-beijing.aliyuncs.com/xpack-license/validator_linux_$GOARCH && \
    wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/ko-encrypt/encrypt_linux_$GOARCH && \
    chmod 500 validator_linux_$GOARCH && \
    chmod 500 encrypt_linux_$GOARCH && \
    yum remove -y wget && \
    yum clean all && \
    rm -rf /var/cache/yum/*

WORKDIR /

COPY --from=stage-build /build/ko/dist/etc /etc/
COPY --from=stage-build /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

COPY --from=stage-build /build/ko/dist/etc /etc/
COPY --from=stage-build /build/ko/dist/usr /usr/

RUN chmod 550 /usr/local/bin/ko-server

EXPOSE 8080

CMD ["ko-server"]
