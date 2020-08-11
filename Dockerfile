FROM golang:1.14-alpine as stage-build
LABEL stage=stage-build
WORKDIR /build/ko
ARG GOPROXY
ARG GOARCH

ENV GOARCH=$GOARCH
ENV GOPROXY=$GOPROXY
ENV GOARCH=$GOARCH
ENV GO111MODULE=on
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN  apk update \
  && apk add git \
  && apk add make \
  && apk add bash \
  && apk add binutils-gold
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

FROM alpine:3.11

COPY --from=stage-build /build/ko/dist/etc /etc/
COPY --from=stage-build /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

COPY --from=stage-build /build/ko/dist/etc /etc/
COPY --from=stage-build /build/ko/dist/usr /usr/

EXPOSE 8080

CMD ["ko-server"]
