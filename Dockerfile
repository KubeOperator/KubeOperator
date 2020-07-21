FROM golang:1.14-alpine as stage-build
LABEL stage=stage-build
WORKDIR /build/ko
ARG GOPROXY
ENV GOPROXY=$GOPROXY
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk update \
  && apk add git \
  && apk add make \
  && apk add bash
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build_server_linux

FROM alpine:3.11

COPY --from=stage-build /build/ko/dist/etc /etc/
COPY --from=stage-build /build/ko/dist/usr /usr/

EXPOSE 8080

CMD ["ko-server"]
