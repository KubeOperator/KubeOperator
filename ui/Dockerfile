FROM nginx:alpine
# Build dist in other ci
COPY ./dist/ /opt/fit2openshift-ui
COPY nginx.conf /etc/nginx/conf.d/default.conf
