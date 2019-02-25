# Copyright (c) 2016-present Sonatype, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM centos:centos7

MAINTAINER Sonatype <cloud-ops@sonatype.com>

LABEL vendor=Sonatype \
      com.sonatype.license="Apache License, Version 2.0" \
      com.sonatype.name="Nexus Repository Manager base image"

ARG NEXUS_VERSION=3.15.0-01
ARG NEXUS_DOWNLOAD_URL=https://download.sonatype.com/nexus/3/nexus-${NEXUS_VERSION}-unix.tar.gz
ARG NEXUS_DOWNLOAD_SHA256_HASH=7bee286f980eb485738110480daf14915c5eeeee3c805904b4814a2669fb712a

#ARG JAVA_URL=http://download.oracle.com/otn-pub/java/jdk/8u192-b12/750e1c8617c5452694857ad95c3ee230/server-jre-8u192-linux-x64.tar.gz
ARG JAVA_URL=http://fit2cloud-openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/server-jre-8u192-linux-x64.tar.gz 
ARG JAVA_DOWNLOAD_SHA256_HASH=3d811a5ec65dc6fc261f488757bae86ecfe285a79992363b016f60cdb4dbe7e6

ENV JAVA_HOME=/opt/java

# configure nexus runtime
ENV SONATYPE_DIR=/opt/sonatype
ENV NEXUS_HOME=${SONATYPE_DIR}/nexus \
    NEXUS_DATA=/nexus-data \
    NEXUS_CONTEXT='' \
    SONATYPE_WORK=${SONATYPE_DIR}/sonatype-work \
    DOCKER_TYPE='docker'

ARG NEXUS_REPOSITORY_MANAGER_COOKBOOK_VERSION="release-0.5.20180828-231916.863899c"
ARG NEXUS_REPOSITORY_MANAGER_COOKBOOK_URL="https://github.com/sonatype/chef-nexus-repository-manager/releases/download/${NEXUS_REPOSITORY_MANAGER_COOKBOOK_VERSION}/chef-nexus-repository-manager.tar.gz"
ENV  OFFLINE_PACKAGE=/opt/fit2openshift/offline_packages
ADD ./okd-3.11-meta/solo.json.erb /var/chef/solo.json.erb
ADD ./okd-3.11-meta/dependencies/chef-14.9.13-1.el7.x86_64.rpm /opt

RUN rpm -ihv /opt/chef-14.9.13-1.el7.x86_64.rpm

# Install using chef-solo
RUN echo hello chef \
    && /opt/chef/embedded/bin/erb /var/chef/solo.json.erb > /var/chef/solo.json \
    && chef-solo \
       --recipe-url ${NEXUS_REPOSITORY_MANAGER_COOKBOOK_URL} \
       --json-attributes /var/chef/solo.json \
    && rpm -qa *chef* | xargs rpm -e \
    && rpm --rebuilddb \
    && rm -rf /etc/chef \
    && rm -rf /opt/chefdk \
    && rm -rf /var/cache/yum \
    && rm -rf /var/chef

COPY ./okd-3.11-meta/nexus-data  ${NEXUS_DATA}

COPY ./resource ${OFFLINE_PACKAGE}/okd-v3.11.0-0cbc58b

RUN chmod -R 777 ${NEXUS_DATA} && chmod -R 777 ${OFFLINE_PACKAGE}


EXPOSE 8081 8082
USER nexus

ENV INSTALL4J_ADD_VM_PARAMS="-Xms1200m -Xmx1200m -XX:MaxDirectMemorySize=2g -Djava.util.prefs.userRoot=${NEXUS_DATA}/javaprefs"

CMD ["sh", "-c", "${SONATYPE_DIR}/start-nexus-repository-manager.sh"]
