FROM node as stage-build
WORKDIR /data
ADD ./package.json /data/package.json
RUN alias cnpm="npm --registry=https://registry.npm.taobao.org \
--cache=$HOME/.npm/.cache/cnpm \
--disturl=https://npm.taobao.org/dist \
--userconfig=$HOME/.cnpmrc"
RUN cnpm i
ADD . /data
RUN npm run-script build

FROM nginx:alpine
COPY --from=stage-build /data/dist /opt/fit2openshift-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf
