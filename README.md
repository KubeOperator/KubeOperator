# KubeOperator - Kubernetes 集群部署和管理平台

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

## 什么是 KubeOperator？

KubeOperator 是一个开源项目，帮助运维人员通过 Web 控制台，在完全离线环境下实现生产级别的 Kubernetes 集群的可视化部署及生命周期管理。KubeOperator 尤其适合用于在 VMware 云平台上部署和管理生产级别的 Kubernetes 集群。

## 为什么需要 KubeOperator？

-  按需创建：快速创建 Kubernetes 集群，避免手动创建。
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率。
-  按需修补：快速升级和修补 Kubernetes 集群，保证安全性和版本同步。
-  健康检查：主动式健康检测，及时发现潜在问题。
-  自我修复：通过重建故障节点确保集群可用性。
-  Multi-AZ支持：通过把集群节点分布在不同的故障域上确保集群的高可用。

## KubeOperator 的版本规划

 v1.0

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 支持 Kubernetes 常用组件安装，包括 Registry，Promethus，Dashboard等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 作为网络方案；

 v2.0

- [ ] 支持集群扩容；
- [ ] 支持集群升级；
- [ ] 支持调用 VMware vCenter 接口自动创建集群所需资源；
- [ ] 支持 VMware vSAN / VMFS 作为持久化存储；

 v2.1
 
- [ ] 支持多主多节点模式（Multi AZ，分布在不同故障域） 
- [ ] 支持 VMware NSX-T；
- [ ] 支持集群备份及恢复；

## 安装 KubeOperator

 [安装手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/install.md)

## 使用 KubeOperator

 [使用手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/user-guide.md)

## 致谢

- 感谢 [kubeasz](https://github.com/easzlab/kubeasz) 提供各种 Kubernetes Ansible 脚本.

## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.
