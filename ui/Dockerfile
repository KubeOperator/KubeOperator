FROM node:10-alpine as stage-build
WORKDIR /data

COPY ./package.json /data/package.json
COPY ./package-lock.json /data/package-lock.json
RUN npm install
COPY . /data
RUN npm run-script build

FROM nginx:alpine

RUN rm -rf libjpeg.so.8.2.2 libturbojpeg.so.0.2.0

COPY --from=stage-build /data/dist /opt/kubeOperator-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf