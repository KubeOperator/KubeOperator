FROM centos:7

MAINTAINER OpenShift Team <dev@lists.openshift.redhat.com>

USER root

LABEL name="openshift/origin-ansible" \
      summary="OpenShift's installation and configuration tool" \
      description="A containerized openshift-ansible image to use in CI - includes necessary packages to test clusters on AWS/GCP/Azure" \
      url="https://github.com/openshift/openshift-ansible" \
      io.k8s.display-name="openshift-ansible" \
      io.k8s.description="A containerized openshift-ansible image to use in CI - includes necessary packages to test clusters on AWS/GCP/Azure" \
      io.openshift.expose-services="" \
      io.openshift.tags="openshift,install,upgrade,ansible" \
      atomic.run="once"

ENV USER_UID=1001 \
    HOME=/opt/app-root/src \
    WORK_DIR=/usr/share/ansible/openshift-ansible \
    OPTS="-v"

# Add image scripts and files for running as a system container
COPY images/installer/root /
# Add origin repo for including the oc client
COPY images/installer/origin-extra-root /
# Install openshift-ansible RPMs
RUN yum install -y epel-release && \
    rm -rf /etc/yum.repos.d/centos-openshift-origin.repo && \
    yum-config-manager --enable built > /dev/null && \
    INSTALL_PKGS="openssh google-cloud-sdk azure-cli-2.0.46" \
    yum install --setopt=tsflags=nodocs -y $INSTALL_PKGS openshift-ansible-test && \
    yum clean all

RUN /usr/local/bin/user_setup \
 && rm /usr/local/bin/usage.ocp

USER ${USER_UID}

WORKDIR ${WORK_DIR}
ENTRYPOINT [ "/usr/local/bin/entrypoint" ]
CMD [ "/usr/local/bin/run" ]
