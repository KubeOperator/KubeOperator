# KubeOperator - K8S 集群部署和管理平台

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)



## 什么是 KubeOperator？

KubeOperator 是一个开源项目，帮助运维人员通过 Web 控制台，在完全离线环境下实现 K8S 集群的可视化部署及管理。

## 为什么需要 KubeOperator？

-  按需创建：快速创建 k8s 集群，避免手动创建。
-  按需伸缩：快速伸缩 k8s 集群，优化资源使用效率。
-  按需修补：快速修补 k8s 集群，保持安全性。
-  健康检查：主动式健康检测，及时发现潜在问题。
-  自我修复：通过重建故障节点确保集群可用性。
-  Multi-AZ支持：通过把集群节点分布在不同的故障域上确保集群的高可用。

## KubeOperator 的版本规划

 v1.0

- [x] 提供 K8S 标准版的离线包仓库；
- [x] 支持两种部署模式：一主多节点模式，多主多节点模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 支持 K8S 常用组件安装，包括 EFK，Registry，Promethus，Dashboard等；
- [x] 提供简易明了的集群运行状况面板；
- [x] 支持 NFS 作为外部持久化存储；
- [x] 支持 vSAN 作为外部持久化存储；
- [x] 支持 Flannel 作为网络方案；
- [x] 支持集群进行扩容；

 v1.1

- [ ] 支持集群的升级；
- [ ] 支持节点操作系统补丁升级；
- [ ] 支持集群的备份及恢复；
- [ ] 支持集群的自动模式: 调用云平台接口自动创建K8S 集群所需要的资源（支持 VMware )；

v2.0

- [ ] Multi-AZ支持。


## KubeOperator 支持哪些 K8S 版本？

- [x] K8S 1.13.5

## KubeOperator 背后是谁在支持？

[FIT2CLOUD](https://www.fit2cloud.com) 是 KubeOperator 项目的发起者及核心贡献者。该项目由第一开源堡垒机 [Jumpserver](http://www.jumpserver.org/) 的原班团队打造。


## 安装 KubeOperator

 [安装手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/install.md)

## 使用 KubeOperator

 [使用手册](https://github.com/fit2anything/KubeOperator/blob/master/docs/user-guide.md)
 
## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.

## 致谢

- 感谢 [kubeasz](https://github.com/easzlab/kubeasz) 提供各种 K8S Ansible 脚本.
