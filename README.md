# 目前该项目仍在开发中，未正式发布。



# Welcome to FIT2OPENSHIFT (Fit to OpenShift)

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.4.2.0-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)


## 什么是 FIT2OPENSHIFT？

FIT2OpenShift 是一个开源项目，帮助运维人员通过 Web 控制台，在完全离线环境下实现 [OpenShift](https://www.openshift.com/) 社区版（[okd](https://www.okd.io/)）集群的可视化部署及运维管理。

## 为什么需要 FIT2OpenShift? 

如果说 K8S 是未来的 Linux，那么 OpenShift 是首选的 K8S 发行版。 OpenShift 不像 CentOS，其社区版部署及后续持续运营的门槛非常高，尤其在国内特殊的网络环境下。

## FIT2OPENSHIFT 有什么功能？

1. 提供各种部署模板，比如 all-in-one 模式，单 master 模式，多 master 模式等；
2. 根据部署模板输入机器IP、密码等信息；
3. 触发部署后，系统将可视化展示集群部署进展和结果；
4. 部署完成后，可以对集群进行扩容和缩容；
5. 集群的持续升级及回滚；
6. 简易明了的集群运行状况面板；
7. 集群备份及恢复；
8. 成熟的网络、存储及安全方案;
9. 提供 OpenShift 社区版各个版本离线包仓库；

最重要的是，上述的功能可以在完全离线环境下，以可视化方式来实现。

## FIT2OPENSHIFT 的架构图

![架构图](https://raw.githubusercontent.com/fit2anything/fit2openshift/master/docs/overview.jpg)


## 开始使用 FIT2OPENSHIFT

 [安装文档](https://github.com/fit2anything/fit2openshift/blob/master/installer/README.md)



## FIT2OPENSHIFT 项目背后是谁在支持？

[FIT2CLOUD](https://www.fit2cloud.com) 是FIT2OpenShift 项目的发起者及核心贡献者。该项目由第一开源堡垒机 [Jumpserver](http://www.jumpserver.org/) 的原帮团队倾力打造。[FIT2CLOUD](https://www.fit2cloud.com) 飞致云是国内拥有「红帽认证OpenShift 管理员」最多的公司。
