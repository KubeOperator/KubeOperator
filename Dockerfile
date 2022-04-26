FROM golang:1.17 as stage-build
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

FROM kubeoperator/server:euler2sp10-20220111
ARG GOARCH

RUN useradd -u 2004 kops && usermod -aG kops kops

COPY --from=stage-build /build/ko/dist/home /home/
COPY --from=stage-build /build/ko/dist/usr /usr/
COPY --from=stage-build /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

RUN cd /usr/local/bin && \
    chmod -R 550 ko-server validator_linux_$GOARCH encrypt_linux_$GOARCH && \
    chown -R kops:kops ko-server validator_linux_$GOARCH encrypt_linux_$GOARCH && \
    chown -R kops:kops /home/kops

WORKDIR /home/kops

EXPOSE 8080

USER kops

CMD ["ko-server"]
