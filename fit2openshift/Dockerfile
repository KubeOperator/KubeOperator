FROM registry.fit2cloud.com/public/python:v3

WORKDIR /opt/fit2openshift-app

COPY ./api/requirements /tmp/requirements
COPY ./gen-key.exp /opt

RUN rm -f /usr/bin/python && ln -s /usr/bin/python2 /usr/bin/python
RUN yum -y install epel-release && cd /tmp/requirements && \
    yum -y install $(cat rpm_requirements.txt)
RUN cd /opt && python3 -m venv py3
RUN cd /tmp/requirements && /opt/py3/bin/pip install --upgrade pip setuptools && \
    /opt/py3/bin/pip install -r requirements.txt -i https://mirrors.ustc.edu.cn/pypi/web/simple
RUN sed -i "s@'uri': True@'uri': False@g" /opt/py3/lib/python3.6/site-packages/django/db/backends/sqlite3/base.py
RUN yum -y install nginx
RUN rm -rf /usr/share/nginx/html/*
RUN yum -y install expect
RUN /opt/gen-key.exp


ENV LANG=zh_CN.UTF-8
ENV LC_ALL=zh_CN.UTF-8
ENV VENV=/opt/py3
ENV APP_DIR=/opt/fi2openshift-app
ENV PYTHONOPTIMIZE=1
ENV C_FORCE_ROOT=1

RUN mkdir -p /root/.ssh/

ADD ./api/ /opt/fit2openshift-app
ADD ./ui/dist/fit2openshift-portal /usr/share/nginx/html

EXPOSE 8080 80
CMD  sh -c 'nginx && python3 entrypoint.py start'
