# KubeOperator - Kubernetes 集群部署和管理平台

![Total visitor](https://visitor-count-badge.herokuapp.com/total.svg?repo_id=kubeoperator)
![Visitors in today](https://visitor-count-badge.herokuapp.com/today.svg?repo_id=kubeoperator)
[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

## 什么是 KubeOperator？

KubeOperator 是一个开源项目，帮助运维人员通过 Web-based UI，在完全离线环境下部署和管理生产级别的 Kubernetes 集群。KubeOperator 尤其适合在 VMware 云平台上部署和管理 Kubernetes 集群，实现 Kubernetes as a Service。
![overview](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/KubeOperator.jpeg?raw=true)

## 为什么需要 KubeOperator？

-  按需创建：调用云平台 API，一键快速创建和部署 Kubernetes 集群 (即 Kubernetes as a Service)；
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率；
-  按需修补：快速升级和修补 Kubernetes 集群，保证集群安全性，并与社区最新版本同步；
-  自我修复：通过重建故障节点确保集群可用性；
-  离线部署：持续更新包括 Kubernetes 及常用组件的离线包；
-  Multi-AZ 支持：通过把集群节点分布在不同的故障域上确保集群的高可用；

## KubeOperator 的版本规划

 v1.0.0 （已发布）

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 支持 Kubernetes 常用组件安装，包括 Registry，Promethus，Dashboard等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 作为网络方案；
- [x] 支持 Kubernetes 集群手动部署模式（自行准备主机和 NFS）；

 v2.0.0 （开发中）

- [x] 支持调用 VMware vCenter API 自动创建集群主机；
- [x] 支持 VMware vSAN 作为持久化存储；
- [x] 支持通过 F5 BIG-IP Controller 对外暴露服务（Nodeport mode）；
- [x] 支持 Kubernetes 集群扩缩容；

 v3.0.0 （计划中）
 
- [ ] 支持集群升级；
- [ ] 支持集群备份及恢复；
- [ ] 支持 Multi AZ，主节点分布在不同的故障域；
- [ ] 支持 VMware NSX-T；

 v4.0.0 （计划中）
 
- [ ] 支持 Openstack 云平台；
- [ ] 支持 Ceph 作为持久化存储；
- [ ] 支持 Calico 作为网络方案；

## 离线包

KubeOperator 会持续维护包括操作系统、Kubernetes 及常用组件的完整离线包，该离线包能在完全离线的网络环境下部署和升级，并保证最终的一致性。离线包版本命名和 Kubernetes 版本命名保持一致。目前已发布的离线包：

- [v1.5.0](https://github.com/KubeOperator/KubeOperator/blob/master/offline-package/v1.5.0.md)
- [v1.5.2](https://github.com/KubeOperator/KubeOperator/blob/master/offline-package/v1.5.2.md)
 
 ## KubeOperator 使用指南

安装手册：
- [KubeOperator 安装手册（离线包方式）](https://github.com/KubeOperator/KubeOperator/blob/master/docs/install-1.md)
- [KubeOperator 安装手册（OVA 方式）](https://github.com/KubeOperator/KubeOperator/blob/master/docs/install-2.md)

使用手册：
- [KubeOperator 使用手册（手动模式）](https://github.com/KubeOperator/KubeOperator/blob/master/docs/user-guide-1.md)
- [KubeOperator 使用手册（自动模式）](https://github.com/KubeOperator/KubeOperator/blob/master/docs/user-guide-2.md)

> Note:
> - 手动模式：用户需要自行准备主机及 NFS 存储，适合在物理机环境下部署。
> - 自动模式：依赖于 VMware 云平台（包括 vSAN），用户只需绑定 vCenter 相关账号和密码，设置好部署计划，即可实现一键部署、扩容和故障自愈。注：V4.0.0 版本开始会支持 Openstack 云平台。

REST API：
- [KubeOperator REST API](https://github.com/KubeOperator/KubeOperator/blob/master/docs/restapi.md)

## 致谢

- 感谢 [kubeasz](https://github.com/easzlab/kubeasz) 提供各种 Kubernetes Ansible 脚本.

## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.
