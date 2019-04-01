FROM nginx:alpine
# Build dist in other ci
COPY ./dist/fit2openshift-portal /opt/fit2openshift-portal
COPY nginx.conf /etc/nginx/conf.d/default.conf
