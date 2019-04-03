FROM node as stage-build
WORKDIR /data
ADD ./package.json /data/package.json
RUN npm i
ADD . /data
RUN npm run-script build

FROM nginx:alpine
COPY --from=stage-build /data/dist /opt/fit2openshift-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf
