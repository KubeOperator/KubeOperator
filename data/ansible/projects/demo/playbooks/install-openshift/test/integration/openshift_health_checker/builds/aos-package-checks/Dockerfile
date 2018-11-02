FROM test-target-base

RUN yum install -y rpm-build rpmdevtools createrepo && \
    rpmdev-setuptree && \
    mkdir -p /mnt/localrepo
ADD root /

# we will build some RPMs that can be used to break yum update in tests.
RUN cd /root/rpmbuild/SOURCES && \
    mkdir break-yum-update-1.0 && \
    tar zfc foo.tgz break-yum-update-1.0 && \
    rpmbuild -bb /root/break-yum-update.spec  && \
    yum install -y /root/rpmbuild/RPMS/noarch/break-yum-update-1.0-1.noarch.rpm && \
    rpmbuild -bb /root/break-yum-update-2.spec  && \
    mkdir /mnt/localrepo/break-yum && \
    cp /root/rpmbuild/RPMS/noarch/break-yum-update-1.0-2.noarch.rpm /mnt/localrepo/break-yum && \
    createrepo /mnt/localrepo/break-yum

# we'll also build some RPMs that can be used to exercise OCP package version tests.
RUN cd /root/rpmbuild/SOURCES && \
    mkdir atomic-openshift-3.2 && \
    mkdir atomic-openshift-3.3 && \
    tar zfc ose.tgz atomic-openshift-3.{2,3} && \
    rpmbuild -bb /root/ose-3.2.spec  && \
    rpmbuild -bb /root/ose-3.3.spec  && \
    mkdir /mnt/localrepo/ose-3.{2,3} && \
    cp /root/rpmbuild/RPMS/noarch/atomic-openshift*-3.2-1.noarch.rpm /mnt/localrepo/ose-3.2 && \
    cp /root/rpmbuild/RPMS/noarch/{openvswitch-2.4,docker-1.10}-1.noarch.rpm /mnt/localrepo/ose-3.2 && \
    createrepo /mnt/localrepo/ose-3.2 && \
    cp /root/rpmbuild/RPMS/noarch/atomic-openshift*-3.3-1.noarch.rpm /mnt/localrepo/ose-3.3 && \
    cp /root/rpmbuild/RPMS/noarch/{openvswitch-2.4,docker-1.10}-1.noarch.rpm /mnt/localrepo/ose-3.3 && \
    createrepo /mnt/localrepo/ose-3.3
