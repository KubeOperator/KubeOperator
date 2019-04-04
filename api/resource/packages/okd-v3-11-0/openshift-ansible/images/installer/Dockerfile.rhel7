FROM rhel7.3:7.3-released

MAINTAINER OpenShift Team <dev@lists.openshift.redhat.com>

USER root

# Playbooks, roles, and their dependencies are installed from packages.
RUN INSTALL_PKGS="openshift-ansible atomic-openshift-clients python-boto python2-boto3 python2-crypto openssl iproute httpd-tools" \
 && x86_EXTRA_RPMS=$(if [ "$(uname -m)" == "x86_64" ]; then echo -n google-cloud-sdk ; fi) \
 && yum repolist > /dev/null \
 && yum-config-manager --enable rhel-7-server-ose-3.7-rpms \
 && yum-config-manager --enable rhel-7-server-rh-common-rpms \
 && yum install -y java-1.8.0-openjdk-headless \
 && yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS $x86_EXTRA_RPMS \
 && rpm -q $INSTALL_PKGS $x86_EXTRA_RPMS \
 && yum clean all

LABEL name="openshift3/ose-ansible" \
      summary="OpenShift's installation and configuration tool" \
      description="A containerized openshift-ansible image to let you run playbooks to install, upgrade, maintain and check an OpenShift cluster" \
      url="https://github.com/openshift/openshift-ansible" \
      io.k8s.display-name="openshift-ansible" \
      io.k8s.description="A containerized openshift-ansible image to let you run playbooks to install, upgrade, maintain and check an OpenShift cluster" \
      io.openshift.expose-services="" \
      io.openshift.tags="openshift,install,upgrade,ansible" \
      com.redhat.component="aos3-installation-docker" \
      version="v3.6.0" \
      release="1" \
      architecture="x86_64" \
      atomic.run="once"

ENV USER_UID=1001 \
    HOME=/opt/app-root/src \
    WORK_DIR=/usr/share/ansible/openshift-ansible \
    ANSIBLE_CONFIG=/usr/share/ansible/openshift-ansible/ansible.cfg \
    OPTS="-v"

# Add image scripts and files for running as a system container
COPY root /

RUN /usr/local/bin/user_setup \
 && mv /usr/local/bin/usage{.ocp,}

USER ${USER_UID}

WORKDIR ${WORK_DIR}
ENTRYPOINT [ "/usr/local/bin/entrypoint" ]
CMD [ "/usr/local/bin/run" ]
