FROM node:10-alpine as stage-build
WORKDIR /data
ADD ./package.json /data/package.json
RUN SASS_BINARY_SITE=https://npm.taobao.org/mirrors/node-sass/ npm install node-sass@4.10.0
RUN npm --registry=https://registry.npm.taobao.org --disturl=https://npm.taobao.org/dist install
ADD . /data
RUN npm run-script build

FROM nginx:alpine
COPY --from=stage-build /data/dist /opt/fit2openshift-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf
