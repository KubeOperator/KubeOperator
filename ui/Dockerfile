FROM node as stage-build
ADD . /data
WORKDIR /data
RUN npm i
RUN npm run-script build

FROM nginx:alpine
COPY --from=stage-build /data/dist /opt/fit2openshift-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf
