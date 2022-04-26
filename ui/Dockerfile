FROM node:10-alpine as stage-build
WORKDIR /data

COPY ./package.json /data/package.json
COPY ./package-lock.json /data/package-lock.json
RUN npm install
COPY . /data
RUN npm run-script build

FROM kubeoperator/nginx:alpine

ENV ENV="/root/.ashrc"
RUN echo 'umask 0027' >> /root/.ashrc

COPY --from=stage-build /data/dist /opt/kubeOperator-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]