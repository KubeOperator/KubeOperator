
FROM centos:centos7

MAINTAINER fit2openshift <support@fit2cloud.com>

ENV  OFFLINE_PACKAGE=/opt/fit2openshift/offline_packages

COPY ./resource ${OFFLINE_PACKAGE}/okd-v3.11.0-0cbc58b

ADD ./entrypoint.sh /opt

RUN chmod +x /opt/entrypoint.sh

RUN  chmod -R 777 ${OFFLINE_PACKAGE}

USER root

ENTRYPOINT ["/opt/entrypoint.sh"]
 
