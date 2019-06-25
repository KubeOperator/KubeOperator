# KubeOperator - K8S 集群部署和管理平台

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.4.2.0-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

## 什么是 KubeOperator？

KubeOperator 是一个开源项目，帮助运维人员通过 Web 控制台，在完全离线环境下实现 K8S 集群的可视化部署及管理。

## 为什么需要 KubeOperator? 

K8S 是未来的 Linux。K8S 高可用集群部署、升级的门槛较高，尤其是在完全离线环境下。

## KubeOperator 有什么功能？

- [x] 提供 K8S 标准版 及 OpenShift 社区版的离线包仓库；
- [x] 支持两种部署模式：单节点模式，高可用模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 支持 K8S 常用组件安装，包括 EFK，Harbor，Promethus，Dashboard等；
- [x] 提供简易明了的集群运行状况面板；
- [x] 支持 NFS 作为外部持久化存储；
- [x] 支持 vSAN 作为外部持久化存储；
- [x] 支持 AD/LDAP 对接(仅 OpenShift)；
- [ ] 支持其他外部持久化存储（比如 Ceph，Gluster；
- [ ] 支持 F5 Big-IP 对接；
- [ ] 支持集群的升级；
- [ ] 支持集群进行扩容；
- [ ] 支持集群的备份及恢复；

## KubeOperator 支持哪些 K8S 版本？

- [x] K8S 1.13.5
- [x] OpenShift OKD 3.11

## KubeOperator 背后是谁在支持？

[FIT2CLOUD](https://www.fit2cloud.com) 是 KubeOperator 项目的发起者及核心贡献者。该项目由第一开源堡垒机 [Jumpserver](http://www.jumpserver.org/) 的原班团队打造。

## KubeOperator 的架构

![架构图](https://raw.githubusercontent.com/fit2anything/KubeOperator/master/docs/images/overview.png)

## 安装 KubeOperator

 [安装手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/install.md)

## 使用 KubeOperator

 [使用手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/user-guide.md)
 
## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.

## 致谢

- 感谢 [kubeasz](https://github.com/easzlab/kubeasz) 提供各种 K8S Ansible 脚本.
