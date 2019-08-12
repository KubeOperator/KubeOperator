FROM registry.fit2cloud.com/public/python:v3
MAINTAINER Fit2Openshift Team <support@fit2cloud.com>
WORKDIR /opt/kubeOperator-api

RUN echo -e "[mysql]\nname=mysql\nbaseurl=https://mirrors.tuna.tsinghua.edu.cn/mysql/yum/mysql57-community-el6/\ngpgcheck=0\nenabled=1" > /etc/yum.repos.d/mysql.repo
COPY ./requirements /tmp/requirements
RUN cd /tmp/requirements && yum -y install epel-release && yum -y install $(cat rpm_requirements.txt)
RUN cd /tmp/requirements && \
     pip install --upgrade pip setuptools -i https://mirrors.aliyun.com/pypi/simple/ && \
     pip install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/ || pip install -r requirements.txt
RUN test -f /root/.ssh/id_rsa || ssh-keygen -f /root/.ssh/id_rsa -t rsa -P ''
RUN echo -e "Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile /dev/null" > /root/.ssh/config
COPY ./requirements/install_curl.sh /opt/kubeOperator-api
COPY ./requirements/install_terraform.sh /opt/kubeOperator-api
RUN bash /opt/kubeOperator-api/install_curl.sh
RUN bash /opt/kubeOperator-api/install_terraform.sh
COPY extra/dnsmasq.conf /etc/
COPY extra/kubeops.dns.conf /etc/dnsmasq.d
COPY . /opt/kubeOperator-api
RUN echo > config.yml

ENV LANG=zh_CN.UTF-8
ENV LC_ALL=zh_CN.UTF-8
ENV PYTHONOPTIMIZE=1
ENV C_FORCE_ROOT=1

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
