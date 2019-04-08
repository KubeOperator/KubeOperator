FROM registry.fit2cloud.com/public/python:v3
MAINTAINER Fit2Openshift Team <support@fit2cloud.com>
WORKDIR /opt/fit2openshift-api

COPY ./requirements /tmp/requirements
RUN mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
RUN wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-6.repo
RUN wget-O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-6.repo
RUN rpm -ivh https://repo.mysql.com/mysql57-community-release-el6.rpm
RUN cd /tmp/requirements && yum -y install $(cat rpm_requirements.txt)
RUN cd /tmp/requirements && pip install --upgrade pip setuptools && \
     pip install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/ || pip install -r requirements.txt
RUN test -f /root/.ssh/id_rsa || ssh-keygen -f /root/.ssh/id_rsa -t rsa -P ''
RUN echo -e "Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile /dev/null" > /root/.ssh/config

COPY . /opt/fit2openshift-api
RUN echo > config.yml

ENV LANG=zh_CN.UTF-8
ENV LC_ALL=zh_CN.UTF-8
ENV PYTHONOPTIMIZE=1
ENV C_FORCE_ROOT=1

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
