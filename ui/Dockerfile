FROM node:10-alpine as stage-build
WORKDIR /data

COPY ./package.json /data/package.json
COPY ./package-lock.json /data/package-lock.json
RUN npm install
COPY . /data
RUN npm run-script build

FROM nginx:alpine
ARG GOARCH

RUN echo > /etc/apk/repositories && echo -e "https://dl-cdn.alpinelinux.org/alpine/latest-stable/main\nhttps://dl-cdn.alpinelinux.org/alpine/latest-stable/community" >> /etc/apk/repositories && apk update && apk upgrade

RUN apk del curl

RUN rm -rf /usr/lib/libpng16* /usr/lib/libjpeg* /usr/lib/libturbojpeg* /usr/lib/libncursesw*

RUN wget https://kubeoperator.oss-cn-beijing.aliyuncs.com/busybox/1.33.1/$GOARCH/busybox.tar.gz \
    && tar zxvf busybox.tar.gz -C /bin \
    && rm -rf busybox.tar.gz

COPY --from=stage-build /data/dist /opt/kubeOperator-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf